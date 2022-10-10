// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Rfc2307 rfc2307
//
// swagger:model rfc2307
type Rfc2307 struct {

	// attribute
	Attribute *Rfc2307AttributeType `json:"attribute,omitempty"`

	// cn
	Cn *Rfc2307Cn `json:"cn,omitempty"`

	// member
	Member *Rfc2307Member `json:"member,omitempty"`

	// nis
	Nis *Rfc2307Nis `json:"nis,omitempty"`

	// posix
	Posix *Rfc2307Posix `json:"posix,omitempty"`
}

// Validate validates this rfc2307
func (m *Rfc2307) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAttribute(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMember(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNis(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePosix(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Rfc2307) validateAttribute(formats strfmt.Registry) error {
	if swag.IsZero(m.Attribute) { // not required
		return nil
	}

	if m.Attribute != nil {
		if err := m.Attribute.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("attribute")
			}
			return err
		}
	}

	return nil
}

func (m *Rfc2307) validateCn(formats strfmt.Registry) error {
	if swag.IsZero(m.Cn) { // not required
		return nil
	}

	if m.Cn != nil {
		if err := m.Cn.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cn")
			}
			return err
		}
	}

	return nil
}

func (m *Rfc2307) validateMember(formats strfmt.Registry) error {
	if swag.IsZero(m.Member) { // not required
		return nil
	}

	if m.Member != nil {
		if err := m.Member.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("member")
			}
			return err
		}
	}

	return nil
}

func (m *Rfc2307) validateNis(formats strfmt.Registry) error {
	if swag.IsZero(m.Nis) { // not required
		return nil
	}

	if m.Nis != nil {
		if err := m.Nis.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("nis")
			}
			return err
		}
	}

	return nil
}

func (m *Rfc2307) validatePosix(formats strfmt.Registry) error {
	if swag.IsZero(m.Posix) { // not required
		return nil
	}

	if m.Posix != nil {
		if err := m.Posix.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("posix")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this rfc2307 based on the context it is used
func (m *Rfc2307) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAttribute(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCn(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMember(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNis(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePosix(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Rfc2307) contextValidateAttribute(ctx context.Context, formats strfmt.Registry) error {

	if m.Attribute != nil {
		if err := m.Attribute.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("attribute")
			}
			return err
		}
	}

	return nil
}

func (m *Rfc2307) contextValidateCn(ctx context.Context, formats strfmt.Registry) error {

	if m.Cn != nil {
		if err := m.Cn.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cn")
			}
			return err
		}
	}

	return nil
}

func (m *Rfc2307) contextValidateMember(ctx context.Context, formats strfmt.Registry) error {

	if m.Member != nil {
		if err := m.Member.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("member")
			}
			return err
		}
	}

	return nil
}

func (m *Rfc2307) contextValidateNis(ctx context.Context, formats strfmt.Registry) error {

	if m.Nis != nil {
		if err := m.Nis.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("nis")
			}
			return err
		}
	}

	return nil
}

func (m *Rfc2307) contextValidatePosix(ctx context.Context, formats strfmt.Registry) error {

	if m.Posix != nil {
		if err := m.Posix.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("posix")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Rfc2307) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Rfc2307) UnmarshalBinary(b []byte) error {
	var res Rfc2307
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Rfc2307AttributeType rfc2307 attribute type
//
// swagger:model Rfc2307AttributeType
type Rfc2307AttributeType struct {

	// RFC 2307 gecos attribute.
	// Example: name
	Gecos string `json:"gecos,omitempty"`

	// RFC 2307 gidNumber attribute.
	// Example: msSFU30GidNumber
	GidNumber string `json:"gid_number,omitempty"`

	// RFC 2307 homeDirectory attribute.
	// Example: msSFU30HomeDirectory
	HomeDirectory string `json:"home_directory,omitempty"`

	// RFC 2307 loginShell attribute.
	// Example: msSFU30LoginShell
	LoginShell string `json:"login_shell,omitempty"`

	// RFC 1274 userid attribute used by RFC 2307 as UID.
	// Example: sAMAccountName
	UID string `json:"uid,omitempty"`

	// RFC 2307 uidNumber attribute.
	// Example: msSFU30UidNumber
	UIDNumber string `json:"uid_number,omitempty"`

	// RFC 2256 userPassword attribute used by RFC 2307.
	// Example: msSFU30Password
	UserPassword string `json:"user_password,omitempty"`
}

// Validate validates this rfc2307 attribute type
func (m *Rfc2307AttributeType) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this rfc2307 attribute type based on context it is used
func (m *Rfc2307AttributeType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Rfc2307AttributeType) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Rfc2307AttributeType) UnmarshalBinary(b []byte) error {
	var res Rfc2307AttributeType
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Rfc2307Cn rfc2307 cn
//
// swagger:model Rfc2307Cn
type Rfc2307Cn struct {

	// RFC 2256 cn attribute used by RFC 2307 when working with groups.
	// Example: cn
	Group string `json:"group,omitempty"`

	// RFC 2256 cn attribute used by RFC 2307 when working with netgroups.
	// Example: name
	Netgroup string `json:"netgroup,omitempty"`
}

// Validate validates this rfc2307 cn
func (m *Rfc2307Cn) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this rfc2307 cn based on context it is used
func (m *Rfc2307Cn) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Rfc2307Cn) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Rfc2307Cn) UnmarshalBinary(b []byte) error {
	var res Rfc2307Cn
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Rfc2307Member rfc2307 member
//
// swagger:model Rfc2307Member
type Rfc2307Member struct {

	// RFC 2307 memberNisNetgroup attribute.
	// Example: msSFU30MemberNisNetgroup
	NisNetgroup string `json:"nis_netgroup,omitempty"`

	// RFC 2307 memberUid attribute.
	// Example: msSFU30MemberUid
	UID string `json:"uid,omitempty"`
}

// Validate validates this rfc2307 member
func (m *Rfc2307Member) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this rfc2307 member based on context it is used
func (m *Rfc2307Member) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Rfc2307Member) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Rfc2307Member) UnmarshalBinary(b []byte) error {
	var res Rfc2307Member
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Rfc2307Nis rfc2307 nis
//
// swagger:model Rfc2307Nis
type Rfc2307Nis struct {

	// RFC 2307 nisMapEntry attribute.
	// Example: msSFU30NisMapEntry
	Mapentry string `json:"mapentry,omitempty"`

	// RFC 2307 nisMapName attribute.
	// Example: msSFU30NisMapName
	Mapname string `json:"mapname,omitempty"`

	// RFC 2307 nisNetgroup object class.
	// Example: msSFU30NisNetGroup
	Netgroup string `json:"netgroup,omitempty"`

	// RFC 2307 nisNetgroupTriple attribute.
	// Example: msSFU30MemberOfNisNetgroup
	NetgroupTriple string `json:"netgroup_triple,omitempty"`

	// RFC 2307 nisObject object class.
	// Example: msSFU30NisObject
	Object string `json:"object,omitempty"`
}

// Validate validates this rfc2307 nis
func (m *Rfc2307Nis) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this rfc2307 nis based on context it is used
func (m *Rfc2307Nis) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Rfc2307Nis) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Rfc2307Nis) UnmarshalBinary(b []byte) error {
	var res Rfc2307Nis
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Rfc2307Posix rfc2307 posix
//
// swagger:model Rfc2307Posix
type Rfc2307Posix struct {

	// RFC 2307 posixAccount object class.
	// Example: User
	Account string `json:"account,omitempty"`

	// RFC 2307 posixGroup object class.
	// Example: Group
	Group string `json:"group,omitempty"`
}

// Validate validates this rfc2307 posix
func (m *Rfc2307Posix) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this rfc2307 posix based on context it is used
func (m *Rfc2307Posix) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Rfc2307Posix) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Rfc2307Posix) UnmarshalBinary(b []byte) error {
	var res Rfc2307Posix
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}