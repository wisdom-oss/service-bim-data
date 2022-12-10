package handlers

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	e "microservice/errors"
	"microservice/helpers"
	"microservice/vars"
)

func AuthorizationCheck(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			logger := log.WithFields(
				log.Fields{
					"middleware": true,
					"title":      "AuthorizationCheck",
				},
			)
			logger.Debug("Checking the incoming request for authorization information set by the gateway")
			if request.URL.Path == "/ping" {
				nextHandler.ServeHTTP(responseWriter, request)
				return
			}
			// Get the scopes the requesting user has
			scopes := request.Header.Get("X-Authenticated-Scope")
			// Check if the string is empty
			if strings.TrimSpace(scopes) == "" {
				logger.Warning("Unauthorized request detected. The required header had no content or was not set")
				helpers.SendRequestError(e.UnauthorizedRequest, responseWriter)
				return
			}

			scopeList := strings.Split(scopes, ",")
			if !helpers.StringArrayContains(scopeList, vars.ScopeConfiguration.ScopeValue) {
				logger.Error("Request rejected. The user is missing the scope needed for accessing this service")
				helpers.SendRequestError(e.MissingScope, responseWriter)
				return
			}
			// Call the next handler which will continue handling the request
			nextHandler.ServeHTTP(responseWriter, request)
		},
	)
}

/*
PingHandler

This handler is used to test if the service is able to ping itself. This is done to run a healthcheck on the container
*/
func PingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

/*
RequestHandler

This handler shows how a basic handler works and how to send back a message
*/
func RequestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	logger := log.WithFields(
		log.Fields{
			"middleware": true,
			"title":      "RequestHandler",
		},
	)
	// Check if the request contains a modelID and instanceID
	modelIDSet := request.URL.Query().Has("modelID")
	instanceIDSet := request.URL.Query().Has("instanceID")

	if !modelIDSet || !instanceIDSet {
		logger.Warning("incoming request did not contain the needed query parameters")
		helpers.SendRequestError(e.MissingQueryParameter, responseWriter)
		return
	}
}
