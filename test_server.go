//provides a basic http server for testing, acts as an stub
package vtputils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type TestReply struct {
	Code  int    //HTTP reply code
	Value string //reply value
}
type TestServer struct {
	Reply        TestReply
	Port         string
	EndPoint     string //endpoint to test = "/test"
	stop         chan int
	CustomHander Handler
	mux          *http.ServeMux
}

type Handler func(http.ResponseWriter, *http.Request)

func serve404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	log.Printf("Error- %s %s %s", r.Method, r.URL, r.Proto)
	fmt.Fprint(w, "Test Error")
}

func (s *TestServer) HandleEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		log.Printf("Header[%q] = %q\n", k, v)
	}
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		log.Printf("Form[%q]=%q\n", k, v)
	}
	defer r.Body.Close()
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		log.Println("Body value:")
		log.Println(string(b))
	}

	log.Printf("responding with value %d\n", s.Reply.Code)
	log.Println(s.Reply.Value)
	w.WriteHeader(s.Reply.Code)
	fmt.Fprint(w, s.Reply.Value)
}
func (s *TestServer) Done() {
	s.stop <- 0
}

//add routes to test server
func (s *TestServer) AddRoute(endpoint string, handler Handler) {
	log.Println("adding route " + endpoint)
	if s.mux == nil {
		s.Init()
	}
	s.mux.HandleFunc(endpoint, handler)
}

func (s *TestServer) Init() {
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", serve404)
}

func (s *TestServer) StartTestServer() {
	if s.EndPoint != "" {
		if s.CustomHander == nil {
			s.AddRoute(s.EndPoint, s.HandleEndpoint)
		} else {
			s.AddRoute(s.EndPoint, s.CustomHander)
		}
	}
	s.stop = make(chan int)

	ln, err := net.Listen("tcp", "localhost:"+s.Port)
	if err != nil {
		log.Fatalf("Can't listen: %s", err)
	} else {
		log.Println("started to test server in " + s.Port)
	}
	go http.Serve(ln, s.mux)
	<-s.stop
	ln.Close()
	log.Println("bye bye..")
}
