package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/charmbracelet/log"
)

// * BLOCKING * //
func AndyProxy(key string, port int) error {
	log.Info("starting andyAPI...", "model", "auto", "port", port)

	target, err := url.Parse("https://andy.mindcraft-ce.com/api/")
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Rewrite = func(r *httputil.ProxyRequest) {
		r.SetURL(target)
		r.Out.Host = target.Host
		r.Out.Header.Set("Content-Type", "application/json")
		if key != "" {
			r.Out.Header.Set("Authorization", "Bearer "+key)
		}
	}

	server := &http.Server{
		Addr:    "127.0.0.1:" + strconv.Itoa(port),
		Handler: proxy,
	}

	log.Info("API server running on http://127.0.0.1:" + strconv.Itoa(port))
	return server.ListenAndServe()
}
