//provides a basic http server for testing, acts as an stub
package vtputils

import (
	"fmt"
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
	CustomHander func(http.ResponseWriter, *http.Request)
}

func serve404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "Test Error")
}

func (s *TestServer) handleEndpoint(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("responding with value %d\n", s.Reply.Code)
	log.Println(s.Reply.Value)
	w.WriteHeader(s.Reply.Code)
	fmt.Fprint(w, s.Reply.Value)
}
func (s *TestServer) Done() {
	s.stop <- 0
}
func (s *TestServer) StartTestServer() {
	myMux := http.NewServeMux()
	myMux.HandleFunc("/", serve404)
	if s.CustomHander == nil {
		myMux.HandleFunc(s.EndPoint, s.handleEndpoint)
	} else {
		myMux.HandleFunc(s.EndPoint, s.CustomHander)
	}
	s.stop = make(chan int)

	ln, err := net.Listen("tcp", "localhost:"+s.Port)
	if err != nil {
		log.Fatalf("Can't listen: %s", err)
	}
	go http.Serve(ln, myMux)
	<-s.stop
	ln.Close()
	log.Println("bye bye..")
}
