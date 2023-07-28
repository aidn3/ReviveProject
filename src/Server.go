package src

import (
	"errors"
	"fmt"
	"istio.io/pkg/cache"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Listen(manager EndpointManager, hypixel Hypixel, cache cache.ExpiringCache, port int) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		serverHandler(manager, hypixel, cache, writer, request)
	})

	portString := ":" + strconv.Itoa(port)
	fmt.Println("Starting server at port " + portString)
	if err := http.ListenAndServe(portString, nil); err != nil {
		log.Fatal(err)
	}
}

func serverHandler(manager EndpointManager, hypixel Hypixel, expiringCache cache.ExpiringCache, w http.ResponseWriter, request *http.Request) {
	fmt.Println("request: " + request.URL.String())

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-name", "ReviveProject")
	w.Header().Set("X-version", "0.0.1")

	if request.Method != http.MethodGet {
		returnError(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	pendingRequest, err := manager.Parse(request)
	if err != nil {
		switch {
		case errors.Is(err, ErrEndpointNotImplemented):
			returnError(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, ErrParametersRequired):
			returnError(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, ErrParametersTooLong):
			returnError(w, err.Error(), http.StatusRequestURITooLong)
			return
		default:
			_, _ = fmt.Fprintln(os.Stderr, err)
			returnError(w, "unknown error. Admin has been notified.", http.StatusInternalServerError)
			return
		}
	}

	cacheKey := PendingRequestToString(*pendingRequest)
	response, ok := expiringCache.Get(cacheKey)
	if ok {
		fmt.Println("cache hit: " + request.URL.String())
		castedResponse := response.(*Response)
		w.WriteHeader(castedResponse.Code)
		_, _ = fmt.Fprintf(w, castedResponse.Data)
		return
	}

	hypixelResponse, err := hypixel.Request(*pendingRequest)
	if err != nil {
		returnError(w, "error while trying to connect to hypixel", http.StatusInternalServerError)
		return
	}

	// Allow developers to change the final output
	Transform(*pendingRequest, hypixelResponse)

	w.WriteHeader(hypixelResponse.Code)
	_, _ = fmt.Fprintf(w, hypixelResponse.Data)

	if hypixelResponse.Cache {
		expiringCache.Set(cacheKey, hypixelResponse)
	}
}

func returnError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	_, _ = fmt.Fprintln(w, "{\"success\": false, \"cause\": \""+message+"\"}")
}
