package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type AvailableEndpoint struct {
	Parameter *string `json:"parameter"`
	MaxLive   int64   `json:"maxLive"`
}

type EndPointManager struct {
	Endpoints map[string]AvailableEndpoint
	Parse     func(request *http.Request) (*PendingRequest, error)
}

var ErrParametersRequired = errors.New("parameter(s) required for this endpoint")
var ErrParametersTooLong = errors.New("parameter(s) value is too long")
var ErrEndpointNotImplemented = errors.New("endpoint not implemented")

func NewEndPointManager(file string) (manage *EndPointManager, err error) {
	endpoints, err := load(file)
	if err != nil {
		return nil, err
	}

	return &EndPointManager{
		Endpoints: *endpoints,
		Parse: func(request *http.Request) (*PendingRequest, error) {
			url := request.URL
			path := url.Path

			endpoint, ok := (*endpoints)[path]
			if !ok {
				return nil, ErrEndpointNotImplemented
			}

			if endpoint.Parameter == nil {
				return &PendingRequest{Path: path, MaxLive: endpoint.MaxLive}, nil
			}

			value := strings.ToLower(url.Query().Get(*endpoint.Parameter))
			if len(value) == 0 {
				return nil, ErrParametersRequired
			} else if len(value) > 36 {
				return nil, ErrParametersTooLong
			}

			return &PendingRequest{
				Path:    path,
				MaxLive: endpoint.MaxLive,
				Key:     endpoint.Parameter,
				Value:   &value,
			}, nil
		},
	}, nil
}

func load(file string) (endpoints *map[string]AvailableEndpoint, err error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	byteValue, _ := io.ReadAll(jsonFile)
	var manager map[string]AvailableEndpoint

	err = json.Unmarshal(byteValue, &manager)
	if err != nil {
		return nil, err
	}

	err = jsonFile.Close()
	if err != nil {
		return nil, err
	}

	return &manager, nil
}
