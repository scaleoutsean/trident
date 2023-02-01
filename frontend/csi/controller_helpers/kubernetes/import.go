// Copyright 2022 NetApp, Inc. All Rights Reserved.

package kubernetes

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/netapp/trident/frontend/csi"
	. "github.com/netapp/trident/logger"
	"github.com/netapp/trident/storage"
	"github.com/netapp/trident/utils"
)

/////////////////////////////////////////////////////////////////////////////
//
// This file contains the event handlers that import existing volumes from backends.
//
/////////////////////////////////////////////////////////////////////////////

func (h *helper) ImportVolume(
	ctx context.Context, request *storage.ImportVolumeRequest,
) (*storage.VolumeExternal, error) {
	Logc(ctx).WithField("request", request).Debug("ImportVolume")

	// Get PVC from ImportVolumeRequest
	jsonData, err := base64.StdEncoding.DecodeString(request.PVCData)
	if err != nil {
		return nil, fmt.Errorf("the pvcData field does not contain valid base64-encoded data: %v", err)
	}

	existingVol, err := h.orchestrator.GetVolumeByInternalName(request.InternalName, ctx)
	if err == nil {
		return nil, utils.FoundError(fmt.Sprintf("PV %s already exists for volume %s",
			existingVol, request.InternalName))
	}

	claim := &v1.PersistentVolumeClaim{}
	err = json.Unmarshal(jsonData, &claim)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON PVC: %v", err)
	}

	Logc(ctx).WithFields(log.Fields{
		"PVC":               claim.Name,
		"PVC_storageClass":  getStorageClassForPVC(claim),
		"PVC_accessModes":   claim.Spec.AccessModes,
		"PVC_annotations":   claim.Annotations,
		"PVC_volume":        claim.Spec.VolumeName,
		"PVC_requestedSize": claim.Spec.Resources.Requests[v1.ResourceStorage],
	}).Debug("ImportVolume: received PVC data.")

	if err = h.checkValidStorageClassReceived(ctx, claim); err != nil {
		return nil, err
	}

	if len(claim.Namespace) == 0 {
		return nil, fmt.Errorf("a valid PVC namespace is required for volume import")
	}

	// Lookup backend ID from given name
	backend, err := h.orchestrator.GetBackend(ctx, request.Backend)
	if err != nil {
		return nil, fmt.Errorf("could not find backend %s; %v", request.Backend, err)
	}

	// Set required annotations
	if claim.Annotations == nil {
		claim.Annotations = map[string]string{}
	}
	claim.Annotations[AnnNotManaged] = strconv.FormatBool(request.NoManage)
	claim.Annotations[AnnImportOriginalName] = request.InternalName
	claim.Annotations[AnnImportBackendUUID] = backend.BackendUUID
	claim.Annotations[AnnStorageProvisioner] = csi.Provisioner

	// Set the PVC's storage field to the actual volume size.
	volExternal, err := h.orchestrator.GetVolumeExternal(ctx, request.InternalName, request.Backend)
	if err != nil {
		return nil, fmt.Errorf("volume import failed to get size of volume: %v", err)
	}

	if claim.Spec.Resources.Requests == nil {
		claim.Spec.Resources.Requests = v1.ResourceList{}
	}
	claim.Spec.Resources.Requests[v1.ResourceStorage] = resource.MustParse(volExternal.Config.Size)
	Logc(ctx).WithFields(log.Fields{
		"size":      volExternal.Config.Size,
		"claimSize": claim.Spec.Resources.Requests[v1.ResourceStorage],
	}).Debug("Volume import determined volume size")

	pvc, pvcErr := h.createImportPVC(ctx, claim)
	if pvcErr != nil {
		Logc(ctx).WithFields(log.Fields{
			"claim": claim.Name,
			"error": pvcErr,
		}).Warn("ImportVolume: error occurred during PVC creation.")
		return nil, pvcErr
	}
	Logc(ctx).WithField("PVC", pvc).Debug("ImportVolume: created pending PVC.")

	pvName := fmt.Sprintf("pvc-%s", pvc.GetUID())

	// Wait for the import to happen and the PV to be created by the sidecar
	// after which we can then request the volume object from the core to return
	_, err = h.waitForCachedPVByName(ctx, pvName, ImportPVCacheWaitPeriod)
	if err != nil {
		return nil, fmt.Errorf("error waiting for PV %s; %v", pvName, err)
	}

	volume, err := h.orchestrator.GetVolume(ctx, pvName)
	if err != nil {
		return nil, fmt.Errorf("error getting volume %s; %v", pvName, err)
	}
	return volume, nil
}

func (h *helper) createImportPVC(
	ctx context.Context, claim *v1.PersistentVolumeClaim,
) (*v1.PersistentVolumeClaim, error) {
	Logc(ctx).WithFields(log.Fields{
		"claim":     claim,
		"namespace": claim.Namespace,
	}).Debug("CreateImportPVC: ready to create PVC")

	createOpts := metav1.CreateOptions{}
	pvc, err := h.kubeClient.CoreV1().PersistentVolumeClaims(claim.Namespace).Create(ctx, claim, createOpts)
	if err != nil {
		return nil, fmt.Errorf("error occurred during PVC creation: %v", err)
	}

	return pvc, nil
}
