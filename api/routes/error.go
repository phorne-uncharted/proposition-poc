package routes

import (
	"net/http"

	log "github.com/unchartedsoftware/plog"
)

var (
	verboseError = false
)

// SetVerboseError sets the flag determining if the client should receive
// error details
func SetVerboseError(verbose bool) {
	verboseError = verbose
}

func handleError(w http.ResponseWriter, err error) {
	handleErrorType(w, err, http.StatusInternalServerError)
}

func handleErrorType(w http.ResponseWriter, err error, code int) {
	log.Errorf("%+v", err)
	errMessage := "An error occured on the server while processing the request"
	if verboseError {
		errMessage = err.Error()
	}
	http.Error(w, errMessage, code)
}
