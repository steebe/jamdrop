package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jamdrop/logger"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	APIBaseURL      = "https://api.spotify.com"
	AccountsBaseURL = "https://accounts.spotify.com"
)

func apiPath(path string, v ...interface{}) string {
	return fmt.Sprintf("%s/%s", APIBaseURL, strings.TrimPrefix(fmt.Sprintf(path, v...), "/"))
}

func accountsPath(path string, v ...interface{}) string {
	return fmt.Sprintf("%s/%s", AccountsBaseURL, strings.TrimPrefix(fmt.Sprintf(path, v...), "/"))
}

func (s *Client) request(r *http.Request) (*http.Response, []byte, error) {
	s.Logger.Debug(fmt.Sprintf("%s %s", r.Method, r.URL.String()))

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to make request")
	}

	var body []byte
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, nil, errors.Wrap(err, "failed to read request body")
	}

	if res.StatusCode < 200 || 299 < res.StatusCode {
		s.Logger.Error("bad status", logger.Field("status", res.StatusCode))
		var errRes errorResponse
		if err := json.Unmarshal(body, &errRes); err != nil {
			return nil, nil, errors.Wrap(err, "failed to unmarshal request response")
		}

		return res, body, errRes.Error
	}

	return res, body, nil
}

func (s *Client) requestToJSON(r *http.Request, v interface{}) error {
	_, body, err := s.request(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return errors.Wrapf(err, "failed to unmarshal request response: %s", body)
	}

	return nil
}
