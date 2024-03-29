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

// FileDirectorySecurityModifyReader is a Reader for the FileDirectorySecurityModify structure.
type FileDirectorySecurityModifyReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *FileDirectorySecurityModifyReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 202:
		result := NewFileDirectorySecurityModifyAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewFileDirectorySecurityModifyDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewFileDirectorySecurityModifyAccepted creates a FileDirectorySecurityModifyAccepted with default headers values
func NewFileDirectorySecurityModifyAccepted() *FileDirectorySecurityModifyAccepted {
	return &FileDirectorySecurityModifyAccepted{}
}

/*
FileDirectorySecurityModifyAccepted describes a response with status code 202, with default header values.

Accepted
*/
type FileDirectorySecurityModifyAccepted struct {
	Payload *models.JobLinkResponse
}

// IsSuccess returns true when this file directory security modify accepted response has a 2xx status code
func (o *FileDirectorySecurityModifyAccepted) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this file directory security modify accepted response has a 3xx status code
func (o *FileDirectorySecurityModifyAccepted) IsRedirect() bool {
	return false
}

// IsClientError returns true when this file directory security modify accepted response has a 4xx status code
func (o *FileDirectorySecurityModifyAccepted) IsClientError() bool {
	return false
}

// IsServerError returns true when this file directory security modify accepted response has a 5xx status code
func (o *FileDirectorySecurityModifyAccepted) IsServerError() bool {
	return false
}

// IsCode returns true when this file directory security modify accepted response a status code equal to that given
func (o *FileDirectorySecurityModifyAccepted) IsCode(code int) bool {
	return code == 202
}

func (o *FileDirectorySecurityModifyAccepted) Error() string {
	return fmt.Sprintf("[PATCH /protocols/file-security/permissions/{svm.uuid}/{path}][%d] fileDirectorySecurityModifyAccepted  %+v", 202, o.Payload)
}

func (o *FileDirectorySecurityModifyAccepted) String() string {
	return fmt.Sprintf("[PATCH /protocols/file-security/permissions/{svm.uuid}/{path}][%d] fileDirectorySecurityModifyAccepted  %+v", 202, o.Payload)
}

func (o *FileDirectorySecurityModifyAccepted) GetPayload() *models.JobLinkResponse {
	return o.Payload
}

func (o *FileDirectorySecurityModifyAccepted) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.JobLinkResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFileDirectorySecurityModifyDefault creates a FileDirectorySecurityModifyDefault with default headers values
func NewFileDirectorySecurityModifyDefault(code int) *FileDirectorySecurityModifyDefault {
	return &FileDirectorySecurityModifyDefault{
		_statusCode: code,
	}
}

/*
FileDirectorySecurityModifyDefault describes a response with status code -1, with default header values.

Error
*/
type FileDirectorySecurityModifyDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the file directory security modify default response
func (o *FileDirectorySecurityModifyDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this file directory security modify default response has a 2xx status code
func (o *FileDirectorySecurityModifyDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this file directory security modify default response has a 3xx status code
func (o *FileDirectorySecurityModifyDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this file directory security modify default response has a 4xx status code
func (o *FileDirectorySecurityModifyDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this file directory security modify default response has a 5xx status code
func (o *FileDirectorySecurityModifyDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this file directory security modify default response a status code equal to that given
func (o *FileDirectorySecurityModifyDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *FileDirectorySecurityModifyDefault) Error() string {
	return fmt.Sprintf("[PATCH /protocols/file-security/permissions/{svm.uuid}/{path}][%d] file_directory_security_modify default  %+v", o._statusCode, o.Payload)
}

func (o *FileDirectorySecurityModifyDefault) String() string {
	return fmt.Sprintf("[PATCH /protocols/file-security/permissions/{svm.uuid}/{path}][%d] file_directory_security_modify default  %+v", o._statusCode, o.Payload)
}

func (o *FileDirectorySecurityModifyDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *FileDirectorySecurityModifyDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
