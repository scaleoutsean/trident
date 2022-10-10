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

// ShadowcopyModifyReader is a Reader for the ShadowcopyModify structure.
type ShadowcopyModifyReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ShadowcopyModifyReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewShadowcopyModifyOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewShadowcopyModifyDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewShadowcopyModifyOK creates a ShadowcopyModifyOK with default headers values
func NewShadowcopyModifyOK() *ShadowcopyModifyOK {
	return &ShadowcopyModifyOK{}
}

/*
ShadowcopyModifyOK describes a response with status code 200, with default header values.

OK
*/
type ShadowcopyModifyOK struct {
	Payload *models.ShadowcopyAddFiles
}

// IsSuccess returns true when this shadowcopy modify o k response has a 2xx status code
func (o *ShadowcopyModifyOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this shadowcopy modify o k response has a 3xx status code
func (o *ShadowcopyModifyOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this shadowcopy modify o k response has a 4xx status code
func (o *ShadowcopyModifyOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this shadowcopy modify o k response has a 5xx status code
func (o *ShadowcopyModifyOK) IsServerError() bool {
	return false
}

// IsCode returns true when this shadowcopy modify o k response a status code equal to that given
func (o *ShadowcopyModifyOK) IsCode(code int) bool {
	return code == 200
}

func (o *ShadowcopyModifyOK) Error() string {
	return fmt.Sprintf("[PATCH /protocols/cifs/shadow-copies/{client_uuid}][%d] shadowcopyModifyOK  %+v", 200, o.Payload)
}

func (o *ShadowcopyModifyOK) String() string {
	return fmt.Sprintf("[PATCH /protocols/cifs/shadow-copies/{client_uuid}][%d] shadowcopyModifyOK  %+v", 200, o.Payload)
}

func (o *ShadowcopyModifyOK) GetPayload() *models.ShadowcopyAddFiles {
	return o.Payload
}

func (o *ShadowcopyModifyOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ShadowcopyAddFiles)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewShadowcopyModifyDefault creates a ShadowcopyModifyDefault with default headers values
func NewShadowcopyModifyDefault(code int) *ShadowcopyModifyDefault {
	return &ShadowcopyModifyDefault{
		_statusCode: code,
	}
}

/*
ShadowcopyModifyDefault describes a response with status code -1, with default header values.

Error
*/
type ShadowcopyModifyDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the shadowcopy modify default response
func (o *ShadowcopyModifyDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this shadowcopy modify default response has a 2xx status code
func (o *ShadowcopyModifyDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this shadowcopy modify default response has a 3xx status code
func (o *ShadowcopyModifyDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this shadowcopy modify default response has a 4xx status code
func (o *ShadowcopyModifyDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this shadowcopy modify default response has a 5xx status code
func (o *ShadowcopyModifyDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this shadowcopy modify default response a status code equal to that given
func (o *ShadowcopyModifyDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *ShadowcopyModifyDefault) Error() string {
	return fmt.Sprintf("[PATCH /protocols/cifs/shadow-copies/{client_uuid}][%d] shadowcopy_modify default  %+v", o._statusCode, o.Payload)
}

func (o *ShadowcopyModifyDefault) String() string {
	return fmt.Sprintf("[PATCH /protocols/cifs/shadow-copies/{client_uuid}][%d] shadowcopy_modify default  %+v", o._statusCode, o.Payload)
}

func (o *ShadowcopyModifyDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ShadowcopyModifyDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}