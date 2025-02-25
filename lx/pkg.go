// Package lx expands on the [net/http] package by providing a Server type and
// convenience functions, improving developer experience when creating API
// services.
//
// lx also provides structured logging for HTTP requests and operations.
//
// Quick start:
//
//	type ServerAdapter struct {
//	  *lx.Server
//	}
//
//	func NewServerAdapter(opts ...lx.ServerOptsFunc) *ServerAdapter {
//	  return &ServerAdapter{
//	    Server: lx.NewServer(opts...),
//	  }
//	}
//
//	type HelloResponse struct {
//	  Message string `json:"message"`
//	}
//
//	func (s *ServerAdapter) SayHello(w http.ResponseWriter, r *http.Request) {
//	  s.Logger.Info("Saying hello from ServerAdapter.SayHello handler!")
//	  lx.WriteJson(w, HelloResponse{
//	    Message: "What's up!",
//	  })
//	}
//
//	func main() {
//	  server := NewServerAdapter(lx.WithAddr("localhost:3000"))
//	  server.HandleFunc("GET /", server.SayHello)
//	  server.ListenAndServe()
//	}
package lx
