// Code generated by go-swagger; DO NOT EDIT.

package application

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netapp/trident/storage_drivers/ontap/api/rest/models"
)

// ConsistencyGroupDeleteReader is a Reader for the ConsistencyGroupDelete structure.
type ConsistencyGroupDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ConsistencyGroupDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewConsistencyGroupDeleteOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 202:
		result := NewConsistencyGroupDeleteAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewConsistencyGroupDeleteDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewConsistencyGroupDeleteOK creates a ConsistencyGroupDeleteOK with default headers values
func NewConsistencyGroupDeleteOK() *ConsistencyGroupDeleteOK {
	return &ConsistencyGroupDeleteOK{}
}

/*
ConsistencyGroupDeleteOK describes a response with status code 200, with default header values.

OK
*/
type ConsistencyGroupDeleteOK struct {
}

// IsSuccess returns true when this consistency group delete o k response has a 2xx status code
func (o *ConsistencyGroupDeleteOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this consistency group delete o k response has a 3xx status code
func (o *ConsistencyGroupDeleteOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this consistency group delete o k response has a 4xx status code
func (o *ConsistencyGroupDeleteOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this consistency group delete o k response has a 5xx status code
func (o *ConsistencyGroupDeleteOK) IsServerError() bool {
	return false
}

// IsCode returns true when this consistency group delete o k response a status code equal to that given
func (o *ConsistencyGroupDeleteOK) IsCode(code int) bool {
	return code == 200
}

func (o *ConsistencyGroupDeleteOK) Error() string {
	return fmt.Sprintf("[DELETE /application/consistency-groups/{uuid}][%d] consistencyGroupDeleteOK ", 200)
}

func (o *ConsistencyGroupDeleteOK) String() string {
	return fmt.Sprintf("[DELETE /application/consistency-groups/{uuid}][%d] consistencyGroupDeleteOK ", 200)
}

func (o *ConsistencyGroupDeleteOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewConsistencyGroupDeleteAccepted creates a ConsistencyGroupDeleteAccepted with default headers values
func NewConsistencyGroupDeleteAccepted() *ConsistencyGroupDeleteAccepted {
	return &ConsistencyGroupDeleteAccepted{}
}

/*
ConsistencyGroupDeleteAccepted describes a response with status code 202, with default header values.

Accepted
*/
type ConsistencyGroupDeleteAccepted struct {
	Payload *models.JobLinkResponse
}

// IsSuccess returns true when this consistency group delete accepted response has a 2xx status code
func (o *ConsistencyGroupDeleteAccepted) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this consistency group delete accepted response has a 3xx status code
func (o *ConsistencyGroupDeleteAccepted) IsRedirect() bool {
	return false
}

// IsClientError returns true when this consistency group delete accepted response has a 4xx status code
func (o *ConsistencyGroupDeleteAccepted) IsClientError() bool {
	return false
}

// IsServerError returns true when this consistency group delete accepted response has a 5xx status code
func (o *ConsistencyGroupDeleteAccepted) IsServerError() bool {
	return false
}

// IsCode returns true when this consistency group delete accepted response a status code equal to that given
func (o *ConsistencyGroupDeleteAccepted) IsCode(code int) bool {
	return code == 202
}

func (o *ConsistencyGroupDeleteAccepted) Error() string {
	return fmt.Sprintf("[DELETE /application/consistency-groups/{uuid}][%d] consistencyGroupDeleteAccepted  %+v", 202, o.Payload)
}

func (o *ConsistencyGroupDeleteAccepted) String() string {
	return fmt.Sprintf("[DELETE /application/consistency-groups/{uuid}][%d] consistencyGroupDeleteAccepted  %+v", 202, o.Payload)
}

func (o *ConsistencyGroupDeleteAccepted) GetPayload() *models.JobLinkResponse {
	return o.Payload
}

func (o *ConsistencyGroupDeleteAccepted) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.JobLinkResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewConsistencyGroupDeleteDefault creates a ConsistencyGroupDeleteDefault with default headers values
func NewConsistencyGroupDeleteDefault(code int) *ConsistencyGroupDeleteDefault {
	return &ConsistencyGroupDeleteDefault{
		_statusCode: code,
	}
}

/*
	ConsistencyGroupDeleteDefault describes a response with status code -1, with default header values.

	ONTAP Error Response Codes

| Error Code | Description |
| ---------- | ----------- |
| 53411842 | Consistency group does not exist. |
| 53411843 | A consistency group with specified UUID was not found. |
| 53411844 | Specified consistency group was not found in the specified SVM. |
| 53411845 | The specified UUID and name refer to different consistency groups. |
| 53411846 | Either name or UUID must be provided. |
*/
type ConsistencyGroupDeleteDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the consistency group delete default response
func (o *ConsistencyGroupDeleteDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this consistency group delete default response has a 2xx status code
func (o *ConsistencyGroupDeleteDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this consistency group delete default response has a 3xx status code
func (o *ConsistencyGroupDeleteDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this consistency group delete default response has a 4xx status code
func (o *ConsistencyGroupDeleteDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this consistency group delete default response has a 5xx status code
func (o *ConsistencyGroupDeleteDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this consistency group delete default response a status code equal to that given
func (o *ConsistencyGroupDeleteDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *ConsistencyGroupDeleteDefault) Error() string {
	return fmt.Sprintf("[DELETE /application/consistency-groups/{uuid}][%d] consistency_group_delete default  %+v", o._statusCode, o.Payload)
}

func (o *ConsistencyGroupDeleteDefault) String() string {
	return fmt.Sprintf("[DELETE /application/consistency-groups/{uuid}][%d] consistency_group_delete default  %+v", o._statusCode, o.Payload)
}

func (o *ConsistencyGroupDeleteDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ConsistencyGroupDeleteDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
