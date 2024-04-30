package main

import "net/http"

type Api struct {
	storage Storage
	address string
}

func NewApi(storage Storage,address string) *Api {
	return &Api{
		storage: storage,
		address: address,
	}
}

func (api *Api) Run() {
	server := http.NewServeMux()

	server.HandleFunc("/", func (res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("some-some"))
		return
	})

	http.ListenAndServe(api.address, server)
}