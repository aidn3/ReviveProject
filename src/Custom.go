package src

import (
	"encoding/json"
	"net/http"
)

// Serve This function is only called on non-cached valid requests
// with a positive code 200 response
func Serve(request Request, response *Response, hypixel Hypixel) {
	// example code

	// custom endpoint
	// declared in "endpoints.json"
	if request.Path == "/example" {
		response.Code = 200
		response.Cache = true
		response.Data = "{\"success\": true, \"example\": \"This is a custom example endpoint\"}"
	}

	// custom response for an existing endpoint
	if request.Path == "/skyblock/profiles" {
		var decoded map[string]any

		err := json.Unmarshal([]byte(response.Data), &decoded)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Data = "{\"success\": false, \"cause\": \"could not parse hypixel json?\"}"
			return
		}

		decoded["info"] = "This is an extra field in the response"
		encoded, _ := json.Marshal(decoded)
		response.Data = string(encoded)
	}
}
