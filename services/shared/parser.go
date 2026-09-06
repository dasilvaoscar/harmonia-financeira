package shared

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseBody(r *http.Request) (map[string]interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	var bodyMap map[string]interface{}

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	return bodyMap, nil
}
