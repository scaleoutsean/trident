// Code generated by go-swagger; DO NOT EDIT.

package security

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netapp/trident/storage_drivers/ontap/api/rest/models"
)

// AzureKeyVaultGetReader is a Reader for the AzureKeyVaultGet structure.
type AzureKeyVaultGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AzureKeyVaultGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAzureKeyVaultGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewAzureKeyVaultGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAzureKeyVaultGetOK creates a AzureKeyVaultGetOK with default headers values
func NewAzureKeyVaultGetOK() *AzureKeyVaultGetOK {
	return &AzureKeyVaultGetOK{}
}

/*
AzureKeyVaultGetOK describes a response with status code 200, with default header values.

OK
*/
type AzureKeyVaultGetOK struct {
	Payload *models.AzureKeyVault
}

// IsSuccess returns true when this azure key vault get o k response has a 2xx status code
func (o *AzureKeyVaultGetOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this azure key vault get o k response has a 3xx status code
func (o *AzureKeyVaultGetOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this azure key vault get o k response has a 4xx status code
func (o *AzureKeyVaultGetOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this azure key vault get o k response has a 5xx status code
func (o *AzureKeyVaultGetOK) IsServerError() bool {
	return false
}

// IsCode returns true when this azure key vault get o k response a status code equal to that given
func (o *AzureKeyVaultGetOK) IsCode(code int) bool {
	return code == 200
}

func (o *AzureKeyVaultGetOK) Error() string {
	return fmt.Sprintf("[GET /security/azure-key-vaults/{uuid}][%d] azureKeyVaultGetOK  %+v", 200, o.Payload)
}

func (o *AzureKeyVaultGetOK) String() string {
	return fmt.Sprintf("[GET /security/azure-key-vaults/{uuid}][%d] azureKeyVaultGetOK  %+v", 200, o.Payload)
}

func (o *AzureKeyVaultGetOK) GetPayload() *models.AzureKeyVault {
	return o.Payload
}

func (o *AzureKeyVaultGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AzureKeyVault)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAzureKeyVaultGetDefault creates a AzureKeyVaultGetDefault with default headers values
func NewAzureKeyVaultGetDefault(code int) *AzureKeyVaultGetDefault {
	return &AzureKeyVaultGetDefault{
		_statusCode: code,
	}
}

/*
AzureKeyVaultGetDefault describes a response with status code -1, with default header values.

Error
*/
type AzureKeyVaultGetDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the azure key vault get default response
func (o *AzureKeyVaultGetDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this azure key vault get default response has a 2xx status code
func (o *AzureKeyVaultGetDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this azure key vault get default response has a 3xx status code
func (o *AzureKeyVaultGetDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this azure key vault get default response has a 4xx status code
func (o *AzureKeyVaultGetDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this azure key vault get default response has a 5xx status code
func (o *AzureKeyVaultGetDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this azure key vault get default response a status code equal to that given
func (o *AzureKeyVaultGetDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *AzureKeyVaultGetDefault) Error() string {
	return fmt.Sprintf("[GET /security/azure-key-vaults/{uuid}][%d] azure_key_vault_get default  %+v", o._statusCode, o.Payload)
}

func (o *AzureKeyVaultGetDefault) String() string {
	return fmt.Sprintf("[GET /security/azure-key-vaults/{uuid}][%d] azure_key_vault_get default  %+v", o._statusCode, o.Payload)
}

func (o *AzureKeyVaultGetDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *AzureKeyVaultGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
