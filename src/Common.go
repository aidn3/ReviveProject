package src

type PendingRequest struct {
	Path  string
	Key   *string
	Value *string

	MaxLive int64
}

type Response struct {
	Code int
	Data string
}

func PendingRequestToString(request PendingRequest) string {
	// It will be converted to valid url query.
	// Valid string as well as easier to read and process
	result := request.Path
	if request.Key != nil {
		result += "?" + *request.Key + "=" + *request.Value
	}
	return result
}
