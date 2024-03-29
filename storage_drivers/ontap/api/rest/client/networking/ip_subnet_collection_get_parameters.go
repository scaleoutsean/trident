// Code generated by go-swagger; DO NOT EDIT.

package networking

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

// NewIPSubnetCollectionGetParams creates a new IPSubnetCollectionGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewIPSubnetCollectionGetParams() *IPSubnetCollectionGetParams {
	return &IPSubnetCollectionGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewIPSubnetCollectionGetParamsWithTimeout creates a new IPSubnetCollectionGetParams object
// with the ability to set a timeout on a request.
func NewIPSubnetCollectionGetParamsWithTimeout(timeout time.Duration) *IPSubnetCollectionGetParams {
	return &IPSubnetCollectionGetParams{
		timeout: timeout,
	}
}

// NewIPSubnetCollectionGetParamsWithContext creates a new IPSubnetCollectionGetParams object
// with the ability to set a context for a request.
func NewIPSubnetCollectionGetParamsWithContext(ctx context.Context) *IPSubnetCollectionGetParams {
	return &IPSubnetCollectionGetParams{
		Context: ctx,
	}
}

// NewIPSubnetCollectionGetParamsWithHTTPClient creates a new IPSubnetCollectionGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewIPSubnetCollectionGetParamsWithHTTPClient(client *http.Client) *IPSubnetCollectionGetParams {
	return &IPSubnetCollectionGetParams{
		HTTPClient: client,
	}
}

/*
IPSubnetCollectionGetParams contains all the parameters to send to the API endpoint

	for the ip subnet collection get operation.

	Typically these are written to a http.Request.
*/
type IPSubnetCollectionGetParams struct {

	/* AvailableCount.

	   Filter by available_count
	*/
	AvailableCountQueryParameter *int64

	/* AvailableIPRangesEnd.

	   Filter by available_ip_ranges.end
	*/
	AvailableIPRangesEndQueryParameter *string

	/* AvailableIPRangesFamily.

	   Filter by available_ip_ranges.family
	*/
	AvailableIPRangesFamilyQueryParameter *string

	/* AvailableIPRangesStart.

	   Filter by available_ip_ranges.start
	*/
	AvailableIPRangesStartQueryParameter *string

	/* BroadcastDomainName.

	   Filter by broadcast_domain.name
	*/
	BroadcastDomainNameQueryParameter *string

	/* BroadcastDomainUUID.

	   Filter by broadcast_domain.uuid
	*/
	BroadcastDomainUUIDQueryParameter *string

	/* Fields.

	   Specify the fields to return.
	*/
	FieldsQueryParameter []string

	/* Gateway.

	   Filter by gateway
	*/
	GatewayQueryParameter *string

	/* IPRangesEnd.

	   Filter by ip_ranges.end
	*/
	IPRangesEndQueryParameter *string

	/* IPRangesFamily.

	   Filter by ip_ranges.family
	*/
	IPRangesFamilyQueryParameter *string

	/* IPRangesStart.

	   Filter by ip_ranges.start
	*/
	IPRangesStartQueryParameter *string

	/* IpspaceName.

	   Filter by ipspace.name
	*/
	IpspaceNameQueryParameter *string

	/* IpspaceUUID.

	   Filter by ipspace.uuid
	*/
	IpspaceUUIDQueryParameter *string

	/* MaxRecords.

	   Limit the number of records returned.
	*/
	MaxRecordsQueryParameter *int64

	/* Name.

	   Filter by name
	*/
	NameQueryParameter *string

	/* OrderBy.

	   Order results by specified fields and optional [asc|desc] direction. Default direction is 'asc' for ascending.
	*/
	OrderByQueryParameter []string

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

	/* SubnetAddress.

	   Filter by subnet.address
	*/
	SubnetAddressQueryParameter *string

	/* SubnetFamily.

	   Filter by subnet.family
	*/
	SubnetFamilyQueryParameter *string

	/* SubnetNetmask.

	   Filter by subnet.netmask
	*/
	SubnetNetmaskQueryParameter *string

	/* TotalCount.

	   Filter by total_count
	*/
	TotalCountQueryParameter *int64

	/* UsedCount.

	   Filter by used_count
	*/
	UsedCountQueryParameter *int64

	/* UUID.

	   Filter by uuid
	*/
	UUIDQueryParameter *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the ip subnet collection get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *IPSubnetCollectionGetParams) WithDefaults() *IPSubnetCollectionGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the ip subnet collection get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *IPSubnetCollectionGetParams) SetDefaults() {
	var (
		returnRecordsQueryParameterDefault = bool(true)

		returnTimeoutQueryParameterDefault = int64(15)
	)

	val := IPSubnetCollectionGetParams{
		ReturnRecordsQueryParameter: &returnRecordsQueryParameterDefault,
		ReturnTimeoutQueryParameter: &returnTimeoutQueryParameterDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithTimeout(timeout time.Duration) *IPSubnetCollectionGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithContext(ctx context.Context) *IPSubnetCollectionGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithHTTPClient(client *http.Client) *IPSubnetCollectionGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAvailableCountQueryParameter adds the availableCount to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithAvailableCountQueryParameter(availableCount *int64) *IPSubnetCollectionGetParams {
	o.SetAvailableCountQueryParameter(availableCount)
	return o
}

// SetAvailableCountQueryParameter adds the availableCount to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetAvailableCountQueryParameter(availableCount *int64) {
	o.AvailableCountQueryParameter = availableCount
}

// WithAvailableIPRangesEndQueryParameter adds the availableIPRangesEnd to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithAvailableIPRangesEndQueryParameter(availableIPRangesEnd *string) *IPSubnetCollectionGetParams {
	o.SetAvailableIPRangesEndQueryParameter(availableIPRangesEnd)
	return o
}

// SetAvailableIPRangesEndQueryParameter adds the availableIpRangesEnd to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetAvailableIPRangesEndQueryParameter(availableIPRangesEnd *string) {
	o.AvailableIPRangesEndQueryParameter = availableIPRangesEnd
}

// WithAvailableIPRangesFamilyQueryParameter adds the availableIPRangesFamily to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithAvailableIPRangesFamilyQueryParameter(availableIPRangesFamily *string) *IPSubnetCollectionGetParams {
	o.SetAvailableIPRangesFamilyQueryParameter(availableIPRangesFamily)
	return o
}

// SetAvailableIPRangesFamilyQueryParameter adds the availableIpRangesFamily to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetAvailableIPRangesFamilyQueryParameter(availableIPRangesFamily *string) {
	o.AvailableIPRangesFamilyQueryParameter = availableIPRangesFamily
}

// WithAvailableIPRangesStartQueryParameter adds the availableIPRangesStart to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithAvailableIPRangesStartQueryParameter(availableIPRangesStart *string) *IPSubnetCollectionGetParams {
	o.SetAvailableIPRangesStartQueryParameter(availableIPRangesStart)
	return o
}

// SetAvailableIPRangesStartQueryParameter adds the availableIpRangesStart to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetAvailableIPRangesStartQueryParameter(availableIPRangesStart *string) {
	o.AvailableIPRangesStartQueryParameter = availableIPRangesStart
}

// WithBroadcastDomainNameQueryParameter adds the broadcastDomainName to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithBroadcastDomainNameQueryParameter(broadcastDomainName *string) *IPSubnetCollectionGetParams {
	o.SetBroadcastDomainNameQueryParameter(broadcastDomainName)
	return o
}

// SetBroadcastDomainNameQueryParameter adds the broadcastDomainName to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetBroadcastDomainNameQueryParameter(broadcastDomainName *string) {
	o.BroadcastDomainNameQueryParameter = broadcastDomainName
}

// WithBroadcastDomainUUIDQueryParameter adds the broadcastDomainUUID to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithBroadcastDomainUUIDQueryParameter(broadcastDomainUUID *string) *IPSubnetCollectionGetParams {
	o.SetBroadcastDomainUUIDQueryParameter(broadcastDomainUUID)
	return o
}

// SetBroadcastDomainUUIDQueryParameter adds the broadcastDomainUuid to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetBroadcastDomainUUIDQueryParameter(broadcastDomainUUID *string) {
	o.BroadcastDomainUUIDQueryParameter = broadcastDomainUUID
}

// WithFieldsQueryParameter adds the fields to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithFieldsQueryParameter(fields []string) *IPSubnetCollectionGetParams {
	o.SetFieldsQueryParameter(fields)
	return o
}

// SetFieldsQueryParameter adds the fields to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetFieldsQueryParameter(fields []string) {
	o.FieldsQueryParameter = fields
}

// WithGatewayQueryParameter adds the gateway to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithGatewayQueryParameter(gateway *string) *IPSubnetCollectionGetParams {
	o.SetGatewayQueryParameter(gateway)
	return o
}

// SetGatewayQueryParameter adds the gateway to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetGatewayQueryParameter(gateway *string) {
	o.GatewayQueryParameter = gateway
}

// WithIPRangesEndQueryParameter adds the iPRangesEnd to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithIPRangesEndQueryParameter(iPRangesEnd *string) *IPSubnetCollectionGetParams {
	o.SetIPRangesEndQueryParameter(iPRangesEnd)
	return o
}

// SetIPRangesEndQueryParameter adds the ipRangesEnd to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetIPRangesEndQueryParameter(iPRangesEnd *string) {
	o.IPRangesEndQueryParameter = iPRangesEnd
}

// WithIPRangesFamilyQueryParameter adds the iPRangesFamily to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithIPRangesFamilyQueryParameter(iPRangesFamily *string) *IPSubnetCollectionGetParams {
	o.SetIPRangesFamilyQueryParameter(iPRangesFamily)
	return o
}

// SetIPRangesFamilyQueryParameter adds the ipRangesFamily to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetIPRangesFamilyQueryParameter(iPRangesFamily *string) {
	o.IPRangesFamilyQueryParameter = iPRangesFamily
}

// WithIPRangesStartQueryParameter adds the iPRangesStart to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithIPRangesStartQueryParameter(iPRangesStart *string) *IPSubnetCollectionGetParams {
	o.SetIPRangesStartQueryParameter(iPRangesStart)
	return o
}

// SetIPRangesStartQueryParameter adds the ipRangesStart to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetIPRangesStartQueryParameter(iPRangesStart *string) {
	o.IPRangesStartQueryParameter = iPRangesStart
}

// WithIpspaceNameQueryParameter adds the ipspaceName to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithIpspaceNameQueryParameter(ipspaceName *string) *IPSubnetCollectionGetParams {
	o.SetIpspaceNameQueryParameter(ipspaceName)
	return o
}

// SetIpspaceNameQueryParameter adds the ipspaceName to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetIpspaceNameQueryParameter(ipspaceName *string) {
	o.IpspaceNameQueryParameter = ipspaceName
}

// WithIpspaceUUIDQueryParameter adds the ipspaceUUID to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithIpspaceUUIDQueryParameter(ipspaceUUID *string) *IPSubnetCollectionGetParams {
	o.SetIpspaceUUIDQueryParameter(ipspaceUUID)
	return o
}

// SetIpspaceUUIDQueryParameter adds the ipspaceUuid to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetIpspaceUUIDQueryParameter(ipspaceUUID *string) {
	o.IpspaceUUIDQueryParameter = ipspaceUUID
}

// WithMaxRecordsQueryParameter adds the maxRecords to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithMaxRecordsQueryParameter(maxRecords *int64) *IPSubnetCollectionGetParams {
	o.SetMaxRecordsQueryParameter(maxRecords)
	return o
}

// SetMaxRecordsQueryParameter adds the maxRecords to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetMaxRecordsQueryParameter(maxRecords *int64) {
	o.MaxRecordsQueryParameter = maxRecords
}

// WithNameQueryParameter adds the name to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithNameQueryParameter(name *string) *IPSubnetCollectionGetParams {
	o.SetNameQueryParameter(name)
	return o
}

// SetNameQueryParameter adds the name to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetNameQueryParameter(name *string) {
	o.NameQueryParameter = name
}

// WithOrderByQueryParameter adds the orderBy to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithOrderByQueryParameter(orderBy []string) *IPSubnetCollectionGetParams {
	o.SetOrderByQueryParameter(orderBy)
	return o
}

// SetOrderByQueryParameter adds the orderBy to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetOrderByQueryParameter(orderBy []string) {
	o.OrderByQueryParameter = orderBy
}

// WithReturnRecordsQueryParameter adds the returnRecords to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithReturnRecordsQueryParameter(returnRecords *bool) *IPSubnetCollectionGetParams {
	o.SetReturnRecordsQueryParameter(returnRecords)
	return o
}

// SetReturnRecordsQueryParameter adds the returnRecords to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetReturnRecordsQueryParameter(returnRecords *bool) {
	o.ReturnRecordsQueryParameter = returnRecords
}

// WithReturnTimeoutQueryParameter adds the returnTimeout to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithReturnTimeoutQueryParameter(returnTimeout *int64) *IPSubnetCollectionGetParams {
	o.SetReturnTimeoutQueryParameter(returnTimeout)
	return o
}

// SetReturnTimeoutQueryParameter adds the returnTimeout to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetReturnTimeoutQueryParameter(returnTimeout *int64) {
	o.ReturnTimeoutQueryParameter = returnTimeout
}

// WithSubnetAddressQueryParameter adds the subnetAddress to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithSubnetAddressQueryParameter(subnetAddress *string) *IPSubnetCollectionGetParams {
	o.SetSubnetAddressQueryParameter(subnetAddress)
	return o
}

// SetSubnetAddressQueryParameter adds the subnetAddress to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetSubnetAddressQueryParameter(subnetAddress *string) {
	o.SubnetAddressQueryParameter = subnetAddress
}

// WithSubnetFamilyQueryParameter adds the subnetFamily to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithSubnetFamilyQueryParameter(subnetFamily *string) *IPSubnetCollectionGetParams {
	o.SetSubnetFamilyQueryParameter(subnetFamily)
	return o
}

// SetSubnetFamilyQueryParameter adds the subnetFamily to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetSubnetFamilyQueryParameter(subnetFamily *string) {
	o.SubnetFamilyQueryParameter = subnetFamily
}

// WithSubnetNetmaskQueryParameter adds the subnetNetmask to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithSubnetNetmaskQueryParameter(subnetNetmask *string) *IPSubnetCollectionGetParams {
	o.SetSubnetNetmaskQueryParameter(subnetNetmask)
	return o
}

// SetSubnetNetmaskQueryParameter adds the subnetNetmask to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetSubnetNetmaskQueryParameter(subnetNetmask *string) {
	o.SubnetNetmaskQueryParameter = subnetNetmask
}

// WithTotalCountQueryParameter adds the totalCount to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithTotalCountQueryParameter(totalCount *int64) *IPSubnetCollectionGetParams {
	o.SetTotalCountQueryParameter(totalCount)
	return o
}

// SetTotalCountQueryParameter adds the totalCount to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetTotalCountQueryParameter(totalCount *int64) {
	o.TotalCountQueryParameter = totalCount
}

// WithUsedCountQueryParameter adds the usedCount to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithUsedCountQueryParameter(usedCount *int64) *IPSubnetCollectionGetParams {
	o.SetUsedCountQueryParameter(usedCount)
	return o
}

// SetUsedCountQueryParameter adds the usedCount to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetUsedCountQueryParameter(usedCount *int64) {
	o.UsedCountQueryParameter = usedCount
}

// WithUUIDQueryParameter adds the uuid to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) WithUUIDQueryParameter(uuid *string) *IPSubnetCollectionGetParams {
	o.SetUUIDQueryParameter(uuid)
	return o
}

// SetUUIDQueryParameter adds the uuid to the ip subnet collection get params
func (o *IPSubnetCollectionGetParams) SetUUIDQueryParameter(uuid *string) {
	o.UUIDQueryParameter = uuid
}

// WriteToRequest writes these params to a swagger request
func (o *IPSubnetCollectionGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.AvailableCountQueryParameter != nil {

		// query param available_count
		var qrAvailableCount int64

		if o.AvailableCountQueryParameter != nil {
			qrAvailableCount = *o.AvailableCountQueryParameter
		}
		qAvailableCount := swag.FormatInt64(qrAvailableCount)
		if qAvailableCount != "" {

			if err := r.SetQueryParam("available_count", qAvailableCount); err != nil {
				return err
			}
		}
	}

	if o.AvailableIPRangesEndQueryParameter != nil {

		// query param available_ip_ranges.end
		var qrAvailableIPRangesEnd string

		if o.AvailableIPRangesEndQueryParameter != nil {
			qrAvailableIPRangesEnd = *o.AvailableIPRangesEndQueryParameter
		}
		qAvailableIPRangesEnd := qrAvailableIPRangesEnd
		if qAvailableIPRangesEnd != "" {

			if err := r.SetQueryParam("available_ip_ranges.end", qAvailableIPRangesEnd); err != nil {
				return err
			}
		}
	}

	if o.AvailableIPRangesFamilyQueryParameter != nil {

		// query param available_ip_ranges.family
		var qrAvailableIPRangesFamily string

		if o.AvailableIPRangesFamilyQueryParameter != nil {
			qrAvailableIPRangesFamily = *o.AvailableIPRangesFamilyQueryParameter
		}
		qAvailableIPRangesFamily := qrAvailableIPRangesFamily
		if qAvailableIPRangesFamily != "" {

			if err := r.SetQueryParam("available_ip_ranges.family", qAvailableIPRangesFamily); err != nil {
				return err
			}
		}
	}

	if o.AvailableIPRangesStartQueryParameter != nil {

		// query param available_ip_ranges.start
		var qrAvailableIPRangesStart string

		if o.AvailableIPRangesStartQueryParameter != nil {
			qrAvailableIPRangesStart = *o.AvailableIPRangesStartQueryParameter
		}
		qAvailableIPRangesStart := qrAvailableIPRangesStart
		if qAvailableIPRangesStart != "" {

			if err := r.SetQueryParam("available_ip_ranges.start", qAvailableIPRangesStart); err != nil {
				return err
			}
		}
	}

	if o.BroadcastDomainNameQueryParameter != nil {

		// query param broadcast_domain.name
		var qrBroadcastDomainName string

		if o.BroadcastDomainNameQueryParameter != nil {
			qrBroadcastDomainName = *o.BroadcastDomainNameQueryParameter
		}
		qBroadcastDomainName := qrBroadcastDomainName
		if qBroadcastDomainName != "" {

			if err := r.SetQueryParam("broadcast_domain.name", qBroadcastDomainName); err != nil {
				return err
			}
		}
	}

	if o.BroadcastDomainUUIDQueryParameter != nil {

		// query param broadcast_domain.uuid
		var qrBroadcastDomainUUID string

		if o.BroadcastDomainUUIDQueryParameter != nil {
			qrBroadcastDomainUUID = *o.BroadcastDomainUUIDQueryParameter
		}
		qBroadcastDomainUUID := qrBroadcastDomainUUID
		if qBroadcastDomainUUID != "" {

			if err := r.SetQueryParam("broadcast_domain.uuid", qBroadcastDomainUUID); err != nil {
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

	if o.GatewayQueryParameter != nil {

		// query param gateway
		var qrGateway string

		if o.GatewayQueryParameter != nil {
			qrGateway = *o.GatewayQueryParameter
		}
		qGateway := qrGateway
		if qGateway != "" {

			if err := r.SetQueryParam("gateway", qGateway); err != nil {
				return err
			}
		}
	}

	if o.IPRangesEndQueryParameter != nil {

		// query param ip_ranges.end
		var qrIPRangesEnd string

		if o.IPRangesEndQueryParameter != nil {
			qrIPRangesEnd = *o.IPRangesEndQueryParameter
		}
		qIPRangesEnd := qrIPRangesEnd
		if qIPRangesEnd != "" {

			if err := r.SetQueryParam("ip_ranges.end", qIPRangesEnd); err != nil {
				return err
			}
		}
	}

	if o.IPRangesFamilyQueryParameter != nil {

		// query param ip_ranges.family
		var qrIPRangesFamily string

		if o.IPRangesFamilyQueryParameter != nil {
			qrIPRangesFamily = *o.IPRangesFamilyQueryParameter
		}
		qIPRangesFamily := qrIPRangesFamily
		if qIPRangesFamily != "" {

			if err := r.SetQueryParam("ip_ranges.family", qIPRangesFamily); err != nil {
				return err
			}
		}
	}

	if o.IPRangesStartQueryParameter != nil {

		// query param ip_ranges.start
		var qrIPRangesStart string

		if o.IPRangesStartQueryParameter != nil {
			qrIPRangesStart = *o.IPRangesStartQueryParameter
		}
		qIPRangesStart := qrIPRangesStart
		if qIPRangesStart != "" {

			if err := r.SetQueryParam("ip_ranges.start", qIPRangesStart); err != nil {
				return err
			}
		}
	}

	if o.IpspaceNameQueryParameter != nil {

		// query param ipspace.name
		var qrIpspaceName string

		if o.IpspaceNameQueryParameter != nil {
			qrIpspaceName = *o.IpspaceNameQueryParameter
		}
		qIpspaceName := qrIpspaceName
		if qIpspaceName != "" {

			if err := r.SetQueryParam("ipspace.name", qIpspaceName); err != nil {
				return err
			}
		}
	}

	if o.IpspaceUUIDQueryParameter != nil {

		// query param ipspace.uuid
		var qrIpspaceUUID string

		if o.IpspaceUUIDQueryParameter != nil {
			qrIpspaceUUID = *o.IpspaceUUIDQueryParameter
		}
		qIpspaceUUID := qrIpspaceUUID
		if qIpspaceUUID != "" {

			if err := r.SetQueryParam("ipspace.uuid", qIpspaceUUID); err != nil {
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

	if o.NameQueryParameter != nil {

		// query param name
		var qrName string

		if o.NameQueryParameter != nil {
			qrName = *o.NameQueryParameter
		}
		qName := qrName
		if qName != "" {

			if err := r.SetQueryParam("name", qName); err != nil {
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

	if o.SubnetAddressQueryParameter != nil {

		// query param subnet.address
		var qrSubnetAddress string

		if o.SubnetAddressQueryParameter != nil {
			qrSubnetAddress = *o.SubnetAddressQueryParameter
		}
		qSubnetAddress := qrSubnetAddress
		if qSubnetAddress != "" {

			if err := r.SetQueryParam("subnet.address", qSubnetAddress); err != nil {
				return err
			}
		}
	}

	if o.SubnetFamilyQueryParameter != nil {

		// query param subnet.family
		var qrSubnetFamily string

		if o.SubnetFamilyQueryParameter != nil {
			qrSubnetFamily = *o.SubnetFamilyQueryParameter
		}
		qSubnetFamily := qrSubnetFamily
		if qSubnetFamily != "" {

			if err := r.SetQueryParam("subnet.family", qSubnetFamily); err != nil {
				return err
			}
		}
	}

	if o.SubnetNetmaskQueryParameter != nil {

		// query param subnet.netmask
		var qrSubnetNetmask string

		if o.SubnetNetmaskQueryParameter != nil {
			qrSubnetNetmask = *o.SubnetNetmaskQueryParameter
		}
		qSubnetNetmask := qrSubnetNetmask
		if qSubnetNetmask != "" {

			if err := r.SetQueryParam("subnet.netmask", qSubnetNetmask); err != nil {
				return err
			}
		}
	}

	if o.TotalCountQueryParameter != nil {

		// query param total_count
		var qrTotalCount int64

		if o.TotalCountQueryParameter != nil {
			qrTotalCount = *o.TotalCountQueryParameter
		}
		qTotalCount := swag.FormatInt64(qrTotalCount)
		if qTotalCount != "" {

			if err := r.SetQueryParam("total_count", qTotalCount); err != nil {
				return err
			}
		}
	}

	if o.UsedCountQueryParameter != nil {

		// query param used_count
		var qrUsedCount int64

		if o.UsedCountQueryParameter != nil {
			qrUsedCount = *o.UsedCountQueryParameter
		}
		qUsedCount := swag.FormatInt64(qrUsedCount)
		if qUsedCount != "" {

			if err := r.SetQueryParam("used_count", qUsedCount); err != nil {
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

// bindParamIPSubnetCollectionGet binds the parameter fields
func (o *IPSubnetCollectionGetParams) bindParamFields(formats strfmt.Registry) []string {
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

// bindParamIPSubnetCollectionGet binds the parameter order_by
func (o *IPSubnetCollectionGetParams) bindParamOrderBy(formats strfmt.Registry) []string {
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
