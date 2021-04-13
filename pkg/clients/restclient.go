package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type RestClient interface {
	Get(url string) ([]byte, error)
}

type RestHttpClient struct {
	UserAgent string
}

func NewRestHttpClient() *http.Client {
	client := http.DefaultClient
	return client
}

func (c *RestHttpClient) newRequest(method, path string, body interface{}) (*http.Request, error) {

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, path, buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *RestHttpClient) do(req *http.Request) (*http.Response, error) {
	client := NewRestHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *RestHttpClient) Get(url string) ([]byte, error) {
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Unable to get response body " + err.Error())
			}
			log.Debugf("response for call " + url + "is " + string(bodyBytes))
			return bodyBytes, nil
		} else {
			log.Error("Status code is " + resp.Status)
			respBody, _ := ioutil.ReadAll(resp.Body)
			log.Error("response body is " + string(respBody))
			return nil, errors.New(string(respBody))
		}
	} else {
		return nil, errors.New("Empty response for " + url)
	}
}
