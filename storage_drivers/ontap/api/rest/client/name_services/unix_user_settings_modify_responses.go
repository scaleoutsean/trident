// Code generated by go-swagger; DO NOT EDIT.

package name_services

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netapp/trident/storage_drivers/ontap/api/rest/models"
)

// UnixUserSettingsModifyReader is a Reader for the UnixUserSettingsModify structure.
type UnixUserSettingsModifyReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UnixUserSettingsModifyReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUnixUserSettingsModifyOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewUnixUserSettingsModifyDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUnixUserSettingsModifyOK creates a UnixUserSettingsModifyOK with default headers values
func NewUnixUserSettingsModifyOK() *UnixUserSettingsModifyOK {
	return &UnixUserSettingsModifyOK{}
}

/* UnixUserSettingsModifyOK describes a response with status code 200, with default header values.

OK
*/
type UnixUserSettingsModifyOK struct {
}

func (o *UnixUserSettingsModifyOK) Error() string {
	return fmt.Sprintf("[PATCH /name-services/cache/unix-user/settings/{svm.uuid}][%d] unixUserSettingsModifyOK ", 200)
}

func (o *UnixUserSettingsModifyOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUnixUserSettingsModifyDefault creates a UnixUserSettingsModifyDefault with default headers values
func NewUnixUserSettingsModifyDefault(code int) *UnixUserSettingsModifyDefault {
	return &UnixUserSettingsModifyDefault{
		_statusCode: code,
	}
}

/* UnixUserSettingsModifyDefault describes a response with status code -1, with default header values.

 ONTAP Error Response Codes
| Error Code | Description |
| ---------- | ----------- |
| 23724055 | Internal error. Configuration for Vserver failed. Verify that the cluster is healthy, then try the command again. For further assistance, contact technical support. |

*/
type UnixUserSettingsModifyDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the unix user settings modify default response
func (o *UnixUserSettingsModifyDefault) Code() int {
	return o._statusCode
}

func (o *UnixUserSettingsModifyDefault) Error() string {
	return fmt.Sprintf("[PATCH /name-services/cache/unix-user/settings/{svm.uuid}][%d] unix_user_settings_modify default  %+v", o._statusCode, o.Payload)
}
func (o *UnixUserSettingsModifyDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UnixUserSettingsModifyDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}