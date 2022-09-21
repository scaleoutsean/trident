// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// IPSubnet ip subnet
//
// swagger:model ip_subnet
type IPSubnet struct {

	// links
	Links *IPSubnetLinks `json:"_links,omitempty"`

	// available count
	// Read Only: true
	AvailableCount int64 `json:"available_count,omitempty"`

	// available ip ranges
	// Read Only: true
	AvailableIPRanges []*IPAddressRange `json:"available_ip_ranges,omitempty"`

	// broadcast domain
	BroadcastDomain *IPSubnetBroadcastDomain `json:"broadcast_domain,omitempty"`

	// This action will fail if any existing interface is using an IP address in the ranges provided. Set this to false to associate any manually addressed interfaces with the subnet and allow the action to succeed.
	FailIfLifsConflict *bool `json:"fail_if_lifs_conflict,omitempty"`

	// The IP address of the gateway for this subnet.
	// Example: 10.1.1.1
	Gateway string `json:"gateway,omitempty"`

	// ip ranges
	IPRanges []*IPAddressRange `json:"ip_ranges,omitempty"`

	// ipspace
	Ipspace *IPSubnetIpspace `json:"ipspace,omitempty"`

	// Subnet name
	// Example: subnet1
	Name string `json:"name,omitempty"`

	// subnet
	Subnet *IPInfo `json:"subnet,omitempty"`

	// total count
	// Read Only: true
	TotalCount int64 `json:"total_count,omitempty"`

	// used count
	// Read Only: true
	UsedCount int64 `json:"used_count,omitempty"`

	// The UUID that uniquely identifies the subnet.
	// Example: 1cd8a442-86d1-11e0-ae1c-123478563412
	// Read Only: true
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this ip subnet
func (m *IPSubnet) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAvailableIPRanges(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBroadcastDomain(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIPRanges(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIpspace(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubnet(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnet) validateLinks(formats strfmt.Registry) error {
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

func (m *IPSubnet) validateAvailableIPRanges(formats strfmt.Registry) error {
	if swag.IsZero(m.AvailableIPRanges) { // not required
		return nil
	}

	for i := 0; i < len(m.AvailableIPRanges); i++ {
		if swag.IsZero(m.AvailableIPRanges[i]) { // not required
			continue
		}

		if m.AvailableIPRanges[i] != nil {
			if err := m.AvailableIPRanges[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("available_ip_ranges" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *IPSubnet) validateBroadcastDomain(formats strfmt.Registry) error {
	if swag.IsZero(m.BroadcastDomain) { // not required
		return nil
	}

	if m.BroadcastDomain != nil {
		if err := m.BroadcastDomain.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("broadcast_domain")
			}
			return err
		}
	}

	return nil
}

func (m *IPSubnet) validateIPRanges(formats strfmt.Registry) error {
	if swag.IsZero(m.IPRanges) { // not required
		return nil
	}

	for i := 0; i < len(m.IPRanges); i++ {
		if swag.IsZero(m.IPRanges[i]) { // not required
			continue
		}

		if m.IPRanges[i] != nil {
			if err := m.IPRanges[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("ip_ranges" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *IPSubnet) validateIpspace(formats strfmt.Registry) error {
	if swag.IsZero(m.Ipspace) { // not required
		return nil
	}

	if m.Ipspace != nil {
		if err := m.Ipspace.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ipspace")
			}
			return err
		}
	}

	return nil
}

func (m *IPSubnet) validateSubnet(formats strfmt.Registry) error {
	if swag.IsZero(m.Subnet) { // not required
		return nil
	}

	if m.Subnet != nil {
		if err := m.Subnet.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subnet")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this ip subnet based on the context it is used
func (m *IPSubnet) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateAvailableCount(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateAvailableIPRanges(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateBroadcastDomain(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateIPRanges(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateIpspace(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSubnet(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTotalCount(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUsedCount(ctx, formats); err != nil {
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

func (m *IPSubnet) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

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

func (m *IPSubnet) contextValidateAvailableCount(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "available_count", "body", int64(m.AvailableCount)); err != nil {
		return err
	}

	return nil
}

func (m *IPSubnet) contextValidateAvailableIPRanges(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "available_ip_ranges", "body", []*IPAddressRange(m.AvailableIPRanges)); err != nil {
		return err
	}

	for i := 0; i < len(m.AvailableIPRanges); i++ {

		if m.AvailableIPRanges[i] != nil {
			if err := m.AvailableIPRanges[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("available_ip_ranges" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *IPSubnet) contextValidateBroadcastDomain(ctx context.Context, formats strfmt.Registry) error {

	if m.BroadcastDomain != nil {
		if err := m.BroadcastDomain.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("broadcast_domain")
			}
			return err
		}
	}

	return nil
}

func (m *IPSubnet) contextValidateIPRanges(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.IPRanges); i++ {

		if m.IPRanges[i] != nil {
			if err := m.IPRanges[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("ip_ranges" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *IPSubnet) contextValidateIpspace(ctx context.Context, formats strfmt.Registry) error {

	if m.Ipspace != nil {
		if err := m.Ipspace.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ipspace")
			}
			return err
		}
	}

	return nil
}

func (m *IPSubnet) contextValidateSubnet(ctx context.Context, formats strfmt.Registry) error {

	if m.Subnet != nil {
		if err := m.Subnet.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subnet")
			}
			return err
		}
	}

	return nil
}

func (m *IPSubnet) contextValidateTotalCount(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "total_count", "body", int64(m.TotalCount)); err != nil {
		return err
	}

	return nil
}

func (m *IPSubnet) contextValidateUsedCount(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "used_count", "body", int64(m.UsedCount)); err != nil {
		return err
	}

	return nil
}

func (m *IPSubnet) contextValidateUUID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "uuid", "body", string(m.UUID)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *IPSubnet) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IPSubnet) UnmarshalBinary(b []byte) error {
	var res IPSubnet
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// IPSubnetBroadcastDomain The broadcast domain that the subnet is associated with. Either the UUID or name must be supplied on POST.
//
// swagger:model IPSubnetBroadcastDomain
type IPSubnetBroadcastDomain struct {

	// links
	Links *IPSubnetBroadcastDomainLinks `json:"_links,omitempty"`

	// Name of the broadcast domain, scoped to its IPspace
	// Example: bd1
	Name string `json:"name,omitempty"`

	// Broadcast domain UUID
	// Example: 1cd8a442-86d1-11e0-ae1c-123478563412
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this IP subnet broadcast domain
func (m *IPSubnetBroadcastDomain) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetBroadcastDomain) validateLinks(formats strfmt.Registry) error {
	if swag.IsZero(m.Links) { // not required
		return nil
	}

	if m.Links != nil {
		if err := m.Links.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("broadcast_domain" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this IP subnet broadcast domain based on the context it is used
func (m *IPSubnetBroadcastDomain) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetBroadcastDomain) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	if m.Links != nil {
		if err := m.Links.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("broadcast_domain" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *IPSubnetBroadcastDomain) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IPSubnetBroadcastDomain) UnmarshalBinary(b []byte) error {
	var res IPSubnetBroadcastDomain
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// IPSubnetBroadcastDomainLinks IP subnet broadcast domain links
//
// swagger:model IPSubnetBroadcastDomainLinks
type IPSubnetBroadcastDomainLinks struct {

	// self
	Self *Href `json:"self,omitempty"`
}

// Validate validates this IP subnet broadcast domain links
func (m *IPSubnetBroadcastDomainLinks) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSelf(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetBroadcastDomainLinks) validateSelf(formats strfmt.Registry) error {
	if swag.IsZero(m.Self) { // not required
		return nil
	}

	if m.Self != nil {
		if err := m.Self.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("broadcast_domain" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this IP subnet broadcast domain links based on the context it is used
func (m *IPSubnetBroadcastDomainLinks) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSelf(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetBroadcastDomainLinks) contextValidateSelf(ctx context.Context, formats strfmt.Registry) error {

	if m.Self != nil {
		if err := m.Self.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("broadcast_domain" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *IPSubnetBroadcastDomainLinks) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IPSubnetBroadcastDomainLinks) UnmarshalBinary(b []byte) error {
	var res IPSubnetBroadcastDomainLinks
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// IPSubnetIpspace The IPspace that the subnet is associated with. Either the UUID or name must be supplied on POST.
//
// swagger:model IPSubnetIpspace
type IPSubnetIpspace struct {

	// links
	Links *IPSubnetIpspaceLinks `json:"_links,omitempty"`

	// IPspace name
	// Example: exchange
	Name string `json:"name,omitempty"`

	// IPspace UUID
	// Example: 1cd8a442-86d1-11e0-ae1c-123478563412
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this IP subnet ipspace
func (m *IPSubnetIpspace) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetIpspace) validateLinks(formats strfmt.Registry) error {
	if swag.IsZero(m.Links) { // not required
		return nil
	}

	if m.Links != nil {
		if err := m.Links.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ipspace" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this IP subnet ipspace based on the context it is used
func (m *IPSubnetIpspace) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLinks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetIpspace) contextValidateLinks(ctx context.Context, formats strfmt.Registry) error {

	if m.Links != nil {
		if err := m.Links.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ipspace" + "." + "_links")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *IPSubnetIpspace) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IPSubnetIpspace) UnmarshalBinary(b []byte) error {
	var res IPSubnetIpspace
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// IPSubnetIpspaceLinks IP subnet ipspace links
//
// swagger:model IPSubnetIpspaceLinks
type IPSubnetIpspaceLinks struct {

	// self
	Self *Href `json:"self,omitempty"`
}

// Validate validates this IP subnet ipspace links
func (m *IPSubnetIpspaceLinks) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSelf(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetIpspaceLinks) validateSelf(formats strfmt.Registry) error {
	if swag.IsZero(m.Self) { // not required
		return nil
	}

	if m.Self != nil {
		if err := m.Self.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ipspace" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this IP subnet ipspace links based on the context it is used
func (m *IPSubnetIpspaceLinks) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSelf(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetIpspaceLinks) contextValidateSelf(ctx context.Context, formats strfmt.Registry) error {

	if m.Self != nil {
		if err := m.Self.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ipspace" + "." + "_links" + "." + "self")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *IPSubnetIpspaceLinks) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IPSubnetIpspaceLinks) UnmarshalBinary(b []byte) error {
	var res IPSubnetIpspaceLinks
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// IPSubnetLinks IP subnet links
//
// swagger:model IPSubnetLinks
type IPSubnetLinks struct {

	// self
	Self *Href `json:"self,omitempty"`
}

// Validate validates this IP subnet links
func (m *IPSubnetLinks) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSelf(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetLinks) validateSelf(formats strfmt.Registry) error {
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

// ContextValidate validate this IP subnet links based on the context it is used
func (m *IPSubnetLinks) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSelf(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IPSubnetLinks) contextValidateSelf(ctx context.Context, formats strfmt.Registry) error {

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
func (m *IPSubnetLinks) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IPSubnetLinks) UnmarshalBinary(b []byte) error {
	var res IPSubnetLinks
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}