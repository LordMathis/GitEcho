package webhooks

import (
	"net/http"
)

type WebhookServer struct {
	*http.Server
}

func NewWebhookServer(addr string) *WebhookServer {

	return &WebhookServer{
		&http.Server{
			Addr: addr,
		},
	}
}

func (ws *WebhookServer) RegisterWebhookHandler(repo_name string, f func(http.ResponseWriter, *http.Request)) {
	if ws.Handler == nil {
		ws.Handler = http.NewServeMux()
	}

	// Register the handler function for the specified pattern
	mux, ok := ws.Handler.(*http.ServeMux)
	if !ok {
		// If the current handler is not a ServeMux, create a new one
		mux = http.NewServeMux()
		ws.Handler = mux
	}

	mux.HandleFunc("/"+repo_name, f)
}
