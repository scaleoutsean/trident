// Code generated by go-swagger; DO NOT EDIT.

package cluster

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netapp/trident/storage_drivers/ontap/api/rest/models"
)

// MetroclusterOperationCollectionGetReader is a Reader for the MetroclusterOperationCollectionGet structure.
type MetroclusterOperationCollectionGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MetroclusterOperationCollectionGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewMetroclusterOperationCollectionGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewMetroclusterOperationCollectionGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewMetroclusterOperationCollectionGetOK creates a MetroclusterOperationCollectionGetOK with default headers values
func NewMetroclusterOperationCollectionGetOK() *MetroclusterOperationCollectionGetOK {
	return &MetroclusterOperationCollectionGetOK{}
}

/*
MetroclusterOperationCollectionGetOK describes a response with status code 200, with default header values.

OK
*/
type MetroclusterOperationCollectionGetOK struct {
	Payload *models.MetroclusterOperationResponse
}

// IsSuccess returns true when this metrocluster operation collection get o k response has a 2xx status code
func (o *MetroclusterOperationCollectionGetOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this metrocluster operation collection get o k response has a 3xx status code
func (o *MetroclusterOperationCollectionGetOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this metrocluster operation collection get o k response has a 4xx status code
func (o *MetroclusterOperationCollectionGetOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this metrocluster operation collection get o k response has a 5xx status code
func (o *MetroclusterOperationCollectionGetOK) IsServerError() bool {
	return false
}

// IsCode returns true when this metrocluster operation collection get o k response a status code equal to that given
func (o *MetroclusterOperationCollectionGetOK) IsCode(code int) bool {
	return code == 200
}

func (o *MetroclusterOperationCollectionGetOK) Error() string {
	return fmt.Sprintf("[GET /cluster/metrocluster/operations][%d] metroclusterOperationCollectionGetOK  %+v", 200, o.Payload)
}

func (o *MetroclusterOperationCollectionGetOK) String() string {
	return fmt.Sprintf("[GET /cluster/metrocluster/operations][%d] metroclusterOperationCollectionGetOK  %+v", 200, o.Payload)
}

func (o *MetroclusterOperationCollectionGetOK) GetPayload() *models.MetroclusterOperationResponse {
	return o.Payload
}

func (o *MetroclusterOperationCollectionGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.MetroclusterOperationResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewMetroclusterOperationCollectionGetDefault creates a MetroclusterOperationCollectionGetDefault with default headers values
func NewMetroclusterOperationCollectionGetDefault(code int) *MetroclusterOperationCollectionGetDefault {
	return &MetroclusterOperationCollectionGetDefault{
		_statusCode: code,
	}
}

/*
	MetroclusterOperationCollectionGetDefault describes a response with status code -1, with default header values.

	ONTAP Error Response Codes

| Error Code | Description |
| ---------- | ----------- |
| 2425734 | An internal error occurred. Wait a few minutes, and try the operation again. For further assistance, contact technical support. |
*/
type MetroclusterOperationCollectionGetDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the metrocluster operation collection get default response
func (o *MetroclusterOperationCollectionGetDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this metrocluster operation collection get default response has a 2xx status code
func (o *MetroclusterOperationCollectionGetDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this metrocluster operation collection get default response has a 3xx status code
func (o *MetroclusterOperationCollectionGetDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this metrocluster operation collection get default response has a 4xx status code
func (o *MetroclusterOperationCollectionGetDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this metrocluster operation collection get default response has a 5xx status code
func (o *MetroclusterOperationCollectionGetDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this metrocluster operation collection get default response a status code equal to that given
func (o *MetroclusterOperationCollectionGetDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *MetroclusterOperationCollectionGetDefault) Error() string {
	return fmt.Sprintf("[GET /cluster/metrocluster/operations][%d] metrocluster_operation_collection_get default  %+v", o._statusCode, o.Payload)
}

func (o *MetroclusterOperationCollectionGetDefault) String() string {
	return fmt.Sprintf("[GET /cluster/metrocluster/operations][%d] metrocluster_operation_collection_get default  %+v", o._statusCode, o.Payload)
}

func (o *MetroclusterOperationCollectionGetDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *MetroclusterOperationCollectionGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
