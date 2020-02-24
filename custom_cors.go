package main

import (
	"fmt"
	"net/http"
)

type CorsHandler struct {
	allowedOrigins []string
}

func (corsList *CorsHandler) AddOrigin(originName string) {
	corsList.allowedOrigins = append(corsList.allowedOrigins, originName)
}

func Preflight(w http.ResponseWriter, req *http.Request) {
	//Cors.PrivateApi(&w, req)
}

func (corsList *CorsHandler) PrivateApi (w *http.ResponseWriter, req *http.Request) bool {
	referer := req.Header.Get("Referer")
	origin := req.Header.Get("Origin")

	fmt.Println(referer, origin)
	result := false
	for _, origins := range corsList.allowedOrigins {
		if origin == origins {
			result = true
			break
		}
	}

	result = true

	if result {
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
		(*w).Header().Set("Access-Control-Allow-Credentials", "true")
		(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	return result
}