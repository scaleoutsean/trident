// Code generated by go-swagger; DO NOT EDIT.

package storage

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
	"github.com/go-openapi/swag"
)

// NewMonitoredFileDeleteParams creates a new MonitoredFileDeleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewMonitoredFileDeleteParams() *MonitoredFileDeleteParams {
	return &MonitoredFileDeleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewMonitoredFileDeleteParamsWithTimeout creates a new MonitoredFileDeleteParams object
// with the ability to set a timeout on a request.
func NewMonitoredFileDeleteParamsWithTimeout(timeout time.Duration) *MonitoredFileDeleteParams {
	return &MonitoredFileDeleteParams{
		timeout: timeout,
	}
}

// NewMonitoredFileDeleteParamsWithContext creates a new MonitoredFileDeleteParams object
// with the ability to set a context for a request.
func NewMonitoredFileDeleteParamsWithContext(ctx context.Context) *MonitoredFileDeleteParams {
	return &MonitoredFileDeleteParams{
		Context: ctx,
	}
}

// NewMonitoredFileDeleteParamsWithHTTPClient creates a new MonitoredFileDeleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewMonitoredFileDeleteParamsWithHTTPClient(client *http.Client) *MonitoredFileDeleteParams {
	return &MonitoredFileDeleteParams{
		HTTPClient: client,
	}
}

/* MonitoredFileDeleteParams contains all the parameters to send to the API endpoint
   for the monitored file delete operation.

   Typically these are written to a http.Request.
*/
type MonitoredFileDeleteParams struct {

	/* ReturnRecords.

	   The default is false.  If set to true, the records are returned.
	*/
	ReturnRecordsQueryParameter *bool

	// UUID.
	UUIDPathParameter string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the monitored file delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *MonitoredFileDeleteParams) WithDefaults() *MonitoredFileDeleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the monitored file delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *MonitoredFileDeleteParams) SetDefaults() {
	var (
		returnRecordsQueryParameterDefault = bool(false)
	)

	val := MonitoredFileDeleteParams{
		ReturnRecordsQueryParameter: &returnRecordsQueryParameterDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the monitored file delete params
func (o *MonitoredFileDeleteParams) WithTimeout(timeout time.Duration) *MonitoredFileDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the monitored file delete params
func (o *MonitoredFileDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the monitored file delete params
func (o *MonitoredFileDeleteParams) WithContext(ctx context.Context) *MonitoredFileDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the monitored file delete params
func (o *MonitoredFileDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the monitored file delete params
func (o *MonitoredFileDeleteParams) WithHTTPClient(client *http.Client) *MonitoredFileDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the monitored file delete params
func (o *MonitoredFileDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithReturnRecordsQueryParameter adds the returnRecords to the monitored file delete params
func (o *MonitoredFileDeleteParams) WithReturnRecordsQueryParameter(returnRecords *bool) *MonitoredFileDeleteParams {
	o.SetReturnRecordsQueryParameter(returnRecords)
	return o
}

// SetReturnRecordsQueryParameter adds the returnRecords to the monitored file delete params
func (o *MonitoredFileDeleteParams) SetReturnRecordsQueryParameter(returnRecords *bool) {
	o.ReturnRecordsQueryParameter = returnRecords
}

// WithUUIDPathParameter adds the uuid to the monitored file delete params
func (o *MonitoredFileDeleteParams) WithUUIDPathParameter(uuid string) *MonitoredFileDeleteParams {
	o.SetUUIDPathParameter(uuid)
	return o
}

// SetUUIDPathParameter adds the uuid to the monitored file delete params
func (o *MonitoredFileDeleteParams) SetUUIDPathParameter(uuid string) {
	o.UUIDPathParameter = uuid
}

// WriteToRequest writes these params to a swagger request
func (o *MonitoredFileDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.ReturnRecordsQueryParameter != nil {

		// query param return_records
		var qrReturnRecords bool

		if o.ReturnRecordsQueryParameter != nil {
			qrReturnRecords = *o.ReturnRecordsQueryParameter
		}
		qReturnRecords := swag.FormatBool(qrReturnRecords)
		if qReturnRecords != "" {

			if err := r.SetQueryParam("return_records", qReturnRecords); err != nil {
				return err
			}
		}
	}

	// path param uuid
	if err := r.SetPathParam("uuid", o.UUIDPathParameter); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}