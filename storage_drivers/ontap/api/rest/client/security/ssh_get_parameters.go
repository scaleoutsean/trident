// Code generated by go-swagger; DO NOT EDIT.

package security

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

// NewSSHGetParams creates a new SSHGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSSHGetParams() *SSHGetParams {
	return &SSHGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSSHGetParamsWithTimeout creates a new SSHGetParams object
// with the ability to set a timeout on a request.
func NewSSHGetParamsWithTimeout(timeout time.Duration) *SSHGetParams {
	return &SSHGetParams{
		timeout: timeout,
	}
}

// NewSSHGetParamsWithContext creates a new SSHGetParams object
// with the ability to set a context for a request.
func NewSSHGetParamsWithContext(ctx context.Context) *SSHGetParams {
	return &SSHGetParams{
		Context: ctx,
	}
}

// NewSSHGetParamsWithHTTPClient creates a new SSHGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewSSHGetParamsWithHTTPClient(client *http.Client) *SSHGetParams {
	return &SSHGetParams{
		HTTPClient: client,
	}
}

/*
SSHGetParams contains all the parameters to send to the API endpoint

	for the ssh get operation.

	Typically these are written to a http.Request.
*/
type SSHGetParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the ssh get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SSHGetParams) WithDefaults() *SSHGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the ssh get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SSHGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the ssh get params
func (o *SSHGetParams) WithTimeout(timeout time.Duration) *SSHGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the ssh get params
func (o *SSHGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the ssh get params
func (o *SSHGetParams) WithContext(ctx context.Context) *SSHGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the ssh get params
func (o *SSHGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the ssh get params
func (o *SSHGetParams) WithHTTPClient(client *http.Client) *SSHGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the ssh get params
func (o *SSHGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *SSHGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
