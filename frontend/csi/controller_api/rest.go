// Copyright 2022 NetApp, Inc. All Rights Reserved.

package controllerAPI

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/netapp/trident/config"
	. "github.com/netapp/trident/logger"
	"github.com/netapp/trident/utils"
)

const HTTPClientTimeout = time.Second * 30

type ControllerRestClient struct {
	url        string
	httpClient http.Client
}

func CreateTLSRestClient(url, caFile, certFile, keyFile string) (TridentController, error) {
	tlsConfig := &tls.Config{MinVersion: config.MinClientTLSVersion}
	if "" != caFile {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
		tlsConfig.ServerName = config.ServerCertName
	} else {
		tlsConfig.InsecureSkipVerify = true
	}
	if "" != certFile && "" != keyFile {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return &ControllerRestClient{
		url: url,
		httpClient: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
			Timeout: HTTPClientTimeout,
		},
	}, nil
}

// InvokeAPI makes a REST call to the CSI Controller REST endpoint. The body must be a marshaled JSON byte array (
// or nil). The method is the HTTP verb (i.e. GET, POST, ...).  The resource path is appended to the base URL to
// identify the desired server resource; it should start with '/'.
func (c *ControllerRestClient) InvokeAPI(
	ctx context.Context, requestBody []byte, method, resourcePath string, redactRequestBody,
	redactResponseBody bool,
) (*http.Response, []byte, error) {
	// Build URL
	url := c.url + resourcePath

	var request *http.Request
	var err error
	var prettyRequestBuffer bytes.Buffer
	var prettyResponseBuffer bytes.Buffer

	// Create the request
	if requestBody == nil {
		request, err = http.NewRequestWithContext(ctx, method, url, nil)
	} else {
		request, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBody))
	}
	if err != nil {
		return nil, nil, err
	}

	request.Header.Set("X-Request-ID", fmt.Sprint(ctx.Value(ContextKeyRequestID)))
	request.Header.Set("Content-Type", "application/json")

	// Log the request
	if requestBody != nil {
		if err = json.Indent(&prettyRequestBuffer, requestBody, "", "  "); err != nil {
			return nil, nil, fmt.Errorf("error formating request body; %v", err)
		}
	}

	utils.LogHTTPRequest(request, prettyRequestBuffer.Bytes(), redactRequestBody)

	response, err := c.httpClient.Do(request)
	if err != nil {
		err = fmt.Errorf("error communicating with Trident CSI Controller; %v", err)
		return nil, nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body; %v", err)
	}

	if responseBody != nil {
		if err = json.Indent(&prettyResponseBuffer, responseBody, "", "  "); err != nil {
			return nil, nil, fmt.Errorf("error formating response body; %v", err)
		}
	}
	utils.LogHTTPResponse(ctx, response, prettyResponseBuffer.Bytes(), redactResponseBody)

	return response, responseBody, err
}

type CreateNodeResponse struct {
	TopologyLabels map[string]string `json:"topologyLabels"`
}

// CreateNode registers the node with the CSI controller server
func (c *ControllerRestClient) CreateNode(ctx context.Context, node *utils.Node) (CreateNodeResponse, error) {
	nodeData, err := json.MarshalIndent(node, "", " ")
	if err != nil {
		return CreateNodeResponse{}, fmt.Errorf("error parsing create node request; %v", err)
	}
	resp, respBody, err := c.InvokeAPI(ctx, nodeData, "PUT", config.NodeURL+"/"+node.Name, false, false)
	if err != nil {
		return CreateNodeResponse{}, fmt.Errorf("could not log into the Trident CSI Controller: %v", err)
	}
	createResponse := CreateNodeResponse{}
	if err := json.Unmarshal(respBody, &createResponse); err != nil {
		return createResponse, fmt.Errorf("could not parse node : %s; %v", string(respBody), err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return createResponse, fmt.Errorf("could not add CSI node")
	}
	return createResponse, nil
}

type ListNodesResponse struct {
	Nodes []string `json:"nodes"`
	Error string   `json:"error,omitempty"`
}

// GetNodes returns a list of nodes registered with the controller
func (c *ControllerRestClient) GetNodes(ctx context.Context) ([]string, error) {
	resp, respBody, err := c.InvokeAPI(ctx, nil, "GET", config.NodeURL, false, false)
	if err != nil {
		return nil, fmt.Errorf("could not log into the Trident CSI Controller: %v", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not list the CSI nodes")
	}

	// Parse JSON data
	respData := ListNodesResponse{}
	if err := json.Unmarshal(respBody, &respData); err != nil {
		return nil, fmt.Errorf("could not parse node list: %s; %v", string(respBody), err)
	}

	return respData.Nodes, nil
}

// DeleteNode deregisters the node with the CSI controller server
func (c *ControllerRestClient) DeleteNode(ctx context.Context, name string) error {
	resp, _, err := c.InvokeAPI(ctx, nil, "DELETE", config.NodeURL+"/"+name, false, false)
	if err != nil {
		return fmt.Errorf("could not log into the Trident CSI Controller: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNoContent:
	case http.StatusUnprocessableEntity:
	case http.StatusNotFound:
	case http.StatusGone:
		break
	default:
		return fmt.Errorf("could not delete the node")
	}
	return nil
}

type GetCHAPResponse struct {
	CHAP  *utils.IscsiChapInfo `json:"chap"`
	Error string               `json:"error,omitempty"`
}

// GetChap requests the current CHAP credentials for a given volume/node pair from the Trident controller
func (c *ControllerRestClient) GetChap(ctx context.Context, volume, node string) (*utils.IscsiChapInfo, error) {
	resp, respBody, err := c.InvokeAPI(ctx, nil, "GET", config.ChapURL+"/"+volume+"/"+node, false, true)
	if err != nil {
		return &utils.IscsiChapInfo{}, fmt.Errorf("could not communicate with the Trident CSI Controller: %v", err)
	}
	createResponse := GetCHAPResponse{}
	if err := json.Unmarshal(respBody, &createResponse); err != nil {
		return nil, fmt.Errorf("could not parse CHAP info : %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		msg := "could not add get CHAP info"
		Logc(ctx).WithError(fmt.Errorf(createResponse.Error)).Errorf(msg)
		return nil, fmt.Errorf(msg)
	}
	return createResponse.CHAP, nil
}

// UpdateVolumePublication updates a volume publication in the Trident controller.  At present, the only value
// modified is NotSafeToAttach.
func (c *ControllerRestClient) UpdateVolumePublication(
	ctx context.Context, publication *utils.VolumePublicationExternal,
) error {
	body, err := json.Marshal(publication)
	if err != nil {
		return fmt.Errorf("error marshaling update publication request; %v", err)
	}
	url := config.PublicationURL + "/" + publication.VolumeName + "/" + publication.NodeName
	resp, _, err := c.InvokeAPI(ctx, body, "PUT", url, false, false)
	if err != nil {
		return fmt.Errorf("could not log into the Trident CSI Controller: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not update volume publication")
	}
	return nil
}

func (c *ControllerRestClient) UpdateVolumeLUKSPassphraseNames(
	ctx context.Context, volumeName string, passphraseNames []string,
) error {
	operations := passphraseNames
	body, err := json.Marshal(operations)
	if err != nil {
		return fmt.Errorf("could not marshal JSON; %v", err)
	}
	url := config.VolumeURL + "/" + volumeName + "/luksPassphraseNames"
	resp, _, err := c.InvokeAPI(ctx, body, "PUT", url, false, false)
	if err != nil {
		return fmt.Errorf("could not log into the Trident CSI Controller: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not update volume LUKS passphrase names")
	}
	return nil
}
