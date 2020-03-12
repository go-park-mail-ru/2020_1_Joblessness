package main

import (
	"joblessness/haha/server"
	"joblessness/haha/utils/cors"
)

func main() {
	cors.Cors.AddOrigin("https://91.210.170.6:8000")
	cors.Cors.AddOrigin("http://91.210.170.6:8000")
	cors.Cors.AddOrigin("http://localhost:8000")
	cors.Cors.AddOrigin("http://localhost:8080")
	cors.Cors.AddOrigin("https://compassionate-wescoff-a0cb89.netlify.com")

	app := server.NewApp()
	app.StartRouter()
}