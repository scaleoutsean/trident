// Copyright 2022 NetApp, Inc. All Rights Reserved.

package api

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/netapp/trident/storage_drivers/ontap/api/azgo"

	. "github.com/netapp/trident/logger"
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

// //////////////////////////////////////////////////////////////////////////////////////////////////////
// Abstraction layer
// //////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	failureLUNCreate  = "failure_65dc2f4b_adbe_4ed3_8b73_6c61d5eac054"
	failureLUNSetAttr = "failure_7c3a89e2_7d83_457b_9e29_bfdb082c1d8b"
)

type OntapAPI interface {
	APIVersion(context.Context) (string, error)
	SVMName() string

	EmsAutosupportLog(
		ctx context.Context, driverName, appVersion string, autoSupport bool, category string,
		computerName, eventDescription string, eventID int, eventSource string, logLevel int,
	)

	ExportPolicyCreate(ctx context.Context, policy string) error
	ExportPolicyDestroy(ctx context.Context, policy string) error
	ExportPolicyExists(ctx context.Context, policyName string) (bool, error)
	ExportRuleCreate(ctx context.Context, policyName, desiredPolicyRule, nasProtocol string) error
	ExportRuleDestroy(ctx context.Context, policyName string, ruleIndex int) error
	ExportRuleList(ctx context.Context, policyName string) (map[string]int, error)

	FlexgroupCreate(ctx context.Context, volume Volume) error
	FlexgroupExists(ctx context.Context, volumeName string) (bool, error)
	FlexgroupInfo(ctx context.Context, volumeName string) (*Volume, error)
	FlexgroupDisableSnapshotDirectoryAccess(ctx context.Context, volumeName string) error
	FlexgroupSetComment(ctx context.Context, volumeNameInternal, volumeNameExternal, comment string) error
	FlexgroupModifyUnixPermissions(
		ctx context.Context, volumeNameInternal, volumeNameExternal, unixPermissions string,
	) error
	FlexgroupMount(ctx context.Context, name, junctionPath string) error
	FlexgroupListByPrefix(ctx context.Context, prefix string) (Volumes, error)
	FlexgroupDestroy(ctx context.Context, volumeName string, force bool) error
	FlexgroupSetSize(ctx context.Context, name, newSize string) error
	FlexgroupSize(ctx context.Context, volumeName string) (uint64, error)
	FlexgroupUnmount(ctx context.Context, name string, force bool) error
	FlexgroupUsedSize(ctx context.Context, volumeName string) (int, error)
	FlexgroupModifyExportPolicy(ctx context.Context, volumeName, policyName string) error
	FlexgroupSnapshotCreate(ctx context.Context, snapshotName, sourceVolume string) error
	FlexgroupSetQosPolicyGroupName(ctx context.Context, name string, qos QosPolicyGroup) error
	FlexgroupCloneSplitStart(ctx context.Context, cloneName string) error
	FlexgroupSnapshotList(ctx context.Context, sourceVolume string) (Snapshots, error)
	FlexgroupSnapshotDelete(ctx context.Context, snapshotName, sourceVolume string) error

	LunList(ctx context.Context, pattern string) (Luns, error)
	LunCreate(ctx context.Context, lun Lun) error
	LunCloneCreate(ctx context.Context, flexvol, source, lunName string, qosPolicyGroup QosPolicyGroup) error
	LunDestroy(ctx context.Context, lunPath string) error
	LunGetGeometry(ctx context.Context, lunPath string) (uint64, error)
	LunGetComment(ctx context.Context, lunPath string) (string, bool, error)
	LunSetAttribute(ctx context.Context, lunPath, attribute, fstype, context, luks string) error
	ParseLunComment(ctx context.Context, commentJSON string) (map[string]string, error)
	LunSetQosPolicyGroup(ctx context.Context, lunPath string, qosPolicyGroup QosPolicyGroup) error
	LunGetByName(ctx context.Context, name string) (*Lun, error)
	LunRename(ctx context.Context, lunPath, newLunPath string) error
	LunMapInfo(ctx context.Context, initiatorGroupName, lunPath string) (int, error)
	EnsureLunMapped(ctx context.Context, initiatorGroupName, lunPath string, importNotManaged bool) (int, error)
	LunUnmap(ctx context.Context, initiatorGroupName, lunPath string) error
	LunSize(ctx context.Context, lunPath string) (int, error)
	LunSetSize(ctx context.Context, lunPath, newSize string) (uint64, error)
	LunMapGetReportingNodes(ctx context.Context, initiatorGroupName, lunPath string) ([]string, error)
	LunListIgroupsMapped(ctx context.Context, lunPath string) ([]string, error)

	IscsiInitiatorGetDefaultAuth(ctx context.Context) (IscsiInitiatorAuth, error)
	IscsiInitiatorSetDefaultAuth(
		ctx context.Context, authType, userName, passphrase, outbountUserName,
		outboundPassphrase string,
	) error
	IscsiInterfaceGet(ctx context.Context, svm string) ([]string, error)
	IscsiNodeGetNameRequest(ctx context.Context) (string, error)

	IgroupCreate(ctx context.Context, initiatorGroupName, initiatorGroupType, osType string) error
	IgroupDestroy(ctx context.Context, initiatorGroupName string) error
	EnsureIgroupAdded(ctx context.Context, initiatorGroupName, initiator string) error
	IgroupRemove(ctx context.Context, initiatorGroupName, initiator string, force bool) error
	IgroupGetByName(ctx context.Context, initiatorGroupName string) (map[string]bool, error)
	IgroupListLUNsMapped(ctx context.Context, initiatorGroupName string) ([]string, error)

	GetSVMAggregateAttributes(ctx context.Context) (map[string]string, error)
	GetSVMAggregateNames(ctx context.Context) ([]string, error)
	GetSVMAggregateSpace(ctx context.Context, aggregate string) ([]SVMAggregateSpace, error)
	GetSVMPeers(ctx context.Context) ([]string, error)

	GetSVMUUID() string

	QtreeExists(ctx context.Context, name, volumePattern string) (bool, string, error)
	QtreeCreate(
		ctx context.Context, name, volumeName, unixPermissions, exportPolicy, securityStyle, qosPolicy string,
	) error
	QtreeDestroyAsync(ctx context.Context, path string, force bool) error
	QtreeRename(ctx context.Context, path, newPath string) error
	QtreeModifyExportPolicy(ctx context.Context, name, volumeName, newExportPolicyName string) error
	QtreeCount(ctx context.Context, volumeName string) (int, error)
	QtreeListByPrefix(ctx context.Context, prefix, volumePrefix string) (Qtrees, error)
	QtreeGetByName(ctx context.Context, name, volumePrefix string) (*Qtree, error)

	QuotaEntryList(ctx context.Context, volumeName string) (QuotaEntries, error)
	QuotaOff(ctx context.Context, volumeName string) error
	QuotaOn(ctx context.Context, volumeName string) error
	QuotaResize(ctx context.Context, volumeName string) error
	QuotaStatus(ctx context.Context, volumeName string) (string, error)
	QuotaSetEntry(ctx context.Context, qtreeName, volumeName, quotaType, diskLimit string) error
	QuotaGetEntry(ctx context.Context, volumeName, qtreeName, quotaType string) (*QuotaEntry, error)

	GetSLMDataLifs(ctx context.Context, ips, reportingNodeNames []string) ([]string, error)
	NetInterfaceGetDataLIFs(ctx context.Context, protocol string) ([]string, error)
	NodeListSerialNumbers(ctx context.Context) ([]string, error)

	SnapshotRestoreVolume(ctx context.Context, snapshotName, sourceVolume string) error
	SnapshotRestoreFlexgroup(ctx context.Context, snapshotName, sourceVolume string) error

	SnapmirrorCreate(
		ctx context.Context, localFlexvolName, localSVMName, remoteFlexvolName, remoteSVMName,
		replicationPolicy, replicationSchedule string,
	) error
	SnapmirrorResync(ctx context.Context, localFlexvolName, localSVMName, remoteFlexvolName, remoteSVMName string) error
	SnapmirrorDelete(ctx context.Context, localFlexvolName, localSVMName, remoteFlexvolName, remoteSVMName string) error
	SnapmirrorGet(ctx context.Context, localFlexvolName, localSVMName, remoteFlexvolName, remoteSVMName string) (
		*Snapmirror, error)
	SnapmirrorInitialize(
		ctx context.Context, localFlexvolName, localSVMName, remoteFlexvolName, remoteSVMName string,
	) error
	SnapmirrorPolicyGet(ctx context.Context, replicationPolicy string) (*SnapmirrorPolicy, error)
	SnapmirrorQuiesce(
		ctx context.Context, localFlexvolName, localSVMName, remoteFlexvolName, remoteSVMName string,
	) error
	SnapmirrorAbort(ctx context.Context, localFlexvolName, localSVMName, remoteFlexvolName, remoteSVMName string) error
	SnapmirrorBreak(
		ctx context.Context, localFlexvolName, localSVMName, remoteFlexvolName, remoteSVMName,
		snapshotName string,
	) error
	JobScheduleExists(ctx context.Context, replicationSchedule string) error
	SupportsFeature(ctx context.Context, feature Feature) bool
	ValidateAPIVersion(ctx context.Context) error
	SnapmirrorDeleteViaDestination(localFlexvolName, localSVMName string) error
	SnapmirrorRelease(localFlexvolName, localSVMName string) error
	IsSVMDRCapable(ctx context.Context) (bool, error)

	VolumeCloneCreate(ctx context.Context, cloneName, sourceName, snapshot string, async bool) error
	VolumeCloneSplitStart(ctx context.Context, cloneName string) error

	VolumeCreate(ctx context.Context, volume Volume) error
	VolumeDestroy(ctx context.Context, volumeName string, force bool) error
	VolumeDisableSnapshotDirectoryAccess(ctx context.Context, name string) error
	VolumeExists(ctx context.Context, volumeName string) (bool, error)
	VolumeInfo(ctx context.Context, volumeName string) (*Volume, error)
	VolumeListByPrefix(ctx context.Context, prefix string) (Volumes, error)
	VolumeListBySnapshotParent(ctx context.Context, snapshotName, sourceVolume string) (VolumeNameList, error)
	VolumeModifyExportPolicy(ctx context.Context, volumeName, policyName string) error
	VolumeModifyUnixPermissions(
		ctx context.Context, volumeNameInternal, volumeNameExternal, unixPermissions string,
	) error
	VolumeMount(ctx context.Context, name, junctionPath string) error
	VolumeListByAttrs(ctx context.Context, volumeAttrs *Volume) (Volumes, error)
	VolumeRename(ctx context.Context, originalName, newName string) error
	VolumeSetComment(ctx context.Context, volumeNameInternal, volumeNameExternal, comment string) error
	VolumeSetQosPolicyGroupName(ctx context.Context, name string, qos QosPolicyGroup) error
	VolumeSetSize(ctx context.Context, name, newSize string) error
	VolumeSize(ctx context.Context, volumeName string) (uint64, error)
	VolumeUsedSize(ctx context.Context, volumeName string) (int, error)
	VolumeSnapshotCreate(ctx context.Context, snapshotName, sourceVolume string) error
	VolumeSnapshotList(ctx context.Context, sourceVolume string) (Snapshots, error)
	VolumeSnapshotDelete(ctx context.Context, snapshotName, sourceVolume string) error

	TieringPolicyValue(ctx context.Context) string
}

type AggregateSpace interface {
	Size() int64
	Used() int64
	Footprint() int64
}

type SVMAggregateSpace struct {
	size      int64
	used      int64
	footprint int64
}

func NewSVMAggregateSpace(size, used, footprint int64) SVMAggregateSpace {
	return SVMAggregateSpace{
		size:      size,
		used:      used,
		footprint: footprint,
	}
}

func (o SVMAggregateSpace) Size() int64 {
	return o.size
}

func (o SVMAggregateSpace) Used() int64 {
	return o.used
}

func (o SVMAggregateSpace) Footprint() int64 {
	return o.footprint
}

type Response interface {
	APIName() string
	Client() string
	Name() string
	Version() string
	Status() string
	Reason() string
	Errno() string
}

type APIResponse struct {
	apiName string
	status  string
	reason  string
	errno   string
}

// NewAPIResponse factory method to create a new instance of an APIResponse
func NewAPIResponse(status, reason, errno string) *APIResponse {
	result := &APIResponse{
		status: status,
		reason: reason,
		errno:  errno,
	}
	return result
}

func (o APIResponse) APIName() string {
	return o.apiName
}

func (o APIResponse) Status() string {
	return o.status
}

func (o APIResponse) Reason() string {
	return o.reason
}

func (o APIResponse) Errno() string {
	return o.errno
}

// GetError inspects the supplied *APIResponse and error parameters to determine if an error occurred
func GetError(ctx context.Context, response *APIResponse, errorIn error) (errorOut error) {
	defer func() {
		if r := recover(); r != nil {
			Logc(ctx).Errorf("Panic in ontap#GetError. %v\nStack Trace: %v", response, string(debug.Stack()))
			errorOut = azgo.ZapiError{}
		}
	}()

	if errorIn != nil {
		errorOut = errorIn
		return errorOut
	}

	if response == nil {
		errorOut = fmt.Errorf("API error: nil response")
		return errorOut
	}

	responseStatus := response.Status()
	if strings.EqualFold(responseStatus, "passed") {
		errorOut = nil
		return errorOut
	}

	errorOut = fmt.Errorf("API error: %v", response)
	return errorOut
}
