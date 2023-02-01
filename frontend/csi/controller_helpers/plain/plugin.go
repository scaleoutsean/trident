// Copyright 2022 NetApp, Inc. All Rights Reserved.
package plain

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/netapp/trident/config"
	"github.com/netapp/trident/core"
	"github.com/netapp/trident/frontend"
	frontendcommon "github.com/netapp/trident/frontend/common"
	"github.com/netapp/trident/frontend/csi"
	controllerhelpers "github.com/netapp/trident/frontend/csi/controller_helpers"
	. "github.com/netapp/trident/logger"
	"github.com/netapp/trident/storage"
)

type helper struct {
	orchestrator core.Orchestrator
}

// NewHelper instantiates this plugin.
func NewHelper(orchestrator core.Orchestrator) frontend.Plugin {
	log.Info("Initializing plain CSI helper frontend.")

	return &helper{
		orchestrator: orchestrator,
	}
}

// Activate starts this Trident frontend.
func (h *helper) Activate() error {
	log.Info("Activating plain CSI helper frontend.")

	// Configure telemetry
	config.OrchestratorTelemetry.Platform = string(config.PlatformCSI)
	config.OrchestratorTelemetry.PlatformVersion = h.Version()

	return nil
}

// Deactivate stops this Trident frontend.
func (h *helper) Deactivate() error {
	log.Info("Deactivating plain CSI helper frontend.")
	return nil
}

// GetName returns the name of this Trident frontend.
func (h *helper) GetName() string {
	return string(controllerhelpers.PlainCSIHelper)
}

// Version returns the version of this Trident frontend (the Trident version).
func (h *helper) Version() string {
	return csi.Version
}

// GetVolumeConfig accepts the attributes of a volume being requested by the CSI
// provisioner, finds or creates/registers a matching storage class, and returns
// a VolumeConfig structure as needed by Trident to create a new volume.
func (h *helper) GetVolumeConfig(
	ctx context.Context, name string, sizeBytes int64, parameters map[string]string,
	protocol config.Protocol, accessModes []config.AccessMode, volumeMode config.VolumeMode, fsType string,
	requisiteTopology, preferredTopology, accessibleTopology []map[string]string,
) (*storage.VolumeConfig, error) {
	accessMode := frontendcommon.CombineAccessModes(accessModes)

	if parameters == nil {
		parameters = make(map[string]string)
	}

	if _, ok := parameters["fstype"]; !ok {
		parameters["fstype"] = fsType
	}

	// Find a matching storage class, or register a new one
	scConfig, err := frontendcommon.GetStorageClass(ctx, parameters, h.orchestrator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "could not create a storage class from volume request")
	}

	// Create the volume config from all available info from the CSI request
	return frontendcommon.GetVolumeConfig(name, scConfig.Name, sizeBytes, parameters, protocol, accessMode, volumeMode,
		requisiteTopology, preferredTopology)
}

// GetSnapshotConfig accepts the attributes of a snapshot being requested by the CSI
// provisioner and returns a SnapshotConfig structure as needed by Trident to create a new snapshot.
func (h *helper) GetSnapshotConfig(volumeName, snapshotName string) (*storage.SnapshotConfig, error) {
	return &storage.SnapshotConfig{
		Version:    config.OrchestratorAPIVersion,
		Name:       snapshotName,
		VolumeName: volumeName,
	}, nil
}

func (h *helper) GetNodeTopologyLabels(ctx context.Context, nodeName string) (map[string]string, error) {
	return map[string]string{}, nil
}

// RecordVolumeEvent accepts the name of a CSI volume and writes the specified
// event message to the debug log.
func (h *helper) RecordVolumeEvent(ctx context.Context, name, eventType, reason, message string) {
	Logc(ctx).WithFields(log.Fields{
		"name":      name,
		"eventType": eventType,
		"reason":    reason,
		"message":   message,
	}).Debug("Volume event.")
}

// RecordNodeEvent accepts the name of a CSI node and writes the specified
// event message to the debug log.
func (h *helper) RecordNodeEvent(ctx context.Context, name, eventType, reason, message string) {
	Logc(ctx).WithFields(log.Fields{
		"name":      name,
		"eventType": eventType,
		"reason":    reason,
		"message":   message,
	}).Debug("Node event.")
}

// SupportsFeature accepts a CSI feature and returns true if the
// feature exists and is supported.
func (h *helper) SupportsFeature(_ context.Context, feature controllerhelpers.Feature) bool {
	if supported, ok := features[feature]; ok {
		return supported
	} else {
		return false
	}
}
