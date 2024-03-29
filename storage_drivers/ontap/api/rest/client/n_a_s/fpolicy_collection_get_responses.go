// Code generated by go-swagger; DO NOT EDIT.

package n_a_s

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netapp/trident/storage_drivers/ontap/api/rest/models"
)

// FpolicyCollectionGetReader is a Reader for the FpolicyCollectionGet structure.
type FpolicyCollectionGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *FpolicyCollectionGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewFpolicyCollectionGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewFpolicyCollectionGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewFpolicyCollectionGetOK creates a FpolicyCollectionGetOK with default headers values
func NewFpolicyCollectionGetOK() *FpolicyCollectionGetOK {
	return &FpolicyCollectionGetOK{}
}

/*
FpolicyCollectionGetOK describes a response with status code 200, with default header values.

OK
*/
type FpolicyCollectionGetOK struct {
	Payload *models.FpolicyResponse
}

// IsSuccess returns true when this fpolicy collection get o k response has a 2xx status code
func (o *FpolicyCollectionGetOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this fpolicy collection get o k response has a 3xx status code
func (o *FpolicyCollectionGetOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this fpolicy collection get o k response has a 4xx status code
func (o *FpolicyCollectionGetOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this fpolicy collection get o k response has a 5xx status code
func (o *FpolicyCollectionGetOK) IsServerError() bool {
	return false
}

// IsCode returns true when this fpolicy collection get o k response a status code equal to that given
func (o *FpolicyCollectionGetOK) IsCode(code int) bool {
	return code == 200
}

func (o *FpolicyCollectionGetOK) Error() string {
	return fmt.Sprintf("[GET /protocols/fpolicy][%d] fpolicyCollectionGetOK  %+v", 200, o.Payload)
}

func (o *FpolicyCollectionGetOK) String() string {
	return fmt.Sprintf("[GET /protocols/fpolicy][%d] fpolicyCollectionGetOK  %+v", 200, o.Payload)
}

func (o *FpolicyCollectionGetOK) GetPayload() *models.FpolicyResponse {
	return o.Payload
}

func (o *FpolicyCollectionGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.FpolicyResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFpolicyCollectionGetDefault creates a FpolicyCollectionGetDefault with default headers values
func NewFpolicyCollectionGetDefault(code int) *FpolicyCollectionGetDefault {
	return &FpolicyCollectionGetDefault{
		_statusCode: code,
	}
}

/*
FpolicyCollectionGetDefault describes a response with status code -1, with default header values.

Error
*/
type FpolicyCollectionGetDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the fpolicy collection get default response
func (o *FpolicyCollectionGetDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this fpolicy collection get default response has a 2xx status code
func (o *FpolicyCollectionGetDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this fpolicy collection get default response has a 3xx status code
func (o *FpolicyCollectionGetDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this fpolicy collection get default response has a 4xx status code
func (o *FpolicyCollectionGetDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this fpolicy collection get default response has a 5xx status code
func (o *FpolicyCollectionGetDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this fpolicy collection get default response a status code equal to that given
func (o *FpolicyCollectionGetDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *FpolicyCollectionGetDefault) Error() string {
	return fmt.Sprintf("[GET /protocols/fpolicy][%d] fpolicy_collection_get default  %+v", o._statusCode, o.Payload)
}

func (o *FpolicyCollectionGetDefault) String() string {
	return fmt.Sprintf("[GET /protocols/fpolicy][%d] fpolicy_collection_get default  %+v", o._statusCode, o.Payload)
}

func (o *FpolicyCollectionGetDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *FpolicyCollectionGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
