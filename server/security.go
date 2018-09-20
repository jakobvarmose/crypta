package server

import (
	"fmt"
	"net/http"
	"net/url"
)

func securityHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println(req.RequestURI)
		resp.Header().Set("X-Frame-Options", "DENY")
		resp.Header().Set("X-Content-Type-Options", "nosniff")
		handler.ServeHTTP(resp, req)
	})
}

func securityCheck(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		originString := req.Header.Get("Origin")
		referrerString := req.Header.Get("Referer")
		if originString == "" && referrerString == "" {
			http.Error(resp, "403 forbidden - no origin or referrer", 403)
			return
		}
		if originString != "" {
			origin, err := url.ParseRequestURI(originString)
			if err != nil {
				http.Error(resp, "403 forbidden - invalid origin", 403)
				return
			}
			if origin.Host != req.Host {
				http.Error(resp, "403 forbidden - wrong origin", 403)
				return
			}
		}
		if referrerString != "" {
			referrer, err := url.ParseRequestURI(referrerString)
			if err != nil {
				http.Error(resp, "403 forbidden - invalid referrer", 403)
				return
			}
			if referrer.Host != req.Host {
				http.Error(resp, "403 forbidden - wrong referrer", 403)
				return
			}
		}
		handler.ServeHTTP(resp, req)
	})
}
