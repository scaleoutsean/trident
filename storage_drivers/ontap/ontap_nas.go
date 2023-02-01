// Copyright 2022 NetApp, Inc. All Rights Reserved.

package ontap

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/RoaringBitmap/roaring"
	log "github.com/sirupsen/logrus"

	tridentconfig "github.com/netapp/trident/config"
	. "github.com/netapp/trident/logger"
	"github.com/netapp/trident/storage"
	sa "github.com/netapp/trident/storage_attribute"
	drivers "github.com/netapp/trident/storage_drivers"
	"github.com/netapp/trident/storage_drivers/ontap/api"
	"github.com/netapp/trident/utils"
)

// //////////////////////////////////////////////////////////////////////////////////////////
// /             _____________________
// /            |   <<Interface>>    |
// /            |       ONTAPI       |
// /            |____________________|
// /                ^             ^
// /     Implements |             | Implements
// /   ____________________    ____________________
// /  |  ONTAPAPIREST     |   |  ONTAPAPIZAPI     |
// /  |___________________|   |___________________|
// /  | +API: RestClient  |   | +API: *Client     |
// /  |___________________|   |___________________|
// /
// //////////////////////////////////////////////////////////////////////////////////////////

// //////////////////////////////////////////////////////////////////////////////////////////
// Drivers that offer dual support are to call ONTAP REST or ZAPI's
// via abstraction layer (ONTAPI interface)
// //////////////////////////////////////////////////////////////////////////////////////////

// NASStorageDriver is for NFS and SMB storage provisioning
type NASStorageDriver struct {
	initialized bool
	Config      drivers.OntapStorageDriverConfig
	API         api.OntapAPI
	telemetry   *Telemetry

	physicalPools map[string]storage.Pool
	virtualPools  map[string]storage.Pool
}

func (d *NASStorageDriver) GetConfig() *drivers.OntapStorageDriverConfig {
	return &d.Config
}

func (d *NASStorageDriver) GetAPI() api.OntapAPI {
	return d.API
}

func (d *NASStorageDriver) GetTelemetry() *Telemetry {
	return d.telemetry
}

// Name is for returning the name of this driver
func (d *NASStorageDriver) Name() string {
	return drivers.OntapNASStorageDriverName
}

// BackendName returns the name of the backend managed by this driver instance
func (d *NASStorageDriver) BackendName() string {
	if d.Config.BackendName == "" {
		// Use the old naming scheme if no name is specified
		return CleanBackendName("ontapnas_" + d.Config.DataLIF)
	} else {
		return d.Config.BackendName
	}
}

// Initialize from the provided config
func (d *NASStorageDriver) Initialize(
	ctx context.Context, driverContext tridentconfig.DriverContext, configJSON string,
	commonConfig *drivers.CommonStorageDriverConfig, backendSecret map[string]string, backendUUID string,
) error {
	if commonConfig.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Initialize", "Type": "NASStorageDriver"}
		Logc(ctx).WithFields(fields).Debug(">>>> Initialize")
		defer Logc(ctx).WithFields(fields).Debug("<<<< Initialize")
	}

	// Initialize the driver's CommonStorageDriverConfig
	d.Config.CommonStorageDriverConfig = commonConfig

	// Parse the config
	config, err := InitializeOntapConfig(ctx, driverContext, configJSON, commonConfig, backendSecret)
	if err != nil {
		return fmt.Errorf("error initializing %s driver: %v", d.Name(), err)
	}
	d.Config = *config

	// Unit tests mock the API layer, so we only use the real API interface if it doesn't already exist.
	if d.API == nil {
		d.API, err = InitializeOntapDriver(ctx, config)
		if err != nil {
			return fmt.Errorf("error initializing %s driver: %v", d.Name(), err)
		}
	}

	d.Config = *config

	d.physicalPools, d.virtualPools, err = InitializeStoragePoolsCommon(ctx, d,
		d.getStoragePoolAttributes(ctx), d.BackendName())
	if err != nil {
		return fmt.Errorf("could not configure storage pools: %v", err)
	}

	// Validate the none, true/false values
	err = d.validate(ctx)
	if err != nil {
		return fmt.Errorf("error validating %s driver: %v", d.Name(), err)
	}

	// Set up the autosupport heartbeat
	d.telemetry = NewOntapTelemetry(ctx, d)
	d.telemetry.Telemetry = tridentconfig.OrchestratorTelemetry
	d.telemetry.TridentBackendUUID = backendUUID
	d.telemetry.Start(ctx)

	d.initialized = true
	return nil
}

func (d *NASStorageDriver) Initialized() bool {
	return d.initialized
}

func (d *NASStorageDriver) Terminate(ctx context.Context, backendUUID string) {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Terminate", "Type": "NASStorageDriver"}
		Logc(ctx).WithFields(fields).Debug(">>>> Terminate")
		defer Logc(ctx).WithFields(fields).Debug("<<<< Terminate")
	}
	if d.Config.AutoExportPolicy {
		policyName := getExportPolicyName(backendUUID)

		if err := d.API.ExportPolicyDestroy(ctx, policyName); err != nil {
			Logc(ctx).Warn(err)
		}
	}
	if d.telemetry != nil {
		d.telemetry.Stop()
	}
	d.initialized = false
}

// Validate the driver configuration and execution environment
func (d *NASStorageDriver) validate(ctx context.Context) error {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "validate", "Type": "NASStorageDriver"}
		Logc(ctx).WithFields(fields).Debug(">>>> validate")
		defer Logc(ctx).WithFields(fields).Debug("<<<< validate")
	}

	if err := validateReplicationConfig(ctx, d.Config.ReplicationPolicy, d.Config.ReplicationSchedule,
		d.API); err != nil {
		return fmt.Errorf("replication validation failed: %v", err)
	}

	if err := ValidateNASDriver(ctx, d.API, &d.Config); err != nil {
		return fmt.Errorf("driver validation failed: %v", err)
	}

	if err := ValidateStoragePrefix(*d.Config.StoragePrefix); err != nil {
		return err
	}

	if err := ValidateStoragePools(
		ctx, d.physicalPools, d.virtualPools, d, api.MaxNASLabelLength,
	); err != nil {
		return fmt.Errorf("storage pool validation failed: %v", err)
	}

	return nil
}

// Create a volume with the specified options
func (d *NASStorageDriver) Create(
	ctx context.Context, volConfig *storage.VolumeConfig, storagePool storage.Pool, volAttributes map[string]sa.Request,
) error {
	name := volConfig.InternalName

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Create",
			"Type":   "NASStorageDriver",
			"name":   name,
			"attrs":  volAttributes,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> Create")
		defer Logc(ctx).WithFields(fields).Debug("<<<< Create")
	}

	// If the volume already exists, bail out
	volExists, err := d.API.VolumeExists(ctx, name)
	if err != nil {
		return fmt.Errorf("error checking for existing volume: %v", err)
	}
	if volExists {
		return drivers.NewVolumeExistsError(name)
	}

	// If volume shall be mirrored, check that the SVM is peered with the other side
	if volConfig.PeerVolumeHandle != "" {
		if err = checkSVMPeered(ctx, volConfig, d.API.SVMName(), d.API); err != nil {
			return err
		}
	}

	// Get candidate physical pools
	physicalPools, err := getPoolsForCreate(ctx, volConfig, storagePool, volAttributes, d.physicalPools, d.virtualPools)
	if err != nil {
		return err
	}

	// Get options
	opts := d.GetVolumeOpts(ctx, volConfig, volAttributes)

	// get options with default fallback values
	// see also: ontap_common.go#PopulateConfigurationDefaults
	var (
		spaceReserve      = utils.GetV(opts, "spaceReserve", storagePool.InternalAttributes()[SpaceReserve])
		snapshotPolicy    = utils.GetV(opts, "snapshotPolicy", storagePool.InternalAttributes()[SnapshotPolicy])
		snapshotReserve   = utils.GetV(opts, "snapshotReserve", storagePool.InternalAttributes()[SnapshotReserve])
		unixPermissions   = utils.GetV(opts, "unixPermissions", storagePool.InternalAttributes()[UnixPermissions])
		snapshotDir       = utils.GetV(opts, "snapshotDir", storagePool.InternalAttributes()[SnapshotDir])
		exportPolicy      = utils.GetV(opts, "exportPolicy", storagePool.InternalAttributes()[ExportPolicy])
		securityStyle     = utils.GetV(opts, "securityStyle", storagePool.InternalAttributes()[SecurityStyle])
		encryption        = utils.GetV(opts, "encryption", storagePool.InternalAttributes()[Encryption])
		tieringPolicy     = utils.GetV(opts, "tieringPolicy", storagePool.InternalAttributes()[TieringPolicy])
		qosPolicy         = storagePool.InternalAttributes()[QosPolicy]
		adaptiveQosPolicy = storagePool.InternalAttributes()[AdaptiveQosPolicy]
	)

	snapshotReserveInt, err := GetSnapshotReserve(snapshotPolicy, snapshotReserve)
	if err != nil {
		return fmt.Errorf("invalid value for snapshotReserve: %v", err)
	}

	// Determine volume size in bytes
	requestedSize, err := utils.ConvertSizeToBytes(volConfig.Size)
	if err != nil {
		return fmt.Errorf("could not convert volume size %s: %v", volConfig.Size, err)
	}
	sizeBytes, err := strconv.ParseUint(requestedSize, 10, 64)
	if err != nil {
		return fmt.Errorf("%v is an invalid volume size: %v", volConfig.Size, err)
	}
	sizeBytes, err = GetVolumeSize(sizeBytes, storagePool.InternalAttributes()[Size])
	if err != nil {
		return err
	}
	// Get the flexvol size based on the snapshot reserve
	flexvolSize := calculateFlexvolSizeBytes(ctx, name, sizeBytes, snapshotReserveInt)

	size := strconv.FormatUint(flexvolSize, 10)

	if _, _, checkVolumeSizeLimitsError := drivers.CheckVolumeSizeLimits(
		ctx, sizeBytes, d.Config.CommonStorageDriverConfig,
	); checkVolumeSizeLimitsError != nil {
		return checkVolumeSizeLimitsError
	}

	enableSnapshotDir, err := strconv.ParseBool(snapshotDir)
	if err != nil {
		return fmt.Errorf("invalid boolean value for snapshotDir: %v", err)
	}

	enableEncryption, err := GetEncryptionValue(encryption)
	if err != nil {
		return fmt.Errorf("invalid boolean value for encryption: %v", err)
	}

	if tieringPolicy == "" {
		tieringPolicy = d.API.TieringPolicyValue(ctx)
	}

	if d.Config.AutoExportPolicy {
		exportPolicy = getExportPolicyName(storagePool.Backend().BackendUUID())
	}

	qosPolicyGroup, err := api.NewQosPolicyGroup(qosPolicy, adaptiveQosPolicy)
	if err != nil {
		return err
	}
	volConfig.QosPolicy = qosPolicy
	volConfig.AdaptiveQosPolicy = adaptiveQosPolicy

	Logc(ctx).WithFields(log.Fields{
		"name":              name,
		"size":              size,
		"spaceReserve":      spaceReserve,
		"snapshotPolicy":    snapshotPolicy,
		"snapshotReserve":   snapshotReserveInt,
		"unixPermissions":   unixPermissions,
		"snapshotDir":       enableSnapshotDir,
		"exportPolicy":      exportPolicy,
		"securityStyle":     securityStyle,
		"encryption":        utils.GetPrintableBoolPtrValue(enableEncryption),
		"tieringPolicy":     tieringPolicy,
		"qosPolicy":         qosPolicy,
		"adaptiveQosPolicy": adaptiveQosPolicy,
	}).Debug("Creating Flexvol.")

	createErrors := make([]error, 0)
	physicalPoolNames := make([]string, 0)

	for _, physicalPool := range physicalPools {
		aggregate := physicalPool.Name()
		physicalPoolNames = append(physicalPoolNames, aggregate)

		if aggrLimitsErr := checkAggregateLimits(
			ctx, aggregate, spaceReserve, sizeBytes, d.Config, d.GetAPI(),
		); aggrLimitsErr != nil {
			errMessage := fmt.Sprintf("ONTAP-NAS pool %s/%s; error: %v", storagePool.Name(), aggregate, aggrLimitsErr)
			Logc(ctx).Error(errMessage)
			createErrors = append(createErrors, fmt.Errorf(errMessage))
			continue
		}

		labels, err := storagePool.GetLabelsJSON(ctx, storage.ProvisioningLabelTag, api.MaxNASLabelLength)
		if err != nil {
			return err
		}

		// Create the volume
		err = d.API.VolumeCreate(
			ctx, api.Volume{
				Aggregates: []string{
					aggregate,
				},
				Comment:         labels,
				Encrypt:         enableEncryption,
				ExportPolicy:    exportPolicy,
				Name:            name,
				Qos:             qosPolicyGroup,
				Size:            size,
				SpaceReserve:    spaceReserve,
				SnapshotPolicy:  snapshotPolicy,
				SecurityStyle:   securityStyle,
				SnapshotReserve: snapshotReserveInt,
				TieringPolicy:   tieringPolicy,
				UnixPermissions: unixPermissions,
				DPVolume:        volConfig.IsMirrorDestination,
			})

		if err != nil {
			if api.IsVolumeCreateJobExistsError(err) {
				return nil
			}

			errMessage := fmt.Sprintf("ONTAP-NAS pool %s/%s; error creating volume %s: %v", storagePool.Name(),
				aggregate, name, err)
			Logc(ctx).Error(errMessage)
			createErrors = append(createErrors, fmt.Errorf(errMessage))
			continue
		}

		// Disable '.snapshot' to allow official mysql container's chmod-in-init to work
		if !enableSnapshotDir {
			if err := d.API.VolumeDisableSnapshotDirectoryAccess(ctx, name); err != nil {
				createErrors = append(createErrors,
					fmt.Errorf("ONTAP-NAS-FLEXGROUP pool %s; error disabling snapshot directory access for volume %v: %v",
						storagePool.Name(), name, err))
				return drivers.NewBackendIneligibleError(name, createErrors, physicalPoolNames)
			}
		}
		// If a DP volume, skip mounting the volume
		if volConfig.IsMirrorDestination {
			return nil
		}

		// Mount the volume at the specified junction
		if err := d.API.VolumeMount(ctx, name, "/"+name); err != nil {
			return err
		}

		return nil
	}

	// All physical pools that were eligible ultimately failed, so don't try this backend again
	return drivers.NewBackendIneligibleError(name, createErrors, physicalPoolNames)
}

// CreateClone creates a volume clone
func (d *NASStorageDriver) CreateClone(
	ctx context.Context, _, cloneVolConfig *storage.VolumeConfig, storagePool storage.Pool,
) error {
	// Ensure the volume exists
	flexvol, err := d.API.VolumeInfo(ctx, cloneVolConfig.CloneSourceVolumeInternal)
	if err != nil {
		return err
	}

	if d.GetConfig().DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":      "CreateClone",
			"Type":        "NASStorageDriver",
			"name":        cloneVolConfig.InternalName,
			"source":      cloneVolConfig.CloneSourceVolumeInternal,
			"snapshot":    cloneVolConfig.CloneSourceSnapshot,
			"storagePool": storagePool,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> CreateClone")
		defer Logc(ctx).WithFields(fields).Debug("<<<< CreateClone")
	}

	opts := d.GetVolumeOpts(context.Background(), cloneVolConfig, make(map[string]sa.Request))

	labels := flexvol.Comment

	// How "splitOnClone" value gets set:
	// In the Core we first check clone's VolumeConfig for splitOnClone value
	// If it is not set then (again in Core) we check source PV's VolumeConfig for splitOnClone value
	// If we still don't have splitOnClone value then HERE we check for value in the source PV's Storage/Virtual Pool
	// If the value for "splitOnClone" is still empty then HERE we set it to backend config's SplitOnClone value

	// Attempt to get splitOnClone value based on storagePool (source Volume's StoragePool)
	var storagePoolSplitOnCloneVal string

	if storage.IsStoragePoolUnset(storagePool) {
		// Set the base label
		storagePoolTemp := &storage.StoragePool{}
		storagePoolTemp.SetAttributes(map[string]sa.Offer{
			sa.Labels: sa.NewLabelOffer(d.GetConfig().Labels),
		})
		labels, err = storagePoolTemp.GetLabelsJSON(ctx, storage.ProvisioningLabelTag, api.MaxNASLabelLength)
		if err != nil {
			return err
		}
	} else {
		storagePoolSplitOnCloneVal = storagePool.InternalAttributes()[SplitOnClone]
	}

	// If storagePoolSplitOnCloneVal is still unknown, set it to backend's default value
	if storagePoolSplitOnCloneVal == "" {
		storagePoolSplitOnCloneVal = d.GetConfig().SplitOnClone
	}

	split, err := strconv.ParseBool(utils.GetV(opts, "splitOnClone", storagePoolSplitOnCloneVal))
	if err != nil {
		return fmt.Errorf("invalid boolean value for splitOnClone: %v", err)
	}

	qosPolicy := utils.GetV(opts, "qosPolicy", "")
	adaptiveQosPolicy := utils.GetV(opts, "adaptiveQosPolicy", "")
	qosPolicyGroup, err := api.NewQosPolicyGroup(qosPolicy, adaptiveQosPolicy)
	if err != nil {
		return err
	}

	Logc(ctx).WithField("splitOnClone", split).Debug("Creating volume clone.")
	return cloneFlexvol(ctx, cloneVolConfig.InternalName, cloneVolConfig.CloneSourceVolumeInternal,
		cloneVolConfig.CloneSourceSnapshot, labels, split, d.GetConfig(), d.GetAPI(), qosPolicyGroup)
}

// Destroy the volume
func (d *NASStorageDriver) Destroy(ctx context.Context, volConfig *storage.VolumeConfig) error {
	name := volConfig.InternalName

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Destroy",
			"Type":   "NASStorageDriver",
			"name":   name,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> Destroy")
		defer Logc(ctx).WithFields(fields).Debug("<<<< Destroy")
	}

	// TODO: If this is the parent of one or more clones, those clones have to split from this
	// volume before it can be deleted, which means separate copies of those volumes.
	// If there are a lot of clones on this volume, that could seriously balloon the amount of
	// utilized space. Is that what we want? Or should we just deny the delete, and force the
	// user to keep the volume around until all of the clones are gone? If we do that, need a
	// way to list the clones. Maybe volume inspect.

	// First, check to see if the volume has already been deleted out of band
	volumeExists, err := d.API.VolumeExists(ctx, name)
	if err != nil {
		return fmt.Errorf("error checking for volume %v: %v", name, err)
	}
	if !volumeExists {
		// Not an error if the volume no longer exists
		Logc(ctx).WithField("volume", name).Warn("Volume already deleted.")
		return nil
	}

	// If flexvol has been a snapmirror destination
	if err := d.API.SnapmirrorDeleteViaDestination(name, d.API.SVMName()); err != nil {
		if !api.IsNotFoundError(err) {
			return err
		}
	}

	if err := d.API.VolumeDestroy(ctx, name, true); err != nil {
		return err
	}

	return nil
}

func (d *NASStorageDriver) Import(
	ctx context.Context, volConfig *storage.VolumeConfig, originalName string,
) error {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":       "Import",
			"Type":         "NASStorageDriver",
			"originalName": originalName,
			"newName":      volConfig.InternalName,
			"notManaged":   volConfig.ImportNotManaged,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> Import")
		defer Logc(ctx).WithFields(fields).Debug("<<<< Import")
	}

	// Ensure the volume exists
	flexvol, err := d.API.VolumeInfo(ctx, originalName)
	if api.IsVolumeReadError(err) {
		return err
	}

	// Validate the volume is what it should be
	if !api.IsVolumeIdAttributesReadError(err) {
		if flexvol.AccessType != "" && flexvol.AccessType != "rw" {
			Logc(ctx).WithField("originalName", originalName).Error("Could not import volume, type is not rw.")
			return fmt.Errorf("volume %s type is %s, not rw", originalName, flexvol.AccessType)
		}

		// Make sure we're not importing a volume without a junction path when not managed
		if volConfig.ImportNotManaged {
			if flexvol.JunctionPath == "" {
				return fmt.Errorf("junction path is not set for volume %s", originalName)
			}
		}
	} else {
		if volConfig.ImportNotManaged {
			return err
		}
	}

	// Get the volume size
	if api.IsVolumeSpaceAttributesReadError(err) {
		Logc(ctx).WithField("originalName", originalName).Errorf("Could not import volume, size not available")
		return err
	}
	volConfig.Size = flexvol.Size

	// Rename the volume if Trident will manage its lifecycle
	if !volConfig.ImportNotManaged {
		if err := d.API.VolumeRename(ctx, originalName, volConfig.InternalName); err != nil {
			return err
		}
	}

	// Update the volume labels if Trident will manage its lifecycle
	if !volConfig.ImportNotManaged {
		if storage.AllowPoolLabelOverwrite(storage.ProvisioningLabelTag, flexvol.Comment) {
			if err := d.API.VolumeSetComment(ctx, volConfig.InternalName, originalName, ""); err != nil {
				return err
			}
		}
	}

	// Modify unix-permissions of the volume if Trident will manage its lifecycle
	if !volConfig.ImportNotManaged {
		// unixPermissions specified in PVC annotation takes precedence over backend's unixPermissions config
		unixPerms := volConfig.UnixPermissions

		// ONTAP supports unix permissions only with security style "mixed" on SMB volume.
		switch d.Config.NASType {
		case sa.SMB:
			if unixPerms == "" && d.Config.SecurityStyle == "mixed" {
				unixPerms = d.Config.UnixPermissions
			} else if d.Config.SecurityStyle == "ntfs" {
				unixPerms = ""
			}
		case sa.NFS:
			if unixPerms == "" {
				unixPerms = d.Config.UnixPermissions
			}
		}

		if err := d.API.VolumeModifyUnixPermissions(
			ctx, volConfig.InternalName, originalName, unixPerms,
		); err != nil {
			return err
		}
	}

	return nil
}

// Rename changes the name of a volume
func (d *NASStorageDriver) Rename(ctx context.Context, name, newName string) error {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":  "Rename",
			"Type":    "NASStorageDriver",
			"name":    name,
			"newName": newName,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> Rename")
		defer Logc(ctx).WithFields(fields).Debug("<<<< Rename")
	}

	return d.API.VolumeRename(ctx, name, newName)
}

// Publish the volume to the host specified in publishInfo.  This method may or may not be running on the host
// where the volume will be mounted, so it should limit itself to updating access rules, initiator groups, etc.
// that require some host identity (but not locality) as well as storage controller API access.
func (d *NASStorageDriver) Publish(
	ctx context.Context, volConfig *storage.VolumeConfig, publishInfo *utils.VolumePublishInfo,
) error {
	name := volConfig.InternalName

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":  "Publish",
			"DataLIF": d.Config.DataLIF,
			"Type":    "NASStorageDriver",
			"name":    name,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> Publish")
		defer Logc(ctx).WithFields(fields).Debug("<<<< Publish")
	}

	// Determine mount options (volume config wins, followed by backend config)
	mountOptions := d.Config.NfsMountOptions
	if volConfig.MountOptions != "" {
		mountOptions = volConfig.MountOptions
	}

	// Add fields needed by Attach
	// TODO (akerr) Figure out if this is the behavior we want or if we should be changing the junction path for
	//  managed imports
	// publishInfo.NfsPath = fmt.Sprintf("/%s", name)

	if d.Config.NASType == sa.SMB {
		publishInfo.SMBPath = volConfig.AccessInfo.SMBPath
		publishInfo.SMBServer = d.Config.DataLIF
		publishInfo.FilesystemType = sa.SMB
	} else {
		publishInfo.NfsPath = volConfig.AccessInfo.NfsPath
		publishInfo.NfsServerIP = d.Config.DataLIF
		publishInfo.FilesystemType = sa.NFS
		publishInfo.MountOptions = mountOptions
	}

	return publishShare(ctx, d.API, &d.Config, publishInfo, name, d.API.VolumeModifyExportPolicy)
}

// CanSnapshot determines whether a snapshot as specified in the provided snapshot config may be taken.
func (d *NASStorageDriver) CanSnapshot(_ context.Context, _ *storage.SnapshotConfig, _ *storage.VolumeConfig) error {
	return nil
}

// getFlexvolSnapshot gets a snapshot.  To distinguish between an API error reading the snapshot
// and a non-existent snapshot, this method may return (nil, nil).
func getFlexvolSnapshot(
	ctx context.Context, snapConfig *storage.SnapshotConfig, config *drivers.OntapStorageDriverConfig,
	client api.OntapAPI,
) (*storage.Snapshot, error) {
	internalSnapName := snapConfig.InternalName
	internalVolName := snapConfig.VolumeInternalName

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":       "getFlexvolSnapshot",
			"Type":         "NASStorageDriver",
			"snapshotName": internalSnapName,
			"volumeName":   internalVolName,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> getFlexvolSnapshot")
		defer Logc(ctx).WithFields(fields).Debug("<<<< getFlexvolSnapshot")
	}

	size, err := client.VolumeUsedSize(ctx, internalVolName)
	if err != nil {
		return nil, fmt.Errorf("error reading volume size: %v", err)
	}

	snapshots, err := client.VolumeSnapshotList(ctx, internalVolName)
	if err != nil {
		return nil, err
	}

	for _, snap := range snapshots {

		Logc(ctx).WithFields(log.Fields{
			"snapshotName": internalSnapName,
			"volumeName":   internalVolName,
			"created":      snap.CreateTime,
		}).Debug("Found snapshot.")

		if snap.Name == internalSnapName {
			return &storage.Snapshot{
				Config:    snapConfig,
				Created:   snap.CreateTime,
				SizeBytes: int64(size),
				State:     storage.SnapshotStateOnline,
			}, nil
		}
	}

	return nil, nil
}

// GetSnapshot gets a snapshot.  To distinguish between an API error reading the snapshot
// and a non-existent snapshot, this method may return (nil, nil).
func (d *NASStorageDriver) GetSnapshot(
	ctx context.Context, snapConfig *storage.SnapshotConfig, _ *storage.VolumeConfig,
) (*storage.Snapshot, error) {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":       "GetSnapshot",
			"Type":         "NASStorageDriver",
			"snapshotName": snapConfig.InternalName,
			"volumeName":   snapConfig.VolumeInternalName,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> GetSnapshot")
		defer Logc(ctx).WithFields(fields).Debug("<<<< GetSnapshot")
	}

	return getFlexvolSnapshot(ctx, snapConfig, &d.Config, d.API)
}

// getFlexvolSnapshotList returns the list of snapshots associated with the named volume.
func getFlexvolSnapshotList(
	ctx context.Context, volConfig *storage.VolumeConfig, config *drivers.OntapStorageDriverConfig, client api.OntapAPI,
) ([]*storage.Snapshot, error) {
	internalVolName := volConfig.InternalName

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "getFlexvolSnapshotList",
			"Type":       "NASStorageDriver",
			"volumeName": internalVolName,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> getFlexvolSnapshotList")
		defer Logc(ctx).WithFields(fields).Debug("<<<< getFlexvolSnapshotList")
	}

	size, err := client.VolumeUsedSize(ctx, internalVolName)
	if err != nil {
		return nil, fmt.Errorf("error reading volume size: %v", err)
	}

	snapshots, err := client.VolumeSnapshotList(ctx, internalVolName)
	if err != nil {
		return nil, err
	}

	result := make([]*storage.Snapshot, 0)

	for _, snap := range snapshots {

		Logc(ctx).WithFields(log.Fields{
			"name":       snap.Name,
			"accessTime": snap.CreateTime,
		}).Debug("Snapshot")

		snapshot := &storage.Snapshot{
			Config: &storage.SnapshotConfig{
				Version:            tridentconfig.OrchestratorAPIVersion,
				Name:               snap.Name,
				InternalName:       snap.Name,
				VolumeName:         volConfig.Name,
				VolumeInternalName: volConfig.InternalName,
			},
			Created:   snap.CreateTime,
			SizeBytes: int64(size),
			State:     storage.SnapshotStateOnline,
		}

		result = append(result, snapshot)
	}

	return result, nil
}

// GetSnapshots returns the list of snapshots associated with the specified volume
func (d *NASStorageDriver) GetSnapshots(ctx context.Context, volConfig *storage.VolumeConfig) (
	[]*storage.Snapshot, error,
) {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "GetSnapshots",
			"Type":       "NASStorageDriver",
			"volumeName": volConfig.InternalName,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> GetSnapshots")
		defer Logc(ctx).WithFields(fields).Debug("<<<< GetSnapshots")
	}

	return getFlexvolSnapshotList(ctx, volConfig, &d.Config, d.API)
}

// CreateSnapshot creates a snapshot for the given volume
func (d *NASStorageDriver) CreateSnapshot(
	ctx context.Context, snapConfig *storage.SnapshotConfig, _ *storage.VolumeConfig,
) (*storage.Snapshot, error) {
	internalSnapName := snapConfig.InternalName
	internalVolName := snapConfig.VolumeInternalName

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":       "CreateSnapshot",
			"Type":         "NASStorageDriver",
			"snapshotName": internalSnapName,
			"sourceVolume": internalVolName,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> CreateSnapshot")
		defer Logc(ctx).WithFields(fields).Debug("<<<< CreateSnapshot")
	}

	return createFlexvolSnapshot(ctx, snapConfig, &d.Config, d.API, d.API.VolumeUsedSize)
}

// RestoreSnapshot restores a volume (in place) from a snapshot.
func (d *NASStorageDriver) RestoreSnapshot(
	ctx context.Context, snapConfig *storage.SnapshotConfig, _ *storage.VolumeConfig,
) error {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":       "RestoreSnapshot",
			"Type":         "NASStorageDriver",
			"snapshotName": snapConfig.InternalName,
			"volumeName":   snapConfig.VolumeInternalName,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> RestoreSnapshot")
		defer Logc(ctx).WithFields(fields).Debug("<<<< RestoreSnapshot")
	}

	return RestoreSnapshot(ctx, snapConfig, &d.Config, d.API)
}

// DeleteSnapshot creates a snapshot of a volume.
func (d *NASStorageDriver) DeleteSnapshot(
	ctx context.Context, snapConfig *storage.SnapshotConfig, _ *storage.VolumeConfig,
) error {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":       "DeleteSnapshot",
			"Type":         "NASStorageDriver",
			"snapshotName": snapConfig.InternalName,
			"volumeName":   snapConfig.VolumeInternalName,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> DeleteSnapshot")
		defer Logc(ctx).WithFields(fields).Debug("<<<< DeleteSnapshot")
	}

	if err := d.API.VolumeSnapshotDelete(ctx, snapConfig.InternalName, snapConfig.VolumeInternalName); err != nil {
		if api.IsSnapshotBusyError(err) {
			// Start a split here before returning the error so a subsequent delete attempt may succeed.
			_ = SplitVolumeFromBusySnapshot(ctx, snapConfig, &d.Config, d.API, d.API.VolumeCloneSplitStart)
		}
		// we must return the err, even if we started a split, so the snapshot delete is retried
		return err
	}

	Logc(ctx).WithField("snapshotName", snapConfig.InternalName).Debug("Deleted snapshot.")
	return nil
}

// Get tests for the existence of a volume
func (d *NASStorageDriver) Get(ctx context.Context, name string) error {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Get", "Type": "NASStorageDriver"}
		Logc(ctx).WithFields(fields).Debug(">>>> Get")
		defer Logc(ctx).WithFields(fields).Debug("<<<< Get")
	}

	volExists, err := d.API.VolumeExists(ctx, name)
	if err != nil {
		return fmt.Errorf("error checking for existing volume: %v", err)
	}
	if !volExists {
		Logc(ctx).WithField("Volume", name).Debug("Volume not found.")
		return fmt.Errorf("volume %s does not exist", name)
	}

	return nil
}

// GetStorageBackendSpecs retrieves storage backend capabilities
func (d *NASStorageDriver) GetStorageBackendSpecs(
	_ context.Context, backend storage.Backend,
) error {
	return getStorageBackendSpecsCommon(backend, d.physicalPools, d.virtualPools, d.BackendName())
}

// GetStorageBackendPhysicalPoolNames retrieves storage backend physical pools
func (d *NASStorageDriver) GetStorageBackendPhysicalPoolNames(context.Context) []string {
	return getStorageBackendPhysicalPoolNamesCommon(d.physicalPools)
}

func (d *NASStorageDriver) getStoragePoolAttributes(ctx context.Context) map[string]sa.Offer {
	client := d.GetAPI()
	mirroring, _ := client.IsSVMDRCapable(ctx)
	return map[string]sa.Offer{
		sa.BackendType:      sa.NewStringOffer(d.Name()),
		sa.Snapshots:        sa.NewBoolOffer(true),
		sa.Clones:           sa.NewBoolOffer(true),
		sa.Encryption:       sa.NewBoolOffer(true),
		sa.Replication:      sa.NewBoolOffer(mirroring),
		sa.ProvisioningType: sa.NewStringOffer("thick", "thin"),
	}
}

func (d *NASStorageDriver) GetVolumeOpts(
	ctx context.Context, volConfig *storage.VolumeConfig, requests map[string]sa.Request,
) map[string]string {
	return getVolumeOptsCommon(ctx, volConfig, requests)
}

func (d *NASStorageDriver) GetInternalVolumeName(_ context.Context, name string) string {
	return getInternalVolumeNameCommon(d.Config.CommonStorageDriverConfig, name)
}

func (d *NASStorageDriver) CreatePrepare(ctx context.Context, volConfig *storage.VolumeConfig) {
	createPrepareCommon(ctx, d, volConfig)
}

func (d *NASStorageDriver) CreateFollowup(ctx context.Context, volConfig *storage.VolumeConfig) error {
	var accessPath string
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":       "CreateFollowup",
			"Type":         "NASStorageDriver",
			"name":         volConfig.Name,
			"internalName": volConfig.InternalName,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> CreateFollowup")
		defer Logc(ctx).WithFields(fields).Debug("<<<< CreateFollowup")
	}

	volConfig.MirrorHandle = d.API.SVMName() + ":" + volConfig.InternalName
	if d.Config.NASType == sa.SMB {
		volConfig.AccessInfo.SMBServer = d.Config.DataLIF
		volConfig.FileSystem = sa.SMB
	} else {
		volConfig.AccessInfo.NfsServerIP = d.Config.DataLIF
		volConfig.AccessInfo.MountOptions = strings.TrimPrefix(d.Config.NfsMountOptions, "-o ")
		volConfig.FileSystem = sa.NFS
	}

	// Set correct junction path
	flexvol, err := d.API.VolumeInfo(ctx, volConfig.InternalName)
	if err != nil {
		return err
	}

	if flexvol.JunctionPath == "" {
		if flexvol.AccessType == "rw" || flexvol.AccessType == "dp" {
			// Flexvol is not mounted, we need to mount it
			if d.Config.NASType == sa.SMB {
				volConfig.AccessInfo.SMBPath = ConstructOntapNASSMBVolumePath(ctx, d.Config.SMBShare,
					volConfig.InternalName)
				accessPath = volConfig.AccessInfo.SMBPath
			} else {
				volConfig.AccessInfo.NfsPath = "/" + volConfig.InternalName
				accessPath = volConfig.AccessInfo.NfsPath
			}

			err = d.MountVolume(ctx, volConfig.InternalName, accessPath, flexvol)
			if err != nil {
				return err
			}
		}
	} else {
		if d.Config.NASType == sa.SMB {
			volConfig.AccessInfo.SMBPath = ConstructOntapNASSMBVolumePath(ctx, d.Config.SMBShare,
				flexvol.JunctionPath)
		} else {
			volConfig.AccessInfo.NfsPath = flexvol.JunctionPath
		}
	}
	return nil
}

func (d *NASStorageDriver) GetProtocol(context.Context) tridentconfig.Protocol {
	return tridentconfig.File
}

func (d *NASStorageDriver) StoreConfig(_ context.Context, b *storage.PersistentStorageBackendConfig) {
	drivers.SanitizeCommonStorageDriverConfig(d.Config.CommonStorageDriverConfig)
	b.OntapConfig = &d.Config
}

func (d *NASStorageDriver) GetExternalConfig(ctx context.Context) interface{} {
	return getExternalConfig(ctx, d.Config)
}

// GetVolumeExternal queries the storage backend for all relevant info about
// a single container volume managed by this driver and returns a VolumeExternal
// representation of the volume.
func (d *NASStorageDriver) GetVolumeExternal(
	ctx context.Context, name string,
) (*storage.VolumeExternal, error) {
	volume, err := d.API.VolumeInfo(ctx, name)
	if err != nil {
		return nil, err
	}

	return getVolumeExternalCommon(*volume, *d.Config.StoragePrefix, d.API.SVMName()), nil
}

// GetVolumeExternalWrappers queries the storage backend for all relevant info about
// container volumes managed by this driver.  It then writes a VolumeExternal
// representation of each volume to the supplied channel, closing the channel
// when finished.
func (d *NASStorageDriver) GetVolumeExternalWrappers(
	ctx context.Context, channel chan *storage.VolumeExternalWrapper,
) {
	// Let the caller know we're done by closing the channel
	defer close(channel)

	// Get all volumes matching the storage prefix
	volumes, err := d.API.VolumeListByPrefix(ctx, *d.Config.StoragePrefix)
	if err != nil {
		channel <- &storage.VolumeExternalWrapper{Volume: nil, Error: err}
		return
	}

	// Convert all volumes to VolumeExternal and write them to the channel
	for _, volume := range volumes {
		channel <- &storage.VolumeExternalWrapper{
			Volume: getVolumeExternalCommon(*volume, *d.Config.StoragePrefix, d.API.SVMName()),
			Error:  nil,
		}
	}
}

// GetUpdateType returns a bitmap populated with updates to the driver
func (d *NASStorageDriver) GetUpdateType(ctx context.Context, driverOrig storage.Driver) *roaring.Bitmap {
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "GetUpdateType",
			"Type":   "NASStorageDriver",
		}
		Logc(ctx).WithFields(fields).Debug(">>>> GetUpdateType")
		defer Logc(ctx).WithFields(fields).Debug("<<<< GetUpdateType")
	}

	bitmap := roaring.New()
	dOrig, ok := driverOrig.(*NASStorageDriver)
	if !ok {
		bitmap.Add(storage.InvalidUpdate)
		return bitmap
	}

	if d.Config.Password != dOrig.Config.Password {
		bitmap.Add(storage.PasswordChange)
	}

	if d.Config.Username != dOrig.Config.Username {
		bitmap.Add(storage.UsernameChange)
	}

	if !drivers.AreSameCredentials(d.Config.Credentials, dOrig.Config.Credentials) {
		bitmap.Add(storage.CredentialsChange)
	}

	if !reflect.DeepEqual(d.Config.StoragePrefix, dOrig.Config.StoragePrefix) {
		bitmap.Add(storage.PrefixChange)
	}

	return bitmap
}

// Resize expands the volume size.
func (d *NASStorageDriver) Resize(
	ctx context.Context, volConfig *storage.VolumeConfig, requestedSizeBytes uint64,
) error {
	name := volConfig.InternalName
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":             "Resize",
			"Type":               "NASStorageDriver",
			"name":               name,
			"requestedSizeBytes": requestedSizeBytes,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> Resize")
		defer Logc(ctx).WithFields(fields).Debug("<<<< Resize")
	}

	flexvolSize, err := resizeValidation(ctx, name, requestedSizeBytes, d.API.VolumeExists, d.API.VolumeSize)
	if err != nil {
		return err
	}

	volConfig.Size = strconv.FormatUint(flexvolSize, 10)
	if flexvolSize == requestedSizeBytes {
		return nil
	}

	snapshotReserveInt, err := getSnapshotReserveFromOntap(ctx, name, d.API.VolumeInfo)
	if err != nil {
		Logc(ctx).WithField("name", name).Errorf("Could not get the snapshot reserve percentage for volume")
	}

	newFlexvolSize := calculateFlexvolSizeBytes(ctx, name, requestedSizeBytes, snapshotReserveInt)

	if aggrLimitsErr := checkAggregateLimitsForFlexvol(
		ctx, name, newFlexvolSize, d.Config, d.GetAPI(),
	); aggrLimitsErr != nil {
		return aggrLimitsErr
	}

	if _, _, checkVolumeSizeLimitsError := drivers.CheckVolumeSizeLimits(
		ctx, requestedSizeBytes, d.Config.CommonStorageDriverConfig,
	); checkVolumeSizeLimitsError != nil {
		return checkVolumeSizeLimitsError
	}

	if err := d.API.VolumeSetSize(ctx, name, strconv.FormatUint(newFlexvolSize, 10)); err != nil {
		return err
	}

	volConfig.Size = strconv.FormatUint(requestedSizeBytes, 10)
	return nil
}

func (d *NASStorageDriver) ReconcileNodeAccess(
	ctx context.Context, nodes []*utils.Node, backendUUID string,
) error {
	nodeNames := make([]string, 0)
	for _, node := range nodes {
		nodeNames = append(nodeNames, node.Name)
	}
	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "ReconcileNodeAccess",
			"Type":   "NASStorageDriver",
			"Nodes":  nodeNames,
		}
		Logc(ctx).WithFields(fields).Debug(">>>> ReconcileNodeAccess")
		defer Logc(ctx).WithFields(fields).Debug("<<<< ReconcileNodeAccess")
	}

	policyName := getExportPolicyName(backendUUID)

	return reconcileNASNodeAccess(ctx, nodes, &d.Config, d.API, policyName)
}

// String makes NASStorageDriver satisfy the Stringer interface.
func (d NASStorageDriver) String() string {
	return utils.ToStringRedacted(&d, GetOntapDriverRedactList(), d.GetExternalConfig(context.Background()))
}

// GoString makes NASStorageDriver satisfy the GoStringer interface.
func (d NASStorageDriver) GoString() string {
	return d.String()
}

// GetCommonConfig returns driver's CommonConfig
func (d NASStorageDriver) GetCommonConfig(context.Context) *drivers.CommonStorageDriverConfig {
	return d.Config.CommonStorageDriverConfig
}

// EstablishMirror will create a new snapmirror relationship between a RW and a DP volume that have not previously
// had a relationship
func (d *NASStorageDriver) EstablishMirror(
	ctx context.Context, localVolumeHandle, remoteVolumeHandle, replicationPolicy, replicationSchedule string,
) error {
	// If replication policy in TMR is empty use the backend policy
	if replicationPolicy == "" {
		replicationPolicy = d.GetConfig().ReplicationPolicy
	}

	// Validate replication policy, if it is invalid, use the backend policy
	isAsync, err := validateReplicationPolicy(ctx, replicationPolicy, d.API)
	if err != nil {
		Logc(ctx).Debugf("Replication policy given in TMR %s is invalid, using policy %s from backend.",
			replicationPolicy, d.GetConfig().ReplicationPolicy)
		replicationPolicy = d.GetConfig().ReplicationPolicy
		isAsync, err = validateReplicationPolicy(ctx, replicationPolicy, d.API)
		if err != nil {
			Logc(ctx).Debugf("Replication policy %s in backend should be valid.", replicationPolicy)
		}
	}

	// If replication policy is async type, validate the replication schedule from TMR or use backend schedule
	if isAsync {
		if replicationSchedule != "" {
			if err := validateReplicationSchedule(ctx, replicationSchedule, d.API); err != nil {
				Logc(ctx).Debugf("Replication schedule given in TMR %s is invalid, using schedule %s from backend.",
					replicationSchedule, d.GetConfig().ReplicationSchedule)
				replicationSchedule = d.GetConfig().ReplicationSchedule
			}
		} else {
			replicationSchedule = d.GetConfig().ReplicationSchedule
		}
	} else {
		replicationSchedule = ""
	}

	return establishMirror(ctx, localVolumeHandle, remoteVolumeHandle, replicationPolicy, replicationSchedule, d.API)
}

// ReestablishMirror will attempt to resync a snapmirror relationship,
// if and only if the relationship existed previously
func (d *NASStorageDriver) ReestablishMirror(
	ctx context.Context, localVolumeHandle, remoteVolumeHandle, replicationPolicy, replicationSchedule string,
) error {
	// If replication policy in TMR is empty use the backend policy
	if replicationPolicy == "" {
		replicationPolicy = d.GetConfig().ReplicationPolicy
	}

	// Validate replication policy, if it is invalid, use the backend policy
	isAsync, err := validateReplicationPolicy(ctx, replicationPolicy, d.API)
	if err != nil {
		Logc(ctx).Debugf("Replication policy given in TMR %s is invalid, using policy %s from backend.",
			replicationPolicy, d.GetConfig().ReplicationPolicy)
		replicationPolicy = d.GetConfig().ReplicationPolicy
		isAsync, err = validateReplicationPolicy(ctx, replicationPolicy, d.API)
		if err != nil {
			Logc(ctx).Debugf("Replication policy %s in backend should be valid.", replicationPolicy)
		}
	}

	// If replication policy is async type, validate the replication schedule from TMR or use backend schedule
	if isAsync {
		if replicationSchedule != "" {
			if err := validateReplicationSchedule(ctx, replicationSchedule, d.API); err != nil {
				Logc(ctx).Debugf("Replication schedule given in TMR %s is invalid, using schedule %s from backend.",
					replicationSchedule, d.GetConfig().ReplicationSchedule)
				replicationSchedule = d.GetConfig().ReplicationSchedule
			}
		} else {
			replicationSchedule = d.GetConfig().ReplicationSchedule
		}
	} else {
		replicationSchedule = ""
	}

	return reestablishMirror(ctx, localVolumeHandle, remoteVolumeHandle, replicationPolicy, replicationSchedule, d.API)
}

// PromoteMirror will break the snapmirror and make the destination volume RW,
// optionally after a given snapshot has synced
func (d *NASStorageDriver) PromoteMirror(
	ctx context.Context, localVolumeHandle, remoteVolumeHandle, snapshotName string,
) (bool, error) {
	return promoteMirror(ctx, localVolumeHandle, remoteVolumeHandle, snapshotName, d.GetConfig().ReplicationPolicy,
		d.API)
}

// GetMirrorStatus returns the current state of a snapmirror relationship
func (d *NASStorageDriver) GetMirrorStatus(
	ctx context.Context, localVolumeHandle, remoteVolumeHandle string,
) (string, error) {
	return getMirrorStatus(ctx, localVolumeHandle, remoteVolumeHandle, d.API)
}

// ReleaseMirror will release the snapmirror relationship data of the source volume
func (d *NASStorageDriver) ReleaseMirror(ctx context.Context, localVolumeHandle string) error {
	return releaseMirror(ctx, localVolumeHandle, d.API)
}

// GetReplicationDetails returns the replication policy and schedule of a snapmirror relationship
func (d *NASStorageDriver) GetReplicationDetails(
	ctx context.Context, localVolumeHandle, remoteVolumeHandle string,
) (string, string, error) {
	return getReplicationDetails(ctx, localVolumeHandle, remoteVolumeHandle, d.API)
}

// MountVolume returns the volume mount error(if any)
func (d *NASStorageDriver) MountVolume(
	ctx context.Context, name, junctionPath string, flexVol *api.Volume,
) error {
	if err := d.API.VolumeMount(ctx, name, junctionPath); err != nil {
		// An API error is returned if we attempt to mount a DP volume that has not yet been snapmirrored,
		// we expect this to be the case.

		if api.IsApiError(err) && flexVol.DPVolume {
			Logc(ctx).Debugf("Received expected API error when mounting DP volume to junction; %v", err)
		} else {
			return fmt.Errorf("error mounting volume to junction %s; %v", junctionPath, err)
		}
	}
	return nil
}
