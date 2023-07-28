package src

import (
	"fmt"
	"io"
	"net/http"
)

const BaseUrl = "https://api.hypixel.net"

type HypixelApi struct {
	Request func(request Request) (response *Response, err error)
}

func NewHypixelApi(key string) HypixelApi {
	return HypixelApi{
		Request: func(request Request) (response *Response, err error) {
			url := BaseUrl + request.Path
			if request.Parameter != nil {
				Param := *request.Parameter
				url += "?" + Param.Key + "=" + Param.Value
			}

			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Printf("hypixel: error making http request: %s\n", err)
				return nil, err
			}

			req.Header.Set("API-Key", key)
			res, err := client.Do(req)
			if err != nil {
				fmt.Printf("hypixel: error sending http request: %s\n", err)
				return nil, err
			}

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("hypixel: could not read response body: %s\n", err)
				return nil, err
			}

			return &Response{Code: res.StatusCode, Data: string(resBody)}, nil
		},
	}
}
