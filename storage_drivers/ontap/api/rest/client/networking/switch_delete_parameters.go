// Code generated by go-swagger; DO NOT EDIT.

package networking

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewSwitchDeleteParams creates a new SwitchDeleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSwitchDeleteParams() *SwitchDeleteParams {
	return &SwitchDeleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSwitchDeleteParamsWithTimeout creates a new SwitchDeleteParams object
// with the ability to set a timeout on a request.
func NewSwitchDeleteParamsWithTimeout(timeout time.Duration) *SwitchDeleteParams {
	return &SwitchDeleteParams{
		timeout: timeout,
	}
}

// NewSwitchDeleteParamsWithContext creates a new SwitchDeleteParams object
// with the ability to set a context for a request.
func NewSwitchDeleteParamsWithContext(ctx context.Context) *SwitchDeleteParams {
	return &SwitchDeleteParams{
		Context: ctx,
	}
}

// NewSwitchDeleteParamsWithHTTPClient creates a new SwitchDeleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewSwitchDeleteParamsWithHTTPClient(client *http.Client) *SwitchDeleteParams {
	return &SwitchDeleteParams{
		HTTPClient: client,
	}
}

/*
SwitchDeleteParams contains all the parameters to send to the API endpoint

	for the switch delete operation.

	Typically these are written to a http.Request.
*/
type SwitchDeleteParams struct {

	/* Name.

	   Switch Name.
	*/
	NamePathParameter string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the switch delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SwitchDeleteParams) WithDefaults() *SwitchDeleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the switch delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SwitchDeleteParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the switch delete params
func (o *SwitchDeleteParams) WithTimeout(timeout time.Duration) *SwitchDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the switch delete params
func (o *SwitchDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the switch delete params
func (o *SwitchDeleteParams) WithContext(ctx context.Context) *SwitchDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the switch delete params
func (o *SwitchDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the switch delete params
func (o *SwitchDeleteParams) WithHTTPClient(client *http.Client) *SwitchDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the switch delete params
func (o *SwitchDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNamePathParameter adds the name to the switch delete params
func (o *SwitchDeleteParams) WithNamePathParameter(name string) *SwitchDeleteParams {
	o.SetNamePathParameter(name)
	return o
}

// SetNamePathParameter adds the name to the switch delete params
func (o *SwitchDeleteParams) SetNamePathParameter(name string) {
	o.NamePathParameter = name
}

// WriteToRequest writes these params to a swagger request
func (o *SwitchDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.NamePathParameter); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
