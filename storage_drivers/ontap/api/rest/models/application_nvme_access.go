// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ApplicationNvmeAccess Application NVME access
//
// swagger:model application_nvme_access
type ApplicationNvmeAccess struct {

	// backing storage
	BackingStorage *ApplicationNvmeAccessBackingStorage `json:"backing_storage,omitempty"`

	// Clone
	// Read Only: true
	IsClone *bool `json:"is_clone,omitempty"`

	// subsystem map
	SubsystemMap *ApplicationNvmeAccessSubsystemMap `json:"subsystem_map,omitempty"`
}

// Validate validates this application nvme access
func (m *ApplicationNvmeAccess) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBackingStorage(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubsystemMap(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccess) validateBackingStorage(formats strfmt.Registry) error {
	if swag.IsZero(m.BackingStorage) { // not required
		return nil
	}

	if m.BackingStorage != nil {
		if err := m.BackingStorage.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("backing_storage")
			}
			return err
		}
	}

	return nil
}

func (m *ApplicationNvmeAccess) validateSubsystemMap(formats strfmt.Registry) error {
	if swag.IsZero(m.SubsystemMap) { // not required
		return nil
	}

	if m.SubsystemMap != nil {
		if err := m.SubsystemMap.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subsystem_map")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this application nvme access based on the context it is used
func (m *ApplicationNvmeAccess) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBackingStorage(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateIsClone(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSubsystemMap(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccess) contextValidateBackingStorage(ctx context.Context, formats strfmt.Registry) error {

	if m.BackingStorage != nil {
		if err := m.BackingStorage.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("backing_storage")
			}
			return err
		}
	}

	return nil
}

func (m *ApplicationNvmeAccess) contextValidateIsClone(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "is_clone", "body", m.IsClone); err != nil {
		return err
	}

	return nil
}

func (m *ApplicationNvmeAccess) contextValidateSubsystemMap(ctx context.Context, formats strfmt.Registry) error {

	if m.SubsystemMap != nil {
		if err := m.SubsystemMap.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subsystem_map")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationNvmeAccess) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationNvmeAccess) UnmarshalBinary(b []byte) error {
	var res ApplicationNvmeAccess
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ApplicationNvmeAccessBackingStorage application nvme access backing storage
//
// swagger:model ApplicationNvmeAccessBackingStorage
type ApplicationNvmeAccessBackingStorage struct {

	// Backing storage type
	// Read Only: true
	// Enum: [namespace]
	Type string `json:"type,omitempty"`

	// Backing storage UUID
	// Read Only: true
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this application nvme access backing storage
func (m *ApplicationNvmeAccessBackingStorage) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var applicationNvmeAccessBackingStorageTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["namespace"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		applicationNvmeAccessBackingStorageTypeTypePropEnum = append(applicationNvmeAccessBackingStorageTypeTypePropEnum, v)
	}
}

const (

	// BEGIN DEBUGGING
	// ApplicationNvmeAccessBackingStorage
	// ApplicationNvmeAccessBackingStorage
	// type
	// Type
	// namespace
	// END DEBUGGING
	// ApplicationNvmeAccessBackingStorageTypeNamespace captures enum value "namespace"
	ApplicationNvmeAccessBackingStorageTypeNamespace string = "namespace"
)

// prop value enum
func (m *ApplicationNvmeAccessBackingStorage) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, applicationNvmeAccessBackingStorageTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *ApplicationNvmeAccessBackingStorage) validateType(formats strfmt.Registry) error {
	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("backing_storage"+"."+"type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this application nvme access backing storage based on the context it is used
func (m *ApplicationNvmeAccessBackingStorage) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUUID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessBackingStorage) contextValidateType(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "backing_storage"+"."+"type", "body", string(m.Type)); err != nil {
		return err
	}

	return nil
}

func (m *ApplicationNvmeAccessBackingStorage) contextValidateUUID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "backing_storage"+"."+"uuid", "body", string(m.UUID)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationNvmeAccessBackingStorage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationNvmeAccessBackingStorage) UnmarshalBinary(b []byte) error {
	var res ApplicationNvmeAccessBackingStorage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ApplicationNvmeAccessSubsystemMap Subsystem map object
//
// swagger:model ApplicationNvmeAccessSubsystemMap
type ApplicationNvmeAccessSubsystemMap struct {

	// Subsystem ANA group ID
	// Read Only: true
	Anagrpid string `json:"anagrpid,omitempty"`

	// Subsystem namespace ID
	// Read Only: true
	Nsid string `json:"nsid,omitempty"`

	// subsystem
	Subsystem *ApplicationNvmeAccessSubsystemMapSubsystem `json:"subsystem,omitempty"`
}

// Validate validates this application nvme access subsystem map
func (m *ApplicationNvmeAccessSubsystemMap) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSubsystem(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMap) validateSubsystem(formats strfmt.Registry) error {
	if swag.IsZero(m.Subsystem) { // not required
		return nil
	}

	if m.Subsystem != nil {
		if err := m.Subsystem.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subsystem_map" + "." + "subsystem")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this application nvme access subsystem map based on the context it is used
func (m *ApplicationNvmeAccessSubsystemMap) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAnagrpid(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNsid(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSubsystem(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMap) contextValidateAnagrpid(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "subsystem_map"+"."+"anagrpid", "body", string(m.Anagrpid)); err != nil {
		return err
	}

	return nil
}

func (m *ApplicationNvmeAccessSubsystemMap) contextValidateNsid(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "subsystem_map"+"."+"nsid", "body", string(m.Nsid)); err != nil {
		return err
	}

	return nil
}

func (m *ApplicationNvmeAccessSubsystemMap) contextValidateSubsystem(ctx context.Context, formats strfmt.Registry) error {

	if m.Subsystem != nil {
		if err := m.Subsystem.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subsystem_map" + "." + "subsystem")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMap) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMap) UnmarshalBinary(b []byte) error {
	var res ApplicationNvmeAccessSubsystemMap
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ApplicationNvmeAccessSubsystemMapSubsystem application nvme access subsystem map subsystem
//
// swagger:model ApplicationNvmeAccessSubsystemMapSubsystem
type ApplicationNvmeAccessSubsystemMapSubsystem struct {

	// links
	Links *ApplicationNvmeAccessSubsystemMapSubsystemLinks `json:"_links,omitempty"`

	// hosts
	// Read Only: true
	Hosts []*ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0 `json:"hosts,omitempty"`

	// Subsystem name
	// Read Only: true
	Name string `json:"name,omitempty"`

	// Subsystem UUID
	// Read Only: true
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this application nvme access subsystem map subsystem
func (m *ApplicationNvmeAccessSubsystemMapSubsystem) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHosts(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystem) validateLinks(formats strfmt.Registry) error {
	if swag.IsZero(m.Links) { // not required
		return nil
	}

	if m.Links != nil {
		if err := m.Links.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subsystem_map" + "." + "subsystem" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystem) validateHosts(formats strfmt.Registry) error {
	if swag.IsZero(m.Hosts) { // not required
		return nil
	}

	for i := 0; i < len(m.Hosts); i++ {
		if swag.IsZero(m.Hosts[i]) { // not required
			continue
		}

		if m.Hosts[i] != nil {
			if err := m.Hosts[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("subsystem_map" + "." + "subsystem" + "." + "hosts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this application nvme access subsystem map subsystem based on the context it is used
func (m *ApplicationNvmeAccessSubsystemMapSubsystem) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateHosts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUUID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystem) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	if m.Links != nil {
		if err := m.Links.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subsystem_map" + "." + "subsystem" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystem) contextValidateHosts(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "subsystem_map"+"."+"subsystem"+"."+"hosts", "body", []*ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0(m.Hosts)); err != nil {
		return err
	}

	for i := 0; i < len(m.Hosts); i++ {

		if m.Hosts[i] != nil {
			if err := m.Hosts[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("subsystem_map" + "." + "subsystem" + "." + "hosts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystem) contextValidateName(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "subsystem_map"+"."+"subsystem"+"."+"name", "body", string(m.Name)); err != nil {
		return err
	}

	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystem) contextValidateUUID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "subsystem_map"+"."+"subsystem"+"."+"uuid", "body", string(m.UUID)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystem) UnmarshalBinary(b []byte) error {
	var res ApplicationNvmeAccessSubsystemMapSubsystem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0 application nvme access subsystem map subsystem hosts items0
//
// swagger:model ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0
type ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0 struct {

	// links
	Links *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links `json:"_links,omitempty"`

	// Host
	// Read Only: true
	Nqn string `json:"nqn,omitempty"`
}

// Validate validates this application nvme access subsystem map subsystem hosts items0
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0) validateLinks(formats strfmt.Registry) error {
	if swag.IsZero(m.Links) { // not required
		return nil
	}

	if m.Links != nil {
		if err := m.Links.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("_links")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this application nvme access subsystem map subsystem hosts items0 based on the context it is used
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNqn(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	if m.Links != nil {
		if err := m.Links.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("_links")
			}
			return err
		}
	}

	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0) contextValidateNqn(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "nqn", "body", string(m.Nqn)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0) UnmarshalBinary(b []byte) error {
	var res ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links application nvme access subsystem map subsystem hosts items0 links
//
// swagger:model ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links
type ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links struct {

	// self
	Self *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf `json:"self,omitempty"`
}

// Validate validates this application nvme access subsystem map subsystem hosts items0 links
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSelf(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links) validateSelf(formats strfmt.Registry) error {
	if swag.IsZero(m.Self) { // not required
		return nil
	}

	if m.Self != nil {
		if err := m.Self.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this application nvme access subsystem map subsystem hosts items0 links based on the context it is used
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSelf(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links) contextValidateSelf(ctx context.Context, formats strfmt.Registry) error {

	if m.Self != nil {
		if err := m.Self.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links) UnmarshalBinary(b []byte) error {
	var res ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0Links
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf application nvme access subsystem map subsystem hosts items0 links self
//
// swagger:model ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf
type ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf struct {

	// self
	Self *Href `json:"self,omitempty"`
}

// Validate validates this application nvme access subsystem map subsystem hosts items0 links self
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSelf(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf) validateSelf(formats strfmt.Registry) error {
	if swag.IsZero(m.Self) { // not required
		return nil
	}

	if m.Self != nil {
		if err := m.Self.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("_links" + "." + "self" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this application nvme access subsystem map subsystem hosts items0 links self based on the context it is used
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSelf(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf) contextValidateSelf(ctx context.Context, formats strfmt.Registry) error {

	if m.Self != nil {
		if err := m.Self.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("_links" + "." + "self" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf) UnmarshalBinary(b []byte) error {
	var res ApplicationNvmeAccessSubsystemMapSubsystemHostsItems0LinksSelf
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ApplicationNvmeAccessSubsystemMapSubsystemLinks application nvme access subsystem map subsystem links
//
// swagger:model ApplicationNvmeAccessSubsystemMapSubsystemLinks
type ApplicationNvmeAccessSubsystemMapSubsystemLinks struct {

	// self
	Self *Href `json:"self,omitempty"`
}

// Validate validates this application nvme access subsystem map subsystem links
func (m *ApplicationNvmeAccessSubsystemMapSubsystemLinks) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSelf(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystemLinks) validateSelf(formats strfmt.Registry) error {
	if swag.IsZero(m.Self) { // not required
		return nil
	}

	if m.Self != nil {
		if err := m.Self.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subsystem_map" + "." + "subsystem" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this application nvme access subsystem map subsystem links based on the context it is used
func (m *ApplicationNvmeAccessSubsystemMapSubsystemLinks) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSelf(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationNvmeAccessSubsystemMapSubsystemLinks) contextValidateSelf(ctx context.Context, formats strfmt.Registry) error {

	if m.Self != nil {
		if err := m.Self.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subsystem_map" + "." + "subsystem" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystemLinks) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationNvmeAccessSubsystemMapSubsystemLinks) UnmarshalBinary(b []byte) error {
	var res ApplicationNvmeAccessSubsystemMapSubsystemLinks
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
