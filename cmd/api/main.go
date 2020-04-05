package main

import (
	"joblessness/haha/middleware"
	"joblessness/haha/server"
)

func main() {
	corsHandler := middleware.NewCorsHandler()

	corsHandler.AddOrigin("https://5.23.54.85:80")
	corsHandler.AddOrigin("http://5.23.54.85:80")
	corsHandler.AddOrigin("http://localhost:8080")
	corsHandler.AddOrigin("http://localhost:80")
	corsHandler.AddOrigin("https://compassionate-wescoff-a0cb89.netlify.com")

	app := server.NewApp(corsHandler)
	app.StartRouter()
}