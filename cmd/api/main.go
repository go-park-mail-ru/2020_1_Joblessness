package main

import (
	"github.com/kataras/golog"
	"joblessness/haha/server"
	"joblessness/haha/utils/cors"
)

func main() {
	golog.Error("1")
	corsHandler := cors.NewCorsHandler()
	golog.Error("2")
	corsHandler.AddOrigin("https://91.210.170.6:8000")
	corsHandler.AddOrigin("http://91.210.170.6:8000")
	corsHandler.AddOrigin("http://localhost:8000")
	corsHandler.AddOrigin("http://localhost:8080")
	corsHandler.AddOrigin("https://compassionate-wescoff-a0cb89.netlify.com")
	golog.Error("3")

	app := server.NewApp(corsHandler)
	app.StartRouter()
}