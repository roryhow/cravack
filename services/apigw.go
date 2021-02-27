package services

import (
	"bytes"
	"encoding/base64"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func APIGatewayProxyRequestToHTTPRequest(req events.APIGatewayProxyRequest) (*http.Request, error) {
	decodedBody := []byte(req.Body)
	if req.IsBase64Encoded {
		base64Body, err := base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return nil, err
		}
		decodedBody = base64Body
	}

	httpReq, err := http.NewRequest(
		"POST", // req.HTTPMethod is missing a value, we know this is a POST
		req.Path,
		bytes.NewReader(decodedBody),
	)
	if err != nil {
		return nil, err
	}

	if req.MultiValueHeaders != nil {
		for k, values := range req.MultiValueHeaders {
			for _, value := range values {
				httpReq.Header.Add(k, value)
			}
		}
	} else {
		for header := range req.Headers {
			httpReq.Header.Add(header, req.Headers[header])
		}
	}

	httpReq.RequestURI = httpReq.URL.RequestURI()

	return httpReq, nil
}
