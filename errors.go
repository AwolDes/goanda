package goanda

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func newAPIError(request *http.Request, response *http.Response) APIError {
	defer response.Body.Close()

	msg := struct {
		ErrorMessage string
		RejectReason string
	}{}

	apiErr := APIError{
		Response: response,
		Request:  request,
	}

	b, _ := ioutil.ReadAll(response.Body)
	err := json.Unmarshal(b, &msg)
	if err != nil {
		apiErr.Message = string(b)
	} else {
		apiErr.Message = msg.RejectReason + msg.ErrorMessage
	}
	return apiErr
}

// APIError is returned when the Oanda server responds with an error
//
// Message is the returned error message from the server if possible to unmarshal,
// otherwise it is simply the entire body of the response
type APIError struct {
	Request  *http.Request
	Response *http.Response
	Message  string
}

// APIError implements error
func (a APIError) Error() string {
	return fmt.Sprintf("Oanda API Error [Url: %v, Response: %v]: %v",
		a.Request.URL.String(),
		a.Response.Status,
		a.Message,
	)
}
