// Code generated by go-swagger; DO NOT EDIT.

package networking

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netapp/trident/storage_drivers/ontap/api/rest/models"
)

// NetworkIPServicePolicyGetReader is a Reader for the NetworkIPServicePolicyGet structure.
type NetworkIPServicePolicyGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *NetworkIPServicePolicyGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewNetworkIPServicePolicyGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewNetworkIPServicePolicyGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewNetworkIPServicePolicyGetOK creates a NetworkIPServicePolicyGetOK with default headers values
func NewNetworkIPServicePolicyGetOK() *NetworkIPServicePolicyGetOK {
	return &NetworkIPServicePolicyGetOK{}
}

/* NetworkIPServicePolicyGetOK describes a response with status code 200, with default header values.

OK
*/
type NetworkIPServicePolicyGetOK struct {
	Payload *models.IPServicePolicy
}

func (o *NetworkIPServicePolicyGetOK) Error() string {
	return fmt.Sprintf("[GET /network/ip/service-policies/{uuid}][%d] networkIpServicePolicyGetOK  %+v", 200, o.Payload)
}
func (o *NetworkIPServicePolicyGetOK) GetPayload() *models.IPServicePolicy {
	return o.Payload
}

func (o *NetworkIPServicePolicyGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.IPServicePolicy)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewNetworkIPServicePolicyGetDefault creates a NetworkIPServicePolicyGetDefault with default headers values
func NewNetworkIPServicePolicyGetDefault(code int) *NetworkIPServicePolicyGetDefault {
	return &NetworkIPServicePolicyGetDefault{
		_statusCode: code,
	}
}

/* NetworkIPServicePolicyGetDefault describes a response with status code -1, with default header values.

Error
*/
type NetworkIPServicePolicyGetDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the network ip service policy get default response
func (o *NetworkIPServicePolicyGetDefault) Code() int {
	return o._statusCode
}

func (o *NetworkIPServicePolicyGetDefault) Error() string {
	return fmt.Sprintf("[GET /network/ip/service-policies/{uuid}][%d] network_ip_service_policy_get default  %+v", o._statusCode, o.Payload)
}
func (o *NetworkIPServicePolicyGetDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *NetworkIPServicePolicyGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}