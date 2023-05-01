package server

import "net/http"

type Response struct {
	Content string
}

func New(address string) *http.Server {
	initRoutes()

	return &http.Server{
		Addr: address,
	}
}
