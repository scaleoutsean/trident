// Code generated by go-swagger; DO NOT EDIT.

package s_a_n

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netapp/trident/storage_drivers/ontap/api/rest/models"
)

// FcpCollectionPerformanceMetricsGetReader is a Reader for the FcpCollectionPerformanceMetricsGet structure.
type FcpCollectionPerformanceMetricsGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *FcpCollectionPerformanceMetricsGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewFcpCollectionPerformanceMetricsGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewFcpCollectionPerformanceMetricsGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewFcpCollectionPerformanceMetricsGetOK creates a FcpCollectionPerformanceMetricsGetOK with default headers values
func NewFcpCollectionPerformanceMetricsGetOK() *FcpCollectionPerformanceMetricsGetOK {
	return &FcpCollectionPerformanceMetricsGetOK{}
}

/*
FcpCollectionPerformanceMetricsGetOK describes a response with status code 200, with default header values.

OK
*/
type FcpCollectionPerformanceMetricsGetOK struct {
	Payload *models.PerformanceFcpMetricResponse
}

// IsSuccess returns true when this fcp collection performance metrics get o k response has a 2xx status code
func (o *FcpCollectionPerformanceMetricsGetOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this fcp collection performance metrics get o k response has a 3xx status code
func (o *FcpCollectionPerformanceMetricsGetOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this fcp collection performance metrics get o k response has a 4xx status code
func (o *FcpCollectionPerformanceMetricsGetOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this fcp collection performance metrics get o k response has a 5xx status code
func (o *FcpCollectionPerformanceMetricsGetOK) IsServerError() bool {
	return false
}

// IsCode returns true when this fcp collection performance metrics get o k response a status code equal to that given
func (o *FcpCollectionPerformanceMetricsGetOK) IsCode(code int) bool {
	return code == 200
}

func (o *FcpCollectionPerformanceMetricsGetOK) Error() string {
	return fmt.Sprintf("[GET /protocols/san/fcp/services/{svm.uuid}/metrics][%d] fcpCollectionPerformanceMetricsGetOK  %+v", 200, o.Payload)
}

func (o *FcpCollectionPerformanceMetricsGetOK) String() string {
	return fmt.Sprintf("[GET /protocols/san/fcp/services/{svm.uuid}/metrics][%d] fcpCollectionPerformanceMetricsGetOK  %+v", 200, o.Payload)
}

func (o *FcpCollectionPerformanceMetricsGetOK) GetPayload() *models.PerformanceFcpMetricResponse {
	return o.Payload
}

func (o *FcpCollectionPerformanceMetricsGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PerformanceFcpMetricResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFcpCollectionPerformanceMetricsGetDefault creates a FcpCollectionPerformanceMetricsGetDefault with default headers values
func NewFcpCollectionPerformanceMetricsGetDefault(code int) *FcpCollectionPerformanceMetricsGetDefault {
	return &FcpCollectionPerformanceMetricsGetDefault{
		_statusCode: code,
	}
}

/*
FcpCollectionPerformanceMetricsGetDefault describes a response with status code -1, with default header values.

Error
*/
type FcpCollectionPerformanceMetricsGetDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the fcp collection performance metrics get default response
func (o *FcpCollectionPerformanceMetricsGetDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this fcp collection performance metrics get default response has a 2xx status code
func (o *FcpCollectionPerformanceMetricsGetDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this fcp collection performance metrics get default response has a 3xx status code
func (o *FcpCollectionPerformanceMetricsGetDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this fcp collection performance metrics get default response has a 4xx status code
func (o *FcpCollectionPerformanceMetricsGetDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this fcp collection performance metrics get default response has a 5xx status code
func (o *FcpCollectionPerformanceMetricsGetDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this fcp collection performance metrics get default response a status code equal to that given
func (o *FcpCollectionPerformanceMetricsGetDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *FcpCollectionPerformanceMetricsGetDefault) Error() string {
	return fmt.Sprintf("[GET /protocols/san/fcp/services/{svm.uuid}/metrics][%d] fcp_collection_performance_metrics_get default  %+v", o._statusCode, o.Payload)
}

func (o *FcpCollectionPerformanceMetricsGetDefault) String() string {
	return fmt.Sprintf("[GET /protocols/san/fcp/services/{svm.uuid}/metrics][%d] fcp_collection_performance_metrics_get default  %+v", o._statusCode, o.Payload)
}

func (o *FcpCollectionPerformanceMetricsGetDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *FcpCollectionPerformanceMetricsGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
