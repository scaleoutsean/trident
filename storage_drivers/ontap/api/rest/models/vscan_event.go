// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// VscanEvent vscan event
//
// swagger:model vscan_event
type VscanEvent struct {

	// Specifies the reason of the Vscan server disconnection.
	// The available values are:
	// * na                        Not applicable
	// * vscan_disabled            Vscan disabled on the SVM
	// * no_data_lif               SVM does not have data lif on the node
	// * session_uninitialized     Session not initialized
	// * remote_closed             Closure from Server
	// * invalid_protocol_msg      Invalid protocol-message received
	// * invalid_session_id        Invalid session-id received
	// * inactive_connection       No activity on connection
	// * invalid_user              Connection request by invalid user
	// * server_removed            Server removed from the active scanner-pool
	//
	// Read Only: true
	// Enum: [na vscan_disabled no_data_lif session_uninitialized remote_closed invalid_protocol_msg invalid_session_id inactive_connection invalid_user server_removed]
	DisconnectReason string `json:"disconnect_reason,omitempty"`

	// Specifies the Timestamp of the event.
	// Example: 2021-11-25T04:29:41.606Z
	// Format: date-time
	EventTime *strfmt.DateTime `json:"event_time,omitempty"`

	// Specifies the file for which event happened.
	// Example: /1
	FilePath string `json:"file_path,omitempty"`

	// interface
	Interface *VscanEventInterface `json:"interface,omitempty"`

	// node
	Node *VscanEventNode `json:"node,omitempty"`

	// Specifies the IP address of the Vscan server.
	// Example: 192.168.1.1
	Server string `json:"server,omitempty"`

	// svm
	Svm *VscanEventSvm `json:"svm,omitempty"`

	// Specifies the event type.
	// Enum: [scanner_connected scanner_disconnected scanner_updated scan_internal_error scan_failed scan_timedout file_infected file_renamed file_quarantined file_deleted scanner_busy]
	Type string `json:"type,omitempty"`

	// Specifies the scan-engine vendor.
	// Example: mighty master anti-evil scanner
	Vendor string `json:"vendor,omitempty"`

	// Specifies the scan-engine version.
	// Example: 1.0
	Version string `json:"version,omitempty"`
}

// Validate validates this vscan event
func (m *VscanEvent) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDisconnectReason(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEventTime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInterface(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSvm(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var vscanEventTypeDisconnectReasonPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["na","vscan_disabled","no_data_lif","session_uninitialized","remote_closed","invalid_protocol_msg","invalid_session_id","inactive_connection","invalid_user","server_removed"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		vscanEventTypeDisconnectReasonPropEnum = append(vscanEventTypeDisconnectReasonPropEnum, v)
	}
}

const (

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// na
	// END DEBUGGING
	// VscanEventDisconnectReasonNa captures enum value "na"
	VscanEventDisconnectReasonNa string = "na"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// vscan_disabled
	// END DEBUGGING
	// VscanEventDisconnectReasonVscanDisabled captures enum value "vscan_disabled"
	VscanEventDisconnectReasonVscanDisabled string = "vscan_disabled"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// no_data_lif
	// END DEBUGGING
	// VscanEventDisconnectReasonNoDataLif captures enum value "no_data_lif"
	VscanEventDisconnectReasonNoDataLif string = "no_data_lif"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// session_uninitialized
	// END DEBUGGING
	// VscanEventDisconnectReasonSessionUninitialized captures enum value "session_uninitialized"
	VscanEventDisconnectReasonSessionUninitialized string = "session_uninitialized"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// remote_closed
	// END DEBUGGING
	// VscanEventDisconnectReasonRemoteClosed captures enum value "remote_closed"
	VscanEventDisconnectReasonRemoteClosed string = "remote_closed"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// invalid_protocol_msg
	// END DEBUGGING
	// VscanEventDisconnectReasonInvalidProtocolMsg captures enum value "invalid_protocol_msg"
	VscanEventDisconnectReasonInvalidProtocolMsg string = "invalid_protocol_msg"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// invalid_session_id
	// END DEBUGGING
	// VscanEventDisconnectReasonInvalidSessionID captures enum value "invalid_session_id"
	VscanEventDisconnectReasonInvalidSessionID string = "invalid_session_id"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// inactive_connection
	// END DEBUGGING
	// VscanEventDisconnectReasonInactiveConnection captures enum value "inactive_connection"
	VscanEventDisconnectReasonInactiveConnection string = "inactive_connection"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// invalid_user
	// END DEBUGGING
	// VscanEventDisconnectReasonInvalidUser captures enum value "invalid_user"
	VscanEventDisconnectReasonInvalidUser string = "invalid_user"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// disconnect_reason
	// DisconnectReason
	// server_removed
	// END DEBUGGING
	// VscanEventDisconnectReasonServerRemoved captures enum value "server_removed"
	VscanEventDisconnectReasonServerRemoved string = "server_removed"
)

// prop value enum
func (m *VscanEvent) validateDisconnectReasonEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, vscanEventTypeDisconnectReasonPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *VscanEvent) validateDisconnectReason(formats strfmt.Registry) error {
	if swag.IsZero(m.DisconnectReason) { // not required
		return nil
	}

	// value enum
	if err := m.validateDisconnectReasonEnum("disconnect_reason", "body", m.DisconnectReason); err != nil {
		return err
	}

	return nil
}

func (m *VscanEvent) validateEventTime(formats strfmt.Registry) error {
	if swag.IsZero(m.EventTime) { // not required
		return nil
	}

	if err := validate.FormatOf("event_time", "body", "date-time", m.EventTime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *VscanEvent) validateInterface(formats strfmt.Registry) error {
	if swag.IsZero(m.Interface) { // not required
		return nil
	}

	if m.Interface != nil {
		if err := m.Interface.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface")
			}
			return err
		}
	}

	return nil
}

func (m *VscanEvent) validateNode(formats strfmt.Registry) error {
	if swag.IsZero(m.Node) { // not required
		return nil
	}

	if m.Node != nil {
		if err := m.Node.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("node")
			}
			return err
		}
	}

	return nil
}

func (m *VscanEvent) validateSvm(formats strfmt.Registry) error {
	if swag.IsZero(m.Svm) { // not required
		return nil
	}

	if m.Svm != nil {
		if err := m.Svm.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("svm")
			}
			return err
		}
	}

	return nil
}

var vscanEventTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["scanner_connected","scanner_disconnected","scanner_updated","scan_internal_error","scan_failed","scan_timedout","file_infected","file_renamed","file_quarantined","file_deleted","scanner_busy"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		vscanEventTypeTypePropEnum = append(vscanEventTypeTypePropEnum, v)
	}
}

const (

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// scanner_connected
	// END DEBUGGING
	// VscanEventTypeScannerConnected captures enum value "scanner_connected"
	VscanEventTypeScannerConnected string = "scanner_connected"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// scanner_disconnected
	// END DEBUGGING
	// VscanEventTypeScannerDisconnected captures enum value "scanner_disconnected"
	VscanEventTypeScannerDisconnected string = "scanner_disconnected"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// scanner_updated
	// END DEBUGGING
	// VscanEventTypeScannerUpdated captures enum value "scanner_updated"
	VscanEventTypeScannerUpdated string = "scanner_updated"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// scan_internal_error
	// END DEBUGGING
	// VscanEventTypeScanInternalError captures enum value "scan_internal_error"
	VscanEventTypeScanInternalError string = "scan_internal_error"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// scan_failed
	// END DEBUGGING
	// VscanEventTypeScanFailed captures enum value "scan_failed"
	VscanEventTypeScanFailed string = "scan_failed"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// scan_timedout
	// END DEBUGGING
	// VscanEventTypeScanTimedout captures enum value "scan_timedout"
	VscanEventTypeScanTimedout string = "scan_timedout"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// file_infected
	// END DEBUGGING
	// VscanEventTypeFileInfected captures enum value "file_infected"
	VscanEventTypeFileInfected string = "file_infected"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// file_renamed
	// END DEBUGGING
	// VscanEventTypeFileRenamed captures enum value "file_renamed"
	VscanEventTypeFileRenamed string = "file_renamed"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// file_quarantined
	// END DEBUGGING
	// VscanEventTypeFileQuarantined captures enum value "file_quarantined"
	VscanEventTypeFileQuarantined string = "file_quarantined"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// file_deleted
	// END DEBUGGING
	// VscanEventTypeFileDeleted captures enum value "file_deleted"
	VscanEventTypeFileDeleted string = "file_deleted"

	// BEGIN DEBUGGING
	// vscan_event
	// VscanEvent
	// type
	// Type
	// scanner_busy
	// END DEBUGGING
	// VscanEventTypeScannerBusy captures enum value "scanner_busy"
	VscanEventTypeScannerBusy string = "scanner_busy"
)

// prop value enum
func (m *VscanEvent) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, vscanEventTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *VscanEvent) validateType(formats strfmt.Registry) error {
	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this vscan event based on the context it is used
func (m *VscanEvent) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDisconnectReason(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateInterface(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNode(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSvm(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEvent) contextValidateDisconnectReason(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "disconnect_reason", "body", string(m.DisconnectReason)); err != nil {
		return err
	}

	return nil
}

func (m *VscanEvent) contextValidateInterface(ctx context.Context, formats strfmt.Registry) error {

	if m.Interface != nil {
		if err := m.Interface.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface")
			}
			return err
		}
	}

	return nil
}

func (m *VscanEvent) contextValidateNode(ctx context.Context, formats strfmt.Registry) error {

	if m.Node != nil {
		if err := m.Node.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("node")
			}
			return err
		}
	}

	return nil
}

func (m *VscanEvent) contextValidateSvm(ctx context.Context, formats strfmt.Registry) error {

	if m.Svm != nil {
		if err := m.Svm.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("svm")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VscanEvent) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VscanEvent) UnmarshalBinary(b []byte) error {
	var res VscanEvent
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// VscanEventInterface Address of the interface used for the Vscan connection.
//
// swagger:model VscanEventInterface
type VscanEventInterface struct {

	// links
	Links *VscanEventInterfaceLinks `json:"_links,omitempty"`

	// ip
	IP *VscanEventInterfaceIP `json:"ip,omitempty"`

	// The name of the interface. If only the name is provided, the SVM scope
	// must be provided by the object this object is embedded in.
	//
	// Example: lif1
	Name string `json:"name,omitempty"`

	// The UUID that uniquely identifies the interface.
	// Example: 1cd8a442-86d1-11e0-ae1c-123478563412
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this vscan event interface
func (m *VscanEventInterface) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIP(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventInterface) validateLinks(formats strfmt.Registry) error {
	if swag.IsZero(m.Links) { // not required
		return nil
	}

	if m.Links != nil {
		if err := m.Links.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

func (m *VscanEventInterface) validateIP(formats strfmt.Registry) error {
	if swag.IsZero(m.IP) { // not required
		return nil
	}

	if m.IP != nil {
		if err := m.IP.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface" + "." + "ip")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this vscan event interface based on the context it is used
func (m *VscanEventInterface) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateIP(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventInterface) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	if m.Links != nil {
		if err := m.Links.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

func (m *VscanEventInterface) contextValidateIP(ctx context.Context, formats strfmt.Registry) error {

	if m.IP != nil {
		if err := m.IP.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface" + "." + "ip")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VscanEventInterface) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VscanEventInterface) UnmarshalBinary(b []byte) error {
	var res VscanEventInterface
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// VscanEventInterfaceIP IP information
//
// swagger:model VscanEventInterfaceIP
type VscanEventInterfaceIP struct {

	// address
	Address IPAddressReadonly `json:"address,omitempty"`
}

// Validate validates this vscan event interface IP
func (m *VscanEventInterfaceIP) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAddress(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventInterfaceIP) validateAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.Address) { // not required
		return nil
	}

	if err := m.Address.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("interface" + "." + "ip" + "." + "address")
		}
		return err
	}

	return nil
}

// ContextValidate validate this vscan event interface IP based on the context it is used
func (m *VscanEventInterfaceIP) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventInterfaceIP) contextValidateAddress(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Address.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("interface" + "." + "ip" + "." + "address")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VscanEventInterfaceIP) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VscanEventInterfaceIP) UnmarshalBinary(b []byte) error {
	var res VscanEventInterfaceIP
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// VscanEventInterfaceLinks vscan event interface links
//
// swagger:model VscanEventInterfaceLinks
type VscanEventInterfaceLinks struct {

	// self
	Self *Href `json:"self,omitempty"`
}

// Validate validates this vscan event interface links
func (m *VscanEventInterfaceLinks) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSelf(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventInterfaceLinks) validateSelf(formats strfmt.Registry) error {
	if swag.IsZero(m.Self) { // not required
		return nil
	}

	if m.Self != nil {
		if err := m.Self.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this vscan event interface links based on the context it is used
func (m *VscanEventInterfaceLinks) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSelf(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventInterfaceLinks) contextValidateSelf(ctx context.Context, formats strfmt.Registry) error {

	if m.Self != nil {
		if err := m.Self.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VscanEventInterfaceLinks) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VscanEventInterfaceLinks) UnmarshalBinary(b []byte) error {
	var res VscanEventInterfaceLinks
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// VscanEventNode vscan event node
//
// swagger:model VscanEventNode
type VscanEventNode struct {

	// links
	Links *VscanEventNodeLinks `json:"_links,omitempty"`

	// name
	// Example: node1
	Name string `json:"name,omitempty"`

	// uuid
	// Example: 1cd8a442-86d1-11e0-ae1c-123478563412
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this vscan event node
func (m *VscanEventNode) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventNode) validateLinks(formats strfmt.Registry) error {
	if swag.IsZero(m.Links) { // not required
		return nil
	}

	if m.Links != nil {
		if err := m.Links.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("node" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this vscan event node based on the context it is used
func (m *VscanEventNode) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventNode) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	if m.Links != nil {
		if err := m.Links.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("node" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VscanEventNode) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VscanEventNode) UnmarshalBinary(b []byte) error {
	var res VscanEventNode
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// VscanEventNodeLinks vscan event node links
//
// swagger:model VscanEventNodeLinks
type VscanEventNodeLinks struct {

	// self
	Self *Href `json:"self,omitempty"`
}

// Validate validates this vscan event node links
func (m *VscanEventNodeLinks) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSelf(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventNodeLinks) validateSelf(formats strfmt.Registry) error {
	if swag.IsZero(m.Self) { // not required
		return nil
	}

	if m.Self != nil {
		if err := m.Self.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("node" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this vscan event node links based on the context it is used
func (m *VscanEventNodeLinks) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSelf(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventNodeLinks) contextValidateSelf(ctx context.Context, formats strfmt.Registry) error {

	if m.Self != nil {
		if err := m.Self.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("node" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VscanEventNodeLinks) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VscanEventNodeLinks) UnmarshalBinary(b []byte) error {
	var res VscanEventNodeLinks
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// VscanEventSvm vscan event svm
//
// swagger:model VscanEventSvm
type VscanEventSvm struct {

	// links
	Links *VscanEventSvmLinks `json:"_links,omitempty"`

	// The name of the SVM.
	//
	// Example: svm1
	Name string `json:"name,omitempty"`

	// The unique identifier of the SVM.
	//
	// Example: 02c9e252-41be-11e9-81d5-00a0986138f7
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this vscan event svm
func (m *VscanEventSvm) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventSvm) validateLinks(formats strfmt.Registry) error {
	if swag.IsZero(m.Links) { // not required
		return nil
	}

	if m.Links != nil {
		if err := m.Links.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("svm" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this vscan event svm based on the context it is used
func (m *VscanEventSvm) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventSvm) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	if m.Links != nil {
		if err := m.Links.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("svm" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VscanEventSvm) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VscanEventSvm) UnmarshalBinary(b []byte) error {
	var res VscanEventSvm
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// VscanEventSvmLinks vscan event svm links
//
// swagger:model VscanEventSvmLinks
type VscanEventSvmLinks struct {

	// self
	Self *Href `json:"self,omitempty"`
}

// Validate validates this vscan event svm links
func (m *VscanEventSvmLinks) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSelf(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventSvmLinks) validateSelf(formats strfmt.Registry) error {
	if swag.IsZero(m.Self) { // not required
		return nil
	}

	if m.Self != nil {
		if err := m.Self.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("svm" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this vscan event svm links based on the context it is used
func (m *VscanEventSvmLinks) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSelf(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VscanEventSvmLinks) contextValidateSelf(ctx context.Context, formats strfmt.Registry) error {

	if m.Self != nil {
		if err := m.Self.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("svm" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VscanEventSvmLinks) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VscanEventSvmLinks) UnmarshalBinary(b []byte) error {
	var res VscanEventSvmLinks
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
