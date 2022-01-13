// Code generated by go-swagger; DO NOT EDIT.

package n_a_s

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

// NewCifsDomainPreferredDcDeleteParams creates a new CifsDomainPreferredDcDeleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCifsDomainPreferredDcDeleteParams() *CifsDomainPreferredDcDeleteParams {
	return &CifsDomainPreferredDcDeleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCifsDomainPreferredDcDeleteParamsWithTimeout creates a new CifsDomainPreferredDcDeleteParams object
// with the ability to set a timeout on a request.
func NewCifsDomainPreferredDcDeleteParamsWithTimeout(timeout time.Duration) *CifsDomainPreferredDcDeleteParams {
	return &CifsDomainPreferredDcDeleteParams{
		timeout: timeout,
	}
}

// NewCifsDomainPreferredDcDeleteParamsWithContext creates a new CifsDomainPreferredDcDeleteParams object
// with the ability to set a context for a request.
func NewCifsDomainPreferredDcDeleteParamsWithContext(ctx context.Context) *CifsDomainPreferredDcDeleteParams {
	return &CifsDomainPreferredDcDeleteParams{
		Context: ctx,
	}
}

// NewCifsDomainPreferredDcDeleteParamsWithHTTPClient creates a new CifsDomainPreferredDcDeleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewCifsDomainPreferredDcDeleteParamsWithHTTPClient(client *http.Client) *CifsDomainPreferredDcDeleteParams {
	return &CifsDomainPreferredDcDeleteParams{
		HTTPClient: client,
	}
}

/* CifsDomainPreferredDcDeleteParams contains all the parameters to send to the API endpoint
   for the cifs domain preferred dc delete operation.

   Typically these are written to a http.Request.
*/
type CifsDomainPreferredDcDeleteParams struct {

	/* Fqdn.

	   Fully Qualified Domain Name
	*/
	FqdnPathParameter string

	/* ServerIP.

	   Domain Controller IP address
	*/
	ServerIPPathParameter string

	/* SvmUUID.

	   UUID of the SVM to which this object belongs.
	*/
	SVMUUIDPathParameter string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the cifs domain preferred dc delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CifsDomainPreferredDcDeleteParams) WithDefaults() *CifsDomainPreferredDcDeleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the cifs domain preferred dc delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CifsDomainPreferredDcDeleteParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) WithTimeout(timeout time.Duration) *CifsDomainPreferredDcDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) WithContext(ctx context.Context) *CifsDomainPreferredDcDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) WithHTTPClient(client *http.Client) *CifsDomainPreferredDcDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFqdnPathParameter adds the fqdn to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) WithFqdnPathParameter(fqdn string) *CifsDomainPreferredDcDeleteParams {
	o.SetFqdnPathParameter(fqdn)
	return o
}

// SetFqdnPathParameter adds the fqdn to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) SetFqdnPathParameter(fqdn string) {
	o.FqdnPathParameter = fqdn
}

// WithServerIPPathParameter adds the serverIP to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) WithServerIPPathParameter(serverIP string) *CifsDomainPreferredDcDeleteParams {
	o.SetServerIPPathParameter(serverIP)
	return o
}

// SetServerIPPathParameter adds the serverIp to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) SetServerIPPathParameter(serverIP string) {
	o.ServerIPPathParameter = serverIP
}

// WithSVMUUIDPathParameter adds the svmUUID to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) WithSVMUUIDPathParameter(svmUUID string) *CifsDomainPreferredDcDeleteParams {
	o.SetSVMUUIDPathParameter(svmUUID)
	return o
}

// SetSVMUUIDPathParameter adds the svmUuid to the cifs domain preferred dc delete params
func (o *CifsDomainPreferredDcDeleteParams) SetSVMUUIDPathParameter(svmUUID string) {
	o.SVMUUIDPathParameter = svmUUID
}

// WriteToRequest writes these params to a swagger request
func (o *CifsDomainPreferredDcDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param fqdn
	if err := r.SetPathParam("fqdn", o.FqdnPathParameter); err != nil {
		return err
	}

	// path param server_ip
	if err := r.SetPathParam("server_ip", o.ServerIPPathParameter); err != nil {
		return err
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