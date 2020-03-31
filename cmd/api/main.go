package main

import (
	"joblessness/haha/middleware"
	"joblessness/haha/server"
)

func main() {
	corsHandler := middleware.NewCorsHandler()

	corsHandler.AddOrigin("https://91.210.170.6:8000")
	corsHandler.AddOrigin("http://91.210.170.6:8000")
	corsHandler.AddOrigin("http://localhost:8000")
	corsHandler.AddOrigin("http://localhost:8080")
	corsHandler.AddOrigin("https://compassionate-wescoff-a0cb89.netlify.com")

	app := server.NewApp(corsHandler)
	app.StartRouter()
}