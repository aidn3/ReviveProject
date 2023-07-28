package src

type Request struct {
	Endpoint Endpoint

	Path      string
	Parameter *Parameter
}

type Endpoint struct {
	Parameter *string `json:"parameter"`
	MaxLive   int64   `json:"maxLive"`
	Custom    bool    `json:"custom"`
}

type Parameter struct {
	Key   string
	Value string
}

type Response struct {
	Code  int
	Data  string
	Cache bool
}

func PendingRequestToString(request Request) string {
	// It will be converted to valid url query.
	// Valid string as well as easier to read and process
	result := request.Path
	if request.Parameter != nil {
		Param := *request.Parameter
		result += "?" + Param.Key + "=" + Param.Value
	}
	return result
}
