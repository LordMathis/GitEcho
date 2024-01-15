package webhooks

import (
	"errors"
	"net/http"
)

type WebhookServer struct {
	*http.Server
	apiPath string
}

func NewWebhookServer(addr string) *WebhookServer {

	return &WebhookServer{
		&http.Server{
			Addr:    addr,
			Handler: http.NewServeMux(),
		},
		"/api/v1/webhooks",
	}
}

func (ws *WebhookServer) RegisterWebhookHandler(repo_name string, f func(http.ResponseWriter, *http.Request)) error {

	// Register the handler function for the specified pattern
	mux, ok := ws.Handler.(*http.ServeMux)
	if !ok {
		return errors.New("server handler not defined")
	}

	mux.HandleFunc(ws.apiPath+"/"+repo_name, f)
	return nil
}
