package vtputils

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	p := RequestParms{Endpoint: "http://localhost:8000/test", Method: HTTP_GET}
	s := TestServer{Reply: TestReply{200, "Test Ok"}, Port: "8000", EndPoint: "/test"}
	go s.StartTestServer()
	resp, err := Curl(p)
	if err != nil {
		t.Error(err)
	}
	if _, err := HttpStatus(resp); err != nil {
		t.Error(err)
	}
	s.Done()
}

func TestGet404(t *testing.T) {
	p := RequestParms{Endpoint: "http://localhost:8000/test2", Method: HTTP_GET}
	s := TestServer{Reply: TestReply{200, "Test Ok"}, Port: "8000", EndPoint: "/test"}
	go s.StartTestServer()
	resp, err := Curl(p)
	if err != nil {
		t.Error(err)
	}
	if value, _ := HttpStatus(resp); value != 404 {
		t.Error(err)
	}
	s.Done()
}

func TestGetCustomHeader(t *testing.T) {
	p := RequestParms{Endpoint: "http://localhost:8000/test", Method: HTTP_GET}
	s := TestServer{Reply: TestReply{200, "Test Ok"}, Port: "8000", EndPoint: "/test"}
	p.Headers = make(map[string]string)
	p.Headers["X-CUSTOM-HEADER"] = "test"
	s.CustomHander = func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			t.Logf("Header[%q] = %q\n", k, v)
			if k == "X-Custom-Header" {
				w.WriteHeader(200)
				fmt.Fprint(w, "header found")
				return
			}
		}
		w.WriteHeader(404)
		fmt.Fprint(w, "header not found")
	}
	go s.StartTestServer()
	resp, err := Curl(p)
	if err != nil {
		t.Error(err)
	}
	if _, err := HttpStatus(resp); err != nil {
		t.Error(err)
	}
	s.Done()
}
func TestGetAuth(t *testing.T) {
	p := RequestParms{Endpoint: "http://localhost:8000/test", Method: HTTP_GET}
	s := TestServer{Reply: TestReply{200, "Test Ok"}, Port: "8000", EndPoint: "/test"}
	p.Username = "test"
	p.Password = "1234"
	s.CustomHander = func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			t.Logf("Header[%q] = %q\n", k, v)
			if k == "Authorization" {
				w.WriteHeader(200)
				fmt.Fprint(w, "header found")
				return
			}
		}
		w.WriteHeader(404)
		fmt.Fprint(w, "header not found")
	}
	go s.StartTestServer()
	resp, err := Curl(p)
	if err != nil {
		t.Error(err)
	}
	if _, err := HttpStatus(resp); err != nil {
		t.Error(err)
	}
	s.Done()
}
func TestPost(t *testing.T) {
	p := RequestParms{Endpoint: "http://localhost:8000/test", Method: HTTP_POST}
	parm := make(map[string]string)
	parm["q"] = "test"
	p.Params = parm
	s := TestServer{Reply: TestReply{200, "Test Ok"}, Port: "8000", EndPoint: "/test"}
	go s.StartTestServer()
	resp, err := Curl(p)
	if err != nil {
		t.Error(err)
	}
	if _, err := HttpStatus(resp); err != nil {
		t.Error(err)
	}
	s.Done()
}
func TestDelete(t *testing.T) {
	p := RequestParms{Endpoint: "http://localhost:8000/test", Method: HTTP_DELETE}
	parm := make(map[string]string)
	parm["q"] = "test"
	p.Params = parm
	s := TestServer{Reply: TestReply{200, "Test Ok"}, Port: "8000", EndPoint: "/test"}
	go s.StartTestServer()
	resp, err := Curl(p)
	if err != nil {
		t.Error(err)
	}
	if _, err := HttpStatus(resp); err != nil {
		t.Error(err)
	}
	s.Done()
}
func TestJSONPost(t *testing.T) {
	p := RequestParms{Endpoint: "http://localhost:8000/test", Method: HTTP_JSONPOST}
	parm := make(map[string]string)
	parm["json"] = "{\"Test\":\"test\"}"
	p.Params = parm
	s := TestServer{Reply: TestReply{200, "Test Ok"}, Port: "8000", EndPoint: "/test"}
	go s.StartTestServer()
	resp, err := Curl(p)
	if err != nil {
		t.Error(err)
	}
	if _, err := HttpStatus(resp); err != nil {
		t.Error(err)
	}
	s.Done()
}
func TestCustomHanler(t *testing.T) {
	p := RequestParms{Endpoint: "http://localhost:8000/test", Method: HTTP_GET}
	s := TestServer{Port: "8000", EndPoint: "/test"}
	s.CustomHander = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "testing testing...")
	}
	go s.StartTestServer()
	resp, err := Curl(p)
	if err != nil {
		t.Error(err)
	}
	if _, err := HttpStatus(resp); err != nil {
		t.Error(err)
	}
	b, _ := Body(resp)
	t.Log(b)
	s.Done()
}
