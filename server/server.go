package server

import "net/http"

func New(address string) *http.Server {
	return &http.Server{
		Addr: address,
	}
}
