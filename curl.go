//utils package from victor tuson palau
//utilities for interacting with REST apis

package vtputils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//Options for Methods on the query
const (
	HTTP_PUT        = "PUT"
	HTTP_POST       = "POST"
	HTTP_JSONPOST   = "JSONPOST"
	HTTP_DELETE     = "DELETE"
	HTTP_GET        = "GET"
	DEFAULT_TIMEOUT = 10
)

func formValues(params map[string]string) url.Values {
	query := make(url.Values)
	for key, value := range params {
		query.Set(key, value)
	}
	return query
}

func buildGetUrl(params map[string]string, endpoint string) string {
	query := formValues(params)
	result := endpoint
	if (len(params)) > 0 {
		result = endpoint + "?" + query.Encode()
	}
	return result
}

//Request parameters for a api query
//set requestParms.Params key as 'json' to identify payload for a JSON POST query
type RequestParms struct {
	Params   map[string]string //query parameters such as q=hello
	Endpoint string            //full endpoint to be queried e.g. http://localhost:8000/test
	Method   string            // chose from the HTTP_* constants in this package
	Apikey   string            //Apikey is used for basic HTTP authentication header
	Username string            // username for basic authentication
	Password string            //username for basic autthentication
	Headers  map[string]string //headers pair with values
	Timeout  time.Duration     //Timeout duration, if not set uses the default valued of 10sec
}

//Checks the status code and returns an error if appropiate
func HttpStatus(resp *http.Response) (int, error) {
	var err error
	if resp.StatusCode > 399 {
		err = errors.New(strconv.Itoa(resp.StatusCode) + " " + resp.Status)
	}
	return resp.StatusCode, err
}

//unmarshalls a json response into a interface 'result'
func HttpUnmarshall(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(result)
}

//returns the body of http call as string
func Body(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	b, e := ioutil.ReadAll(resp.Body)
	return string(b), e
}

//converts an interface to json []byte - use string() cast to pass a payload to a query
func ToJson(v interface{}) ([]byte, error) {
	jsonEncoded, err := json.MarshalIndent(v, "", " ")
	return jsonEncoded, err
}

func (p RequestParms) populateHeaders(req *http.Request) {
	if p.Apikey != "" {
		req.Header.Add("Authorization", p.Apikey)
	}
	if p.Username != "" {
		req.Header.Add("Authorization", p.Username+":"+p.Password)
	}
	for k, v := range p.Headers {
		req.Header.Add(k, v)
	}
}

//provides a httpResponse for a GET,DELETE,POST or PUT. Can support json data, use "json" as a key on the parmeter
func Curl(p RequestParms) (*http.Response, error) {
	var client http.Client
	var req *http.Request
	var resp *http.Response
	var err error

	client.Timeout = p.Timeout
	if client.Timeout == 0 {
		client.Timeout = time.Second * DEFAULT_TIMEOUT
	}

	if p.Method == HTTP_GET || p.Method == HTTP_DELETE {
		url := buildGetUrl(p.Params, p.Endpoint)
		req, _ = http.NewRequest(p.Method, url, nil)
		p.populateHeaders(req)
		resp, err = client.Do(req)
	}
	if p.Method == HTTP_POST || p.Method == HTTP_PUT {
		url := p.Endpoint
		data := formValues(p.Params)
		req, err = http.NewRequest(p.Method, url, bytes.NewBufferString(data.Encode()))
		if err != nil {
			log.Println(err)
		}

		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		if p.Headers == nil {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		} else {
			p.populateHeaders(req)
		}

		resp, err = client.Do(req)
	}
	if p.Method == HTTP_JSONPOST {
		url := p.Endpoint
		data := p.Params["json"]
		req, err = http.NewRequest("POST", url, bytes.NewBufferString(data))
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(data)))
		p.populateHeaders(req)
		resp, err = client.Do(req)
	}
	return resp, err
}
