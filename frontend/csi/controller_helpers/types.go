// Copyright 2022 NetApp, Inc. All Rights Reserved.

package controllerhelpers

//go:generate mockgen -destination=../../../mocks/mock_frontend/mock_csi/mock_controller_helpers/mock_controller_helpers.go github.com/netapp/trident/frontend/csi/controller_helpers ControllerHelper

import (
	"context"

	"github.com/netapp/trident/config"
	"github.com/netapp/trident/storage"
)

const (
	KubernetesHelper = "k8s_csi_helper"
	PlainCSIHelper   = "plain_csi_helper"

	EventTypeNormal  = "Normal"
	EventTypeWarning = "Warning"
)

type Feature string

// ControllerHelper is the common interface used by the "helper" objects used by
// the CSI controller.  The controller_helpers supply CO-specific details at certain
// points of CSI workflows.
type ControllerHelper interface {
	// GetVolumeConfig accepts the attributes of a volume being requested by the CSI
	// provisioner, adds in any CO-specific details about the new volume, and returns
	// a VolumeConfig structure as needed by Trident to create a new volume.
	GetVolumeConfig(
		ctx context.Context, name string, sizeBytes int64, parameters map[string]string,
		protocol config.Protocol, accessModes []config.AccessMode, volumeMode config.VolumeMode, fsType string,
		requisiteTopology, preferredTopology, accessibleTopology []map[string]string,
	) (*storage.VolumeConfig, error)

	// GetSnapshotConfig accepts the attributes of a snapshot being requested by the CSI
	// provisioner, adds in any CO-specific details about the new volume, and returns
	// a SnapshotConfig structure as needed by Trident to create a new snapshot.
	GetSnapshotConfig(volumeName, snapshotName string) (*storage.SnapshotConfig, error)

	// GetNodeTopologyLabels returns topology labels for a given node
	// Example: map[string]string{"topology.kubernetes.io/region": "us-east1"}
	GetNodeTopologyLabels(ctx context.Context, nodeName string) (map[string]string, error)

	// RecordVolumeEvent accepts the name of a CSI volume and writes the specified
	// event message in a manner appropriate to the container orchestrator.
	RecordVolumeEvent(ctx context.Context, name, eventType, reason, message string)

	// RecordNodeEvent accepts the name of a CSI node and writes the specified
	// event message in a manner appropriate to the container orchestrator.
	RecordNodeEvent(ctx context.Context, name, eventType, reason, message string)

	// SupportsFeature accepts a CSI feature and returns true if the feature is supported.
	SupportsFeature(ctx context.Context, feature Feature) bool

	// Version returns the version of the CO this helper is managing, or the supported
	// CSI version in the plain-CSI case.  This value is reported in Trident's telemetry.
	Version() string
}
