package restez

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HandleGET will parse the query string, and write a JSON response based on ResponseType
func HandleGET[ResponseType any](fn func(map[string]string) (ResponseType, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w)(fn(queryParameters(r)))
	}
}

// queryParameters converts the default map[string][]string to a simpler to use map[string]string.
func queryParameters(r *http.Request) map[string]string {
	queryParameters := map[string]string{}
	for k, v := range r.URL.Query() {
		queryParameters[k] = ""
		if len(v) > 0 {
			queryParameters[k] = v[0]
		}
		if len(v) > 1 {
			log.Printf("request on %s received > 1 values for query parameter %s. extra values discarded.", r.URL.Path, k)
		}
	}
	return queryParameters
}

// HandlePOST will parse the request body with RequestType, and write a JSON response based on ResponseType.
func HandlePOST[RequestType, ResponseType any](fn func(RequestType) (ResponseType, error)) http.HandlerFunc {
	return handleJSONBody(fn)
}

// HandlePUT will parse the request body with RequestType, and write a JSON response based on ResponseType.
func HandlePUT[RequestType, ResponseType any](fn func(RequestType) (ResponseType, error)) http.HandlerFunc {
	return handleJSONBody(fn)
}

// handleJSONBody does the work of parsing JSON, writing JSON.
func handleJSONBody[RequestType, ResponseType any](f func(RequestType) (ResponseType, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RequestType
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			WriteError(w, fmt.Errorf("unable to parse request as type %T: %v", request, err))
			return
		}

		WriteResponse(w)(f(request))
	}
}

func WriteResponse(w http.ResponseWriter) func(any, error) {
	return func(response any, err error) {
		if err != nil {
			WriteError(w, err)
			return
		}

		WriteSuccess(w, response)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	WriteJSON(w, map[string]any{
		"status": "ERROR",
		"error":  err.Error(),
	})
}

func WriteSuccess(w http.ResponseWriter, response any) {
	WriteJSON(w, map[string]any{
		"status":   "OK",
		"response": response,
	})
}

func WriteJSON(w http.ResponseWriter, content any) {

	w.Header().Set("Content-Type", "application/json")

	marshalled, err := json.Marshal(content)
	if err != nil {
		log.Fatalf("failed to marshal content (type %T) to JSON: %v", content, err)
	}

	_, err = w.Write(marshalled)
	if err != nil {
		log.Printf("failed to write content (%s) to client: %v", marshalled, err)
	}
}
