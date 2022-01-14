package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/phorne-uncharted/proposition-poc/api/util"
	"github.com/pkg/errors"
)

func getPostParameters(r *http.Request) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse POST request")
	}

	// could have an empty parameter
	if len(body) == 0 {
		return map[string]interface{}{}, nil
	}

	return util.Unmarshal(body)
}

func handleJSON(w http.ResponseWriter, data interface{}) error {
	// marshal data
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// send response
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
