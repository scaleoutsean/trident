// Code generated by go-swagger; DO NOT EDIT.

package cluster

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

// NewMetroclusterSvmGetParams creates a new MetroclusterSvmGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewMetroclusterSvmGetParams() *MetroclusterSvmGetParams {
	return &MetroclusterSvmGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewMetroclusterSvmGetParamsWithTimeout creates a new MetroclusterSvmGetParams object
// with the ability to set a timeout on a request.
func NewMetroclusterSvmGetParamsWithTimeout(timeout time.Duration) *MetroclusterSvmGetParams {
	return &MetroclusterSvmGetParams{
		timeout: timeout,
	}
}

// NewMetroclusterSvmGetParamsWithContext creates a new MetroclusterSvmGetParams object
// with the ability to set a context for a request.
func NewMetroclusterSvmGetParamsWithContext(ctx context.Context) *MetroclusterSvmGetParams {
	return &MetroclusterSvmGetParams{
		Context: ctx,
	}
}

// NewMetroclusterSvmGetParamsWithHTTPClient creates a new MetroclusterSvmGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewMetroclusterSvmGetParamsWithHTTPClient(client *http.Client) *MetroclusterSvmGetParams {
	return &MetroclusterSvmGetParams{
		HTTPClient: client,
	}
}

/*
MetroclusterSvmGetParams contains all the parameters to send to the API endpoint

	for the metrocluster svm get operation.

	Typically these are written to a http.Request.
*/
type MetroclusterSvmGetParams struct {

	/* ClusterUUID.

	   Cluster ID
	*/
	ClusterUUIDPathParameter string

	/* Fields.

	   Specify the fields to return.
	*/
	FieldsQueryParameter []string

	/* SvmUUID.

	   SVM UUID
	*/
	SVMUUIDPathParameter string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the metrocluster svm get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *MetroclusterSvmGetParams) WithDefaults() *MetroclusterSvmGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the metrocluster svm get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *MetroclusterSvmGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) WithTimeout(timeout time.Duration) *MetroclusterSvmGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) WithContext(ctx context.Context) *MetroclusterSvmGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) WithHTTPClient(client *http.Client) *MetroclusterSvmGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterUUIDPathParameter adds the clusterUUID to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) WithClusterUUIDPathParameter(clusterUUID string) *MetroclusterSvmGetParams {
	o.SetClusterUUIDPathParameter(clusterUUID)
	return o
}

// SetClusterUUIDPathParameter adds the clusterUuid to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) SetClusterUUIDPathParameter(clusterUUID string) {
	o.ClusterUUIDPathParameter = clusterUUID
}

// WithFieldsQueryParameter adds the fields to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) WithFieldsQueryParameter(fields []string) *MetroclusterSvmGetParams {
	o.SetFieldsQueryParameter(fields)
	return o
}

// SetFieldsQueryParameter adds the fields to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) SetFieldsQueryParameter(fields []string) {
	o.FieldsQueryParameter = fields
}

// WithSVMUUIDPathParameter adds the svmUUID to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) WithSVMUUIDPathParameter(svmUUID string) *MetroclusterSvmGetParams {
	o.SetSVMUUIDPathParameter(svmUUID)
	return o
}

// SetSVMUUIDPathParameter adds the svmUuid to the metrocluster svm get params
func (o *MetroclusterSvmGetParams) SetSVMUUIDPathParameter(svmUUID string) {
	o.SVMUUIDPathParameter = svmUUID
}

// WriteToRequest writes these params to a swagger request
func (o *MetroclusterSvmGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param cluster.uuid
	if err := r.SetPathParam("cluster.uuid", o.ClusterUUIDPathParameter); err != nil {
		return err
	}

	if o.FieldsQueryParameter != nil {

		// binding items for fields
		joinedFields := o.bindParamFields(reg)

		// query array param fields
		if err := r.SetQueryParam("fields", joinedFields...); err != nil {
			return err
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

// bindParamMetroclusterSvmGet binds the parameter fields
func (o *MetroclusterSvmGetParams) bindParamFields(formats strfmt.Registry) []string {
	fieldsIR := o.FieldsQueryParameter

	var fieldsIC []string
	for _, fieldsIIR := range fieldsIR { // explode []string

		fieldsIIV := fieldsIIR // string as string
		fieldsIC = append(fieldsIC, fieldsIIV)
	}

	// items.CollectionFormat: "csv"
	fieldsIS := swag.JoinByFormat(fieldsIC, "csv")

	return fieldsIS
}
