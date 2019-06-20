package websocket

import (
	"log"
	"net/http"
)

type Server struct {
	Addr	string
	Uri		string
	Handler Handler
}

func init()  {
	hub := GetHub()
	var handler Handler = &DefaultHandler{}
	hub.SetHandler(&handler)
	go hub.run()
}

func (server *Server) Startup() error {

	if server.Handler == nil {
		log.Panic("the message handler not set")
	}

	hub.SetHandler(&server.Handler)

	if server.Uri == "" {
		server.Uri = "/"
	}

	if server.Addr == "" {
		server.Addr = "0.0.0.0:8771"
	}

	http.HandleFunc(server.Uri, func(w http.ResponseWriter, r *http.Request) {
		NewServerHandler(w, r, hub)
	})
	go http.ListenAndServe(server.Addr,nil)

	log.Println("Websocket server startup success!")
	return nil
}

