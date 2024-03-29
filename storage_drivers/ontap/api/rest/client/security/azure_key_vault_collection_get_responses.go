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

// AzureKeyVaultCollectionGetReader is a Reader for the AzureKeyVaultCollectionGet structure.
type AzureKeyVaultCollectionGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AzureKeyVaultCollectionGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAzureKeyVaultCollectionGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewAzureKeyVaultCollectionGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAzureKeyVaultCollectionGetOK creates a AzureKeyVaultCollectionGetOK with default headers values
func NewAzureKeyVaultCollectionGetOK() *AzureKeyVaultCollectionGetOK {
	return &AzureKeyVaultCollectionGetOK{}
}

/*
AzureKeyVaultCollectionGetOK describes a response with status code 200, with default header values.

OK
*/
type AzureKeyVaultCollectionGetOK struct {
	Payload *models.AzureKeyVaultResponse
}

// IsSuccess returns true when this azure key vault collection get o k response has a 2xx status code
func (o *AzureKeyVaultCollectionGetOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this azure key vault collection get o k response has a 3xx status code
func (o *AzureKeyVaultCollectionGetOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this azure key vault collection get o k response has a 4xx status code
func (o *AzureKeyVaultCollectionGetOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this azure key vault collection get o k response has a 5xx status code
func (o *AzureKeyVaultCollectionGetOK) IsServerError() bool {
	return false
}

// IsCode returns true when this azure key vault collection get o k response a status code equal to that given
func (o *AzureKeyVaultCollectionGetOK) IsCode(code int) bool {
	return code == 200
}

func (o *AzureKeyVaultCollectionGetOK) Error() string {
	return fmt.Sprintf("[GET /security/azure-key-vaults][%d] azureKeyVaultCollectionGetOK  %+v", 200, o.Payload)
}

func (o *AzureKeyVaultCollectionGetOK) String() string {
	return fmt.Sprintf("[GET /security/azure-key-vaults][%d] azureKeyVaultCollectionGetOK  %+v", 200, o.Payload)
}

func (o *AzureKeyVaultCollectionGetOK) GetPayload() *models.AzureKeyVaultResponse {
	return o.Payload
}

func (o *AzureKeyVaultCollectionGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AzureKeyVaultResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAzureKeyVaultCollectionGetDefault creates a AzureKeyVaultCollectionGetDefault with default headers values
func NewAzureKeyVaultCollectionGetDefault(code int) *AzureKeyVaultCollectionGetDefault {
	return &AzureKeyVaultCollectionGetDefault{
		_statusCode: code,
	}
}

/*
AzureKeyVaultCollectionGetDefault describes a response with status code -1, with default header values.

Error
*/
type AzureKeyVaultCollectionGetDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the azure key vault collection get default response
func (o *AzureKeyVaultCollectionGetDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this azure key vault collection get default response has a 2xx status code
func (o *AzureKeyVaultCollectionGetDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this azure key vault collection get default response has a 3xx status code
func (o *AzureKeyVaultCollectionGetDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this azure key vault collection get default response has a 4xx status code
func (o *AzureKeyVaultCollectionGetDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this azure key vault collection get default response has a 5xx status code
func (o *AzureKeyVaultCollectionGetDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this azure key vault collection get default response a status code equal to that given
func (o *AzureKeyVaultCollectionGetDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *AzureKeyVaultCollectionGetDefault) Error() string {
	return fmt.Sprintf("[GET /security/azure-key-vaults][%d] azure_key_vault_collection_get default  %+v", o._statusCode, o.Payload)
}

func (o *AzureKeyVaultCollectionGetDefault) String() string {
	return fmt.Sprintf("[GET /security/azure-key-vaults][%d] azure_key_vault_collection_get default  %+v", o._statusCode, o.Payload)
}

func (o *AzureKeyVaultCollectionGetDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *AzureKeyVaultCollectionGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
