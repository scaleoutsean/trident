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

// Error error
//
// swagger:model error
type Error struct {

	// Message arguments
	// Read Only: true
	Arguments []*ErrorArguments `json:"arguments,omitempty"`

	// Error code
	// Example: 4
	// Read Only: true
	Code string `json:"code,omitempty"`

	// Error message
	// Example: entry doesn't exist
	// Read Only: true
	Message string `json:"message,omitempty"`

	// The target parameter that caused the error.
	// Example: uuid
	// Read Only: true
	Target string `json:"target,omitempty"`
}

// Validate validates this error
func (m *Error) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateArguments(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Error) validateArguments(formats strfmt.Registry) error {
	if swag.IsZero(m.Arguments) { // not required
		return nil
	}

	for i := 0; i < len(m.Arguments); i++ {
		if swag.IsZero(m.Arguments[i]) { // not required
			continue
		}

		if m.Arguments[i] != nil {
			if err := m.Arguments[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("arguments" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this error based on the context it is used
func (m *Error) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateArguments(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCode(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMessage(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTarget(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Error) contextValidateArguments(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "arguments", "body", []*ErrorArguments(m.Arguments)); err != nil {
		return err
	}

	for i := 0; i < len(m.Arguments); i++ {

		if m.Arguments[i] != nil {
			if err := m.Arguments[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("arguments" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Error) contextValidateCode(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "code", "body", string(m.Code)); err != nil {
		return err
	}

	return nil
}

func (m *Error) contextValidateMessage(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "message", "body", string(m.Message)); err != nil {
		return err
	}

	return nil
}

func (m *Error) contextValidateTarget(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "target", "body", string(m.Target)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Error) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Error) UnmarshalBinary(b []byte) error {
	var res Error
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
