package src

type Request struct {
	Path      string
	Parameter *Parameter

	MaxLive int64
}

type Parameter struct {
	Key   string
	Value string
}

type Response struct {
	Code int
	Data string
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
