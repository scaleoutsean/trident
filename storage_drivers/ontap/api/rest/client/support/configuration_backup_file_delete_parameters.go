// Code generated by go-swagger; DO NOT EDIT.

package support

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

// NewConfigurationBackupFileDeleteParams creates a new ConfigurationBackupFileDeleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewConfigurationBackupFileDeleteParams() *ConfigurationBackupFileDeleteParams {
	return &ConfigurationBackupFileDeleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewConfigurationBackupFileDeleteParamsWithTimeout creates a new ConfigurationBackupFileDeleteParams object
// with the ability to set a timeout on a request.
func NewConfigurationBackupFileDeleteParamsWithTimeout(timeout time.Duration) *ConfigurationBackupFileDeleteParams {
	return &ConfigurationBackupFileDeleteParams{
		timeout: timeout,
	}
}

// NewConfigurationBackupFileDeleteParamsWithContext creates a new ConfigurationBackupFileDeleteParams object
// with the ability to set a context for a request.
func NewConfigurationBackupFileDeleteParamsWithContext(ctx context.Context) *ConfigurationBackupFileDeleteParams {
	return &ConfigurationBackupFileDeleteParams{
		Context: ctx,
	}
}

// NewConfigurationBackupFileDeleteParamsWithHTTPClient creates a new ConfigurationBackupFileDeleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewConfigurationBackupFileDeleteParamsWithHTTPClient(client *http.Client) *ConfigurationBackupFileDeleteParams {
	return &ConfigurationBackupFileDeleteParams{
		HTTPClient: client,
	}
}

/*
ConfigurationBackupFileDeleteParams contains all the parameters to send to the API endpoint

	for the configuration backup file delete operation.

	Typically these are written to a http.Request.
*/
type ConfigurationBackupFileDeleteParams struct {

	/* Name.

	   Name of the configuration backup to be deleted.
	*/
	NamePathParameter string

	/* NodeUUID.

	   UUID of the node that owns the configuration backup.
	*/
	NodeUUIDPathParameter string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the configuration backup file delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ConfigurationBackupFileDeleteParams) WithDefaults() *ConfigurationBackupFileDeleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the configuration backup file delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ConfigurationBackupFileDeleteParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) WithTimeout(timeout time.Duration) *ConfigurationBackupFileDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) WithContext(ctx context.Context) *ConfigurationBackupFileDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) WithHTTPClient(client *http.Client) *ConfigurationBackupFileDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNamePathParameter adds the name to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) WithNamePathParameter(name string) *ConfigurationBackupFileDeleteParams {
	o.SetNamePathParameter(name)
	return o
}

// SetNamePathParameter adds the name to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) SetNamePathParameter(name string) {
	o.NamePathParameter = name
}

// WithNodeUUIDPathParameter adds the nodeUUID to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) WithNodeUUIDPathParameter(nodeUUID string) *ConfigurationBackupFileDeleteParams {
	o.SetNodeUUIDPathParameter(nodeUUID)
	return o
}

// SetNodeUUIDPathParameter adds the nodeUuid to the configuration backup file delete params
func (o *ConfigurationBackupFileDeleteParams) SetNodeUUIDPathParameter(nodeUUID string) {
	o.NodeUUIDPathParameter = nodeUUID
}

// WriteToRequest writes these params to a swagger request
func (o *ConfigurationBackupFileDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.NamePathParameter); err != nil {
		return err
	}

	// path param node.uuid
	if err := r.SetPathParam("node.uuid", o.NodeUUIDPathParameter); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
