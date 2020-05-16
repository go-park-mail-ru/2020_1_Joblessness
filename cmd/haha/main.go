package main

import (
	"joblessness/haha/middleware"
	"joblessness/haha/server"
)

func main() {
	corsHandler := middleware.NewCorsHandler()

	corsHandler.AddOrigin("http://5.23.54.85")
	corsHandler.AddOrigin("http://localhost:8080")
	corsHandler.AddOrigin("http://localhost:9090")
	corsHandler.AddOrigin("http://localhost:80")
	corsHandler.AddOrigin("http://localhost:8001")
	corsHandler.AddOrigin("http://localhost:8001")
	corsHandler.AddOrigin("http://hahao.ru:443")
	corsHandler.AddOrigin("https://hahao.ru:443")
	corsHandler.AddOrigin("http://hahao.ru:80")
	corsHandler.AddOrigin("https://hahao.ru:80")

	app := server.NewApp(corsHandler)
	app.StartRouter()
}
