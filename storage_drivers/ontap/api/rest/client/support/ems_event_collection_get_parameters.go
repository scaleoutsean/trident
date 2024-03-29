// Code generated by go-swagger; DO NOT EDIT.

package support

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewEmsEventCollectionGetParams creates a new EmsEventCollectionGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewEmsEventCollectionGetParams() *EmsEventCollectionGetParams {
	return &EmsEventCollectionGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewEmsEventCollectionGetParamsWithTimeout creates a new EmsEventCollectionGetParams object
// with the ability to set a timeout on a request.
func NewEmsEventCollectionGetParamsWithTimeout(timeout time.Duration) *EmsEventCollectionGetParams {
	return &EmsEventCollectionGetParams{
		timeout: timeout,
	}
}

// NewEmsEventCollectionGetParamsWithContext creates a new EmsEventCollectionGetParams object
// with the ability to set a context for a request.
func NewEmsEventCollectionGetParamsWithContext(ctx context.Context) *EmsEventCollectionGetParams {
	return &EmsEventCollectionGetParams{
		Context: ctx,
	}
}

// NewEmsEventCollectionGetParamsWithHTTPClient creates a new EmsEventCollectionGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewEmsEventCollectionGetParamsWithHTTPClient(client *http.Client) *EmsEventCollectionGetParams {
	return &EmsEventCollectionGetParams{
		HTTPClient: client,
	}
}

/*
EmsEventCollectionGetParams contains all the parameters to send to the API endpoint

	for the ems event collection get operation.

	Typically these are written to a http.Request.
*/
type EmsEventCollectionGetParams struct {

	/* ActionPossibleActionsAction.

	   Filter by action.possible_actions.action
	*/
	ActionPossibleActionsActionQueryParameter *string

	/* ActionPossibleActionsInvokeVerb.

	   Filter by action.possible_actions.invoke.verb
	*/
	ActionPossibleActionsInvokeVerbQueryParameter *string

	/* ActionPossibleActionsParametersFormat.

	   Filter by action.possible_actions.parameters.format
	*/
	ActionPossibleActionsParametersFormatQueryParameter *string

	/* ActionPossibleActionsParametersName.

	   Filter by action.possible_actions.parameters.name
	*/
	ActionPossibleActionsParametersNameQueryParameter *string

	/* ActionPossibleActionsParametersType.

	   Filter by action.possible_actions.parameters.type
	*/
	ActionPossibleActionsParametersTypeQueryParameter *string

	/* CreationTime.

	   Filter by creation_time
	*/
	CreationTimeQueryParameter *string

	/* Fields.

	   Specify the fields to return.
	*/
	FieldsQueryParameter []string

	/* FilterName.

	   Filter the collection returned using an event filter
	*/
	FilterNameQueryParameter *string

	/* Index.

	   Filter by index
	*/
	IndexQueryParameter *int64

	/* LastUpdateTime.

	   Filter by last_update_time
	*/
	LastUpdateTimeQueryParameter *string

	/* LogMessage.

	   Filter by log_message
	*/
	LogMessageQueryParameter *string

	/* MaxRecords.

	   Limit the number of records returned.
	*/
	MaxRecordsQueryParameter *int64

	/* MessageName.

	   Filter by message.name
	*/
	MessageNameQueryParameter *string

	/* MessageSeverity.

	   Filter by message.severity
	*/
	MessageSeverityQueryParameter *string

	/* NodeName.

	   Filter by node.name
	*/
	NodeNameQueryParameter *string

	/* NodeUUID.

	   Filter by node.uuid
	*/
	NodeUUIDQueryParameter *string

	/* OrderBy.

	   Order results by specified fields and optional [asc|desc] direction. Default direction is 'asc' for ascending.
	*/
	OrderByQueryParameter []string

	/* ParametersName.

	   Filter by parameters.name
	*/
	ParametersNameQueryParameter *string

	/* ParametersValue.

	   Filter by parameters.value
	*/
	ParametersValueQueryParameter *string

	/* ReturnRecords.

	   The default is true for GET calls.  When set to false, only the number of records is returned.

	   Default: true
	*/
	ReturnRecordsQueryParameter *bool

	/* ReturnTimeout.

	   The number of seconds to allow the call to execute before returning.  When iterating over a collection, the default is 15 seconds.  ONTAP returns earlier if either max records or the end of the collection is reached.

	   Default: 15
	*/
	ReturnTimeoutQueryParameter *int64

	/* Source.

	   Filter by source
	*/
	SourceQueryParameter *string

	/* State.

	   Filter by state
	*/
	StateQueryParameter *string

	/* Stateful.

	   Filter by stateful
	*/
	StatefulQueryParameter *bool

	/* Time.

	   Filter by time
	*/
	TimeQueryParameter *string

	/* UUID.

	   Filter by uuid
	*/
	UUIDQueryParameter *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the ems event collection get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *EmsEventCollectionGetParams) WithDefaults() *EmsEventCollectionGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the ems event collection get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *EmsEventCollectionGetParams) SetDefaults() {
	var (
		returnRecordsQueryParameterDefault = bool(true)

		returnTimeoutQueryParameterDefault = int64(15)
	)

	val := EmsEventCollectionGetParams{
		ReturnRecordsQueryParameter: &returnRecordsQueryParameterDefault,
		ReturnTimeoutQueryParameter: &returnTimeoutQueryParameterDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithTimeout(timeout time.Duration) *EmsEventCollectionGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithContext(ctx context.Context) *EmsEventCollectionGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithHTTPClient(client *http.Client) *EmsEventCollectionGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithActionPossibleActionsActionQueryParameter adds the actionPossibleActionsAction to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithActionPossibleActionsActionQueryParameter(actionPossibleActionsAction *string) *EmsEventCollectionGetParams {
	o.SetActionPossibleActionsActionQueryParameter(actionPossibleActionsAction)
	return o
}

// SetActionPossibleActionsActionQueryParameter adds the actionPossibleActionsAction to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetActionPossibleActionsActionQueryParameter(actionPossibleActionsAction *string) {
	o.ActionPossibleActionsActionQueryParameter = actionPossibleActionsAction
}

// WithActionPossibleActionsInvokeVerbQueryParameter adds the actionPossibleActionsInvokeVerb to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithActionPossibleActionsInvokeVerbQueryParameter(actionPossibleActionsInvokeVerb *string) *EmsEventCollectionGetParams {
	o.SetActionPossibleActionsInvokeVerbQueryParameter(actionPossibleActionsInvokeVerb)
	return o
}

// SetActionPossibleActionsInvokeVerbQueryParameter adds the actionPossibleActionsInvokeVerb to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetActionPossibleActionsInvokeVerbQueryParameter(actionPossibleActionsInvokeVerb *string) {
	o.ActionPossibleActionsInvokeVerbQueryParameter = actionPossibleActionsInvokeVerb
}

// WithActionPossibleActionsParametersFormatQueryParameter adds the actionPossibleActionsParametersFormat to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithActionPossibleActionsParametersFormatQueryParameter(actionPossibleActionsParametersFormat *string) *EmsEventCollectionGetParams {
	o.SetActionPossibleActionsParametersFormatQueryParameter(actionPossibleActionsParametersFormat)
	return o
}

// SetActionPossibleActionsParametersFormatQueryParameter adds the actionPossibleActionsParametersFormat to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetActionPossibleActionsParametersFormatQueryParameter(actionPossibleActionsParametersFormat *string) {
	o.ActionPossibleActionsParametersFormatQueryParameter = actionPossibleActionsParametersFormat
}

// WithActionPossibleActionsParametersNameQueryParameter adds the actionPossibleActionsParametersName to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithActionPossibleActionsParametersNameQueryParameter(actionPossibleActionsParametersName *string) *EmsEventCollectionGetParams {
	o.SetActionPossibleActionsParametersNameQueryParameter(actionPossibleActionsParametersName)
	return o
}

// SetActionPossibleActionsParametersNameQueryParameter adds the actionPossibleActionsParametersName to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetActionPossibleActionsParametersNameQueryParameter(actionPossibleActionsParametersName *string) {
	o.ActionPossibleActionsParametersNameQueryParameter = actionPossibleActionsParametersName
}

// WithActionPossibleActionsParametersTypeQueryParameter adds the actionPossibleActionsParametersType to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithActionPossibleActionsParametersTypeQueryParameter(actionPossibleActionsParametersType *string) *EmsEventCollectionGetParams {
	o.SetActionPossibleActionsParametersTypeQueryParameter(actionPossibleActionsParametersType)
	return o
}

// SetActionPossibleActionsParametersTypeQueryParameter adds the actionPossibleActionsParametersType to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetActionPossibleActionsParametersTypeQueryParameter(actionPossibleActionsParametersType *string) {
	o.ActionPossibleActionsParametersTypeQueryParameter = actionPossibleActionsParametersType
}

// WithCreationTimeQueryParameter adds the creationTime to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithCreationTimeQueryParameter(creationTime *string) *EmsEventCollectionGetParams {
	o.SetCreationTimeQueryParameter(creationTime)
	return o
}

// SetCreationTimeQueryParameter adds the creationTime to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetCreationTimeQueryParameter(creationTime *string) {
	o.CreationTimeQueryParameter = creationTime
}

// WithFieldsQueryParameter adds the fields to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithFieldsQueryParameter(fields []string) *EmsEventCollectionGetParams {
	o.SetFieldsQueryParameter(fields)
	return o
}

// SetFieldsQueryParameter adds the fields to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetFieldsQueryParameter(fields []string) {
	o.FieldsQueryParameter = fields
}

// WithFilterNameQueryParameter adds the filterName to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithFilterNameQueryParameter(filterName *string) *EmsEventCollectionGetParams {
	o.SetFilterNameQueryParameter(filterName)
	return o
}

// SetFilterNameQueryParameter adds the filterName to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetFilterNameQueryParameter(filterName *string) {
	o.FilterNameQueryParameter = filterName
}

// WithIndexQueryParameter adds the index to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithIndexQueryParameter(index *int64) *EmsEventCollectionGetParams {
	o.SetIndexQueryParameter(index)
	return o
}

// SetIndexQueryParameter adds the index to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetIndexQueryParameter(index *int64) {
	o.IndexQueryParameter = index
}

// WithLastUpdateTimeQueryParameter adds the lastUpdateTime to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithLastUpdateTimeQueryParameter(lastUpdateTime *string) *EmsEventCollectionGetParams {
	o.SetLastUpdateTimeQueryParameter(lastUpdateTime)
	return o
}

// SetLastUpdateTimeQueryParameter adds the lastUpdateTime to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetLastUpdateTimeQueryParameter(lastUpdateTime *string) {
	o.LastUpdateTimeQueryParameter = lastUpdateTime
}

// WithLogMessageQueryParameter adds the logMessage to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithLogMessageQueryParameter(logMessage *string) *EmsEventCollectionGetParams {
	o.SetLogMessageQueryParameter(logMessage)
	return o
}

// SetLogMessageQueryParameter adds the logMessage to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetLogMessageQueryParameter(logMessage *string) {
	o.LogMessageQueryParameter = logMessage
}

// WithMaxRecordsQueryParameter adds the maxRecords to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithMaxRecordsQueryParameter(maxRecords *int64) *EmsEventCollectionGetParams {
	o.SetMaxRecordsQueryParameter(maxRecords)
	return o
}

// SetMaxRecordsQueryParameter adds the maxRecords to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetMaxRecordsQueryParameter(maxRecords *int64) {
	o.MaxRecordsQueryParameter = maxRecords
}

// WithMessageNameQueryParameter adds the messageName to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithMessageNameQueryParameter(messageName *string) *EmsEventCollectionGetParams {
	o.SetMessageNameQueryParameter(messageName)
	return o
}

// SetMessageNameQueryParameter adds the messageName to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetMessageNameQueryParameter(messageName *string) {
	o.MessageNameQueryParameter = messageName
}

// WithMessageSeverityQueryParameter adds the messageSeverity to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithMessageSeverityQueryParameter(messageSeverity *string) *EmsEventCollectionGetParams {
	o.SetMessageSeverityQueryParameter(messageSeverity)
	return o
}

// SetMessageSeverityQueryParameter adds the messageSeverity to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetMessageSeverityQueryParameter(messageSeverity *string) {
	o.MessageSeverityQueryParameter = messageSeverity
}

// WithNodeNameQueryParameter adds the nodeName to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithNodeNameQueryParameter(nodeName *string) *EmsEventCollectionGetParams {
	o.SetNodeNameQueryParameter(nodeName)
	return o
}

// SetNodeNameQueryParameter adds the nodeName to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetNodeNameQueryParameter(nodeName *string) {
	o.NodeNameQueryParameter = nodeName
}

// WithNodeUUIDQueryParameter adds the nodeUUID to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithNodeUUIDQueryParameter(nodeUUID *string) *EmsEventCollectionGetParams {
	o.SetNodeUUIDQueryParameter(nodeUUID)
	return o
}

// SetNodeUUIDQueryParameter adds the nodeUuid to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetNodeUUIDQueryParameter(nodeUUID *string) {
	o.NodeUUIDQueryParameter = nodeUUID
}

// WithOrderByQueryParameter adds the orderBy to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithOrderByQueryParameter(orderBy []string) *EmsEventCollectionGetParams {
	o.SetOrderByQueryParameter(orderBy)
	return o
}

// SetOrderByQueryParameter adds the orderBy to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetOrderByQueryParameter(orderBy []string) {
	o.OrderByQueryParameter = orderBy
}

// WithParametersNameQueryParameter adds the parametersName to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithParametersNameQueryParameter(parametersName *string) *EmsEventCollectionGetParams {
	o.SetParametersNameQueryParameter(parametersName)
	return o
}

// SetParametersNameQueryParameter adds the parametersName to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetParametersNameQueryParameter(parametersName *string) {
	o.ParametersNameQueryParameter = parametersName
}

// WithParametersValueQueryParameter adds the parametersValue to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithParametersValueQueryParameter(parametersValue *string) *EmsEventCollectionGetParams {
	o.SetParametersValueQueryParameter(parametersValue)
	return o
}

// SetParametersValueQueryParameter adds the parametersValue to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetParametersValueQueryParameter(parametersValue *string) {
	o.ParametersValueQueryParameter = parametersValue
}

// WithReturnRecordsQueryParameter adds the returnRecords to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithReturnRecordsQueryParameter(returnRecords *bool) *EmsEventCollectionGetParams {
	o.SetReturnRecordsQueryParameter(returnRecords)
	return o
}

// SetReturnRecordsQueryParameter adds the returnRecords to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetReturnRecordsQueryParameter(returnRecords *bool) {
	o.ReturnRecordsQueryParameter = returnRecords
}

// WithReturnTimeoutQueryParameter adds the returnTimeout to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithReturnTimeoutQueryParameter(returnTimeout *int64) *EmsEventCollectionGetParams {
	o.SetReturnTimeoutQueryParameter(returnTimeout)
	return o
}

// SetReturnTimeoutQueryParameter adds the returnTimeout to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetReturnTimeoutQueryParameter(returnTimeout *int64) {
	o.ReturnTimeoutQueryParameter = returnTimeout
}

// WithSourceQueryParameter adds the source to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithSourceQueryParameter(source *string) *EmsEventCollectionGetParams {
	o.SetSourceQueryParameter(source)
	return o
}

// SetSourceQueryParameter adds the source to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetSourceQueryParameter(source *string) {
	o.SourceQueryParameter = source
}

// WithStateQueryParameter adds the state to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithStateQueryParameter(state *string) *EmsEventCollectionGetParams {
	o.SetStateQueryParameter(state)
	return o
}

// SetStateQueryParameter adds the state to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetStateQueryParameter(state *string) {
	o.StateQueryParameter = state
}

// WithStatefulQueryParameter adds the stateful to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithStatefulQueryParameter(stateful *bool) *EmsEventCollectionGetParams {
	o.SetStatefulQueryParameter(stateful)
	return o
}

// SetStatefulQueryParameter adds the stateful to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetStatefulQueryParameter(stateful *bool) {
	o.StatefulQueryParameter = stateful
}

// WithTimeQueryParameter adds the time to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithTimeQueryParameter(time *string) *EmsEventCollectionGetParams {
	o.SetTimeQueryParameter(time)
	return o
}

// SetTimeQueryParameter adds the time to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetTimeQueryParameter(time *string) {
	o.TimeQueryParameter = time
}

// WithUUIDQueryParameter adds the uuid to the ems event collection get params
func (o *EmsEventCollectionGetParams) WithUUIDQueryParameter(uuid *string) *EmsEventCollectionGetParams {
	o.SetUUIDQueryParameter(uuid)
	return o
}

// SetUUIDQueryParameter adds the uuid to the ems event collection get params
func (o *EmsEventCollectionGetParams) SetUUIDQueryParameter(uuid *string) {
	o.UUIDQueryParameter = uuid
}

// WriteToRequest writes these params to a swagger request
func (o *EmsEventCollectionGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.ActionPossibleActionsActionQueryParameter != nil {

		// query param action.possible_actions.action
		var qrActionPossibleActionsAction string

		if o.ActionPossibleActionsActionQueryParameter != nil {
			qrActionPossibleActionsAction = *o.ActionPossibleActionsActionQueryParameter
		}
		qActionPossibleActionsAction := qrActionPossibleActionsAction
		if qActionPossibleActionsAction != "" {

			if err := r.SetQueryParam("action.possible_actions.action", qActionPossibleActionsAction); err != nil {
				return err
			}
		}
	}

	if o.ActionPossibleActionsInvokeVerbQueryParameter != nil {

		// query param action.possible_actions.invoke.verb
		var qrActionPossibleActionsInvokeVerb string

		if o.ActionPossibleActionsInvokeVerbQueryParameter != nil {
			qrActionPossibleActionsInvokeVerb = *o.ActionPossibleActionsInvokeVerbQueryParameter
		}
		qActionPossibleActionsInvokeVerb := qrActionPossibleActionsInvokeVerb
		if qActionPossibleActionsInvokeVerb != "" {

			if err := r.SetQueryParam("action.possible_actions.invoke.verb", qActionPossibleActionsInvokeVerb); err != nil {
				return err
			}
		}
	}

	if o.ActionPossibleActionsParametersFormatQueryParameter != nil {

		// query param action.possible_actions.parameters.format
		var qrActionPossibleActionsParametersFormat string

		if o.ActionPossibleActionsParametersFormatQueryParameter != nil {
			qrActionPossibleActionsParametersFormat = *o.ActionPossibleActionsParametersFormatQueryParameter
		}
		qActionPossibleActionsParametersFormat := qrActionPossibleActionsParametersFormat
		if qActionPossibleActionsParametersFormat != "" {

			if err := r.SetQueryParam("action.possible_actions.parameters.format", qActionPossibleActionsParametersFormat); err != nil {
				return err
			}
		}
	}

	if o.ActionPossibleActionsParametersNameQueryParameter != nil {

		// query param action.possible_actions.parameters.name
		var qrActionPossibleActionsParametersName string

		if o.ActionPossibleActionsParametersNameQueryParameter != nil {
			qrActionPossibleActionsParametersName = *o.ActionPossibleActionsParametersNameQueryParameter
		}
		qActionPossibleActionsParametersName := qrActionPossibleActionsParametersName
		if qActionPossibleActionsParametersName != "" {

			if err := r.SetQueryParam("action.possible_actions.parameters.name", qActionPossibleActionsParametersName); err != nil {
				return err
			}
		}
	}

	if o.ActionPossibleActionsParametersTypeQueryParameter != nil {

		// query param action.possible_actions.parameters.type
		var qrActionPossibleActionsParametersType string

		if o.ActionPossibleActionsParametersTypeQueryParameter != nil {
			qrActionPossibleActionsParametersType = *o.ActionPossibleActionsParametersTypeQueryParameter
		}
		qActionPossibleActionsParametersType := qrActionPossibleActionsParametersType
		if qActionPossibleActionsParametersType != "" {

			if err := r.SetQueryParam("action.possible_actions.parameters.type", qActionPossibleActionsParametersType); err != nil {
				return err
			}
		}
	}

	if o.CreationTimeQueryParameter != nil {

		// query param creation_time
		var qrCreationTime string

		if o.CreationTimeQueryParameter != nil {
			qrCreationTime = *o.CreationTimeQueryParameter
		}
		qCreationTime := qrCreationTime
		if qCreationTime != "" {

			if err := r.SetQueryParam("creation_time", qCreationTime); err != nil {
				return err
			}
		}
	}

	if o.FieldsQueryParameter != nil {

		// binding items for fields
		joinedFields := o.bindParamFields(reg)

		// query array param fields
		if err := r.SetQueryParam("fields", joinedFields...); err != nil {
			return err
		}
	}

	if o.FilterNameQueryParameter != nil {

		// query param filter.name
		var qrFilterName string

		if o.FilterNameQueryParameter != nil {
			qrFilterName = *o.FilterNameQueryParameter
		}
		qFilterName := qrFilterName
		if qFilterName != "" {

			if err := r.SetQueryParam("filter.name", qFilterName); err != nil {
				return err
			}
		}
	}

	if o.IndexQueryParameter != nil {

		// query param index
		var qrIndex int64

		if o.IndexQueryParameter != nil {
			qrIndex = *o.IndexQueryParameter
		}
		qIndex := swag.FormatInt64(qrIndex)
		if qIndex != "" {

			if err := r.SetQueryParam("index", qIndex); err != nil {
				return err
			}
		}
	}

	if o.LastUpdateTimeQueryParameter != nil {

		// query param last_update_time
		var qrLastUpdateTime string

		if o.LastUpdateTimeQueryParameter != nil {
			qrLastUpdateTime = *o.LastUpdateTimeQueryParameter
		}
		qLastUpdateTime := qrLastUpdateTime
		if qLastUpdateTime != "" {

			if err := r.SetQueryParam("last_update_time", qLastUpdateTime); err != nil {
				return err
			}
		}
	}

	if o.LogMessageQueryParameter != nil {

		// query param log_message
		var qrLogMessage string

		if o.LogMessageQueryParameter != nil {
			qrLogMessage = *o.LogMessageQueryParameter
		}
		qLogMessage := qrLogMessage
		if qLogMessage != "" {

			if err := r.SetQueryParam("log_message", qLogMessage); err != nil {
				return err
			}
		}
	}

	if o.MaxRecordsQueryParameter != nil {

		// query param max_records
		var qrMaxRecords int64

		if o.MaxRecordsQueryParameter != nil {
			qrMaxRecords = *o.MaxRecordsQueryParameter
		}
		qMaxRecords := swag.FormatInt64(qrMaxRecords)
		if qMaxRecords != "" {

			if err := r.SetQueryParam("max_records", qMaxRecords); err != nil {
				return err
			}
		}
	}

	if o.MessageNameQueryParameter != nil {

		// query param message.name
		var qrMessageName string

		if o.MessageNameQueryParameter != nil {
			qrMessageName = *o.MessageNameQueryParameter
		}
		qMessageName := qrMessageName
		if qMessageName != "" {

			if err := r.SetQueryParam("message.name", qMessageName); err != nil {
				return err
			}
		}
	}

	if o.MessageSeverityQueryParameter != nil {

		// query param message.severity
		var qrMessageSeverity string

		if o.MessageSeverityQueryParameter != nil {
			qrMessageSeverity = *o.MessageSeverityQueryParameter
		}
		qMessageSeverity := qrMessageSeverity
		if qMessageSeverity != "" {

			if err := r.SetQueryParam("message.severity", qMessageSeverity); err != nil {
				return err
			}
		}
	}

	if o.NodeNameQueryParameter != nil {

		// query param node.name
		var qrNodeName string

		if o.NodeNameQueryParameter != nil {
			qrNodeName = *o.NodeNameQueryParameter
		}
		qNodeName := qrNodeName
		if qNodeName != "" {

			if err := r.SetQueryParam("node.name", qNodeName); err != nil {
				return err
			}
		}
	}

	if o.NodeUUIDQueryParameter != nil {

		// query param node.uuid
		var qrNodeUUID string

		if o.NodeUUIDQueryParameter != nil {
			qrNodeUUID = *o.NodeUUIDQueryParameter
		}
		qNodeUUID := qrNodeUUID
		if qNodeUUID != "" {

			if err := r.SetQueryParam("node.uuid", qNodeUUID); err != nil {
				return err
			}
		}
	}

	if o.OrderByQueryParameter != nil {

		// binding items for order_by
		joinedOrderBy := o.bindParamOrderBy(reg)

		// query array param order_by
		if err := r.SetQueryParam("order_by", joinedOrderBy...); err != nil {
			return err
		}
	}

	if o.ParametersNameQueryParameter != nil {

		// query param parameters.name
		var qrParametersName string

		if o.ParametersNameQueryParameter != nil {
			qrParametersName = *o.ParametersNameQueryParameter
		}
		qParametersName := qrParametersName
		if qParametersName != "" {

			if err := r.SetQueryParam("parameters.name", qParametersName); err != nil {
				return err
			}
		}
	}

	if o.ParametersValueQueryParameter != nil {

		// query param parameters.value
		var qrParametersValue string

		if o.ParametersValueQueryParameter != nil {
			qrParametersValue = *o.ParametersValueQueryParameter
		}
		qParametersValue := qrParametersValue
		if qParametersValue != "" {

			if err := r.SetQueryParam("parameters.value", qParametersValue); err != nil {
				return err
			}
		}
	}

	if o.ReturnRecordsQueryParameter != nil {

		// query param return_records
		var qrReturnRecords bool

		if o.ReturnRecordsQueryParameter != nil {
			qrReturnRecords = *o.ReturnRecordsQueryParameter
		}
		qReturnRecords := swag.FormatBool(qrReturnRecords)
		if qReturnRecords != "" {

			if err := r.SetQueryParam("return_records", qReturnRecords); err != nil {
				return err
			}
		}
	}

	if o.ReturnTimeoutQueryParameter != nil {

		// query param return_timeout
		var qrReturnTimeout int64

		if o.ReturnTimeoutQueryParameter != nil {
			qrReturnTimeout = *o.ReturnTimeoutQueryParameter
		}
		qReturnTimeout := swag.FormatInt64(qrReturnTimeout)
		if qReturnTimeout != "" {

			if err := r.SetQueryParam("return_timeout", qReturnTimeout); err != nil {
				return err
			}
		}
	}

	if o.SourceQueryParameter != nil {

		// query param source
		var qrSource string

		if o.SourceQueryParameter != nil {
			qrSource = *o.SourceQueryParameter
		}
		qSource := qrSource
		if qSource != "" {

			if err := r.SetQueryParam("source", qSource); err != nil {
				return err
			}
		}
	}

	if o.StateQueryParameter != nil {

		// query param state
		var qrState string

		if o.StateQueryParameter != nil {
			qrState = *o.StateQueryParameter
		}
		qState := qrState
		if qState != "" {

			if err := r.SetQueryParam("state", qState); err != nil {
				return err
			}
		}
	}

	if o.StatefulQueryParameter != nil {

		// query param stateful
		var qrStateful bool

		if o.StatefulQueryParameter != nil {
			qrStateful = *o.StatefulQueryParameter
		}
		qStateful := swag.FormatBool(qrStateful)
		if qStateful != "" {

			if err := r.SetQueryParam("stateful", qStateful); err != nil {
				return err
			}
		}
	}

	if o.TimeQueryParameter != nil {

		// query param time
		var qrTime string

		if o.TimeQueryParameter != nil {
			qrTime = *o.TimeQueryParameter
		}
		qTime := qrTime
		if qTime != "" {

			if err := r.SetQueryParam("time", qTime); err != nil {
				return err
			}
		}
	}

	if o.UUIDQueryParameter != nil {

		// query param uuid
		var qrUUID string

		if o.UUIDQueryParameter != nil {
			qrUUID = *o.UUIDQueryParameter
		}
		qUUID := qrUUID
		if qUUID != "" {

			if err := r.SetQueryParam("uuid", qUUID); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindParamEmsEventCollectionGet binds the parameter fields
func (o *EmsEventCollectionGetParams) bindParamFields(formats strfmt.Registry) []string {
	fieldsIR := o.FieldsQueryParameter

	var fieldsIC []string
	for _, fieldsIIR := range fieldsIR { // explode []string

		fieldsIIV := fieldsIIR // string as string
		fieldsIC = append(fieldsIC, fieldsIIV)
	}

	// items.CollectionFormat: "csv"
	fieldsIS := swag.JoinByFormat(fieldsIC, "csv")

	return fieldsIS
}

// bindParamEmsEventCollectionGet binds the parameter order_by
func (o *EmsEventCollectionGetParams) bindParamOrderBy(formats strfmt.Registry) []string {
	orderByIR := o.OrderByQueryParameter

	var orderByIC []string
	for _, orderByIIR := range orderByIR { // explode []string

		orderByIIV := orderByIIR // string as string
		orderByIC = append(orderByIC, orderByIIV)
	}

	// items.CollectionFormat: "csv"
	orderByIS := swag.JoinByFormat(orderByIC, "csv")

	return orderByIS
}
