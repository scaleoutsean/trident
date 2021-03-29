// Code generated by go-swagger; DO NOT EDIT.

package snaplock

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

// NewSnaplockLogDeleteParams creates a new SnaplockLogDeleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSnaplockLogDeleteParams() *SnaplockLogDeleteParams {
	return &SnaplockLogDeleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSnaplockLogDeleteParamsWithTimeout creates a new SnaplockLogDeleteParams object
// with the ability to set a timeout on a request.
func NewSnaplockLogDeleteParamsWithTimeout(timeout time.Duration) *SnaplockLogDeleteParams {
	return &SnaplockLogDeleteParams{
		timeout: timeout,
	}
}

// NewSnaplockLogDeleteParamsWithContext creates a new SnaplockLogDeleteParams object
// with the ability to set a context for a request.
func NewSnaplockLogDeleteParamsWithContext(ctx context.Context) *SnaplockLogDeleteParams {
	return &SnaplockLogDeleteParams{
		Context: ctx,
	}
}

// NewSnaplockLogDeleteParamsWithHTTPClient creates a new SnaplockLogDeleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewSnaplockLogDeleteParamsWithHTTPClient(client *http.Client) *SnaplockLogDeleteParams {
	return &SnaplockLogDeleteParams{
		HTTPClient: client,
	}
}

/* SnaplockLogDeleteParams contains all the parameters to send to the API endpoint
   for the snaplock log delete operation.

   Typically these are written to a http.Request.
*/
type SnaplockLogDeleteParams struct {

	/* ReturnTimeout.

	   The number of seconds to allow the call to execute before returning. When doing a POST, PATCH, or DELETE operation on a single record, the default is 0 seconds.  This means that if an asynchronous operation is started, the server immediately returns HTTP code 202 (Accepted) along with a link to the job.  If a non-zero value is specified for POST, PATCH, or DELETE operations, ONTAP waits that length of time to see if the job completes so it can return something other than 202.
	*/
	ReturnTimeout *int64

	/* SvmUUID.

	   SVM UUID
	*/
	SVMUUIDPathParameter string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the snaplock log delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SnaplockLogDeleteParams) WithDefaults() *SnaplockLogDeleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the snaplock log delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SnaplockLogDeleteParams) SetDefaults() {
	var (
		returnTimeoutDefault = int64(0)
	)

	val := SnaplockLogDeleteParams{
		ReturnTimeout: &returnTimeoutDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the snaplock log delete params
func (o *SnaplockLogDeleteParams) WithTimeout(timeout time.Duration) *SnaplockLogDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the snaplock log delete params
func (o *SnaplockLogDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the snaplock log delete params
func (o *SnaplockLogDeleteParams) WithContext(ctx context.Context) *SnaplockLogDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the snaplock log delete params
func (o *SnaplockLogDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the snaplock log delete params
func (o *SnaplockLogDeleteParams) WithHTTPClient(client *http.Client) *SnaplockLogDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the snaplock log delete params
func (o *SnaplockLogDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithReturnTimeout adds the returnTimeout to the snaplock log delete params
func (o *SnaplockLogDeleteParams) WithReturnTimeout(returnTimeout *int64) *SnaplockLogDeleteParams {
	o.SetReturnTimeout(returnTimeout)
	return o
}

// SetReturnTimeout adds the returnTimeout to the snaplock log delete params
func (o *SnaplockLogDeleteParams) SetReturnTimeout(returnTimeout *int64) {
	o.ReturnTimeout = returnTimeout
}

// WithSVMUUIDPathParameter adds the svmUUID to the snaplock log delete params
func (o *SnaplockLogDeleteParams) WithSVMUUIDPathParameter(svmUUID string) *SnaplockLogDeleteParams {
	o.SetSVMUUIDPathParameter(svmUUID)
	return o
}

// SetSVMUUIDPathParameter adds the svmUuid to the snaplock log delete params
func (o *SnaplockLogDeleteParams) SetSVMUUIDPathParameter(svmUUID string) {
	o.SVMUUIDPathParameter = svmUUID
}

// WriteToRequest writes these params to a swagger request
func (o *SnaplockLogDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.ReturnTimeout != nil {

		// query param return_timeout
		var qrReturnTimeout int64

		if o.ReturnTimeout != nil {
			qrReturnTimeout = *o.ReturnTimeout
		}
		qReturnTimeout := swag.FormatInt64(qrReturnTimeout)
		if qReturnTimeout != "" {

			if err := r.SetQueryParam("return_timeout", qReturnTimeout); err != nil {
				return err
			}
		}
	}

	// path param svm.uuid
	if err := r.SetPathParam("svm.uuid", o.SVMUUIDPathParameter); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}