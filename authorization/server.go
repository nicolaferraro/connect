package authorization

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Server struct {
	Code  chan string
	State string
}

func NewServer(state string) *Server {
	return &Server{
		Code:  make(chan string, 10),
		State: state,
	}
}

func (s *Server) Start() {
	http.HandleFunc("/callback", s.callback)

	println("Server listening on port 3000...")
	println("Callback URL: http://localhost:3000/callback")
	go func() {
		if err := http.ListenAndServe("localhost:3000", nil); err != nil {
			fmt.Sprintf("Error while starting server: %v", err)
			os.Exit(1)
		}
	}()
}

func (s *Server) callback(w http.ResponseWriter, r *http.Request) {
	parsed, err := url.Parse(r.RequestURI)
	if err != nil {
		fmt.Printf("error %v", err)
		w.WriteHeader(400)
		return
	}

	code := parsed.Query().Get("code")
	recvState := parsed.Query().Get("state")

	if code == "" {
		w.WriteHeader(400)
		return
	}

	if recvState != s.State {
		fmt.Printf("wrong state: %s != %s", recvState, s.State)
		w.WriteHeader(403)
		return
	}

	s.Code <- code
	w.Write([]byte("Success! You can close this page"))
}
