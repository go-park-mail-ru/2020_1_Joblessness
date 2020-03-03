package main

import (
	"joblessness/haha/utils/cors"
	"joblessness/haha/utils/routers"
)

func main() {
	cors.Cors.AddOrigin("https://91.210.170.6:8000")
	cors.Cors.AddOrigin("http://91.210.170.6:8000")
	cors.Cors.AddOrigin("http://localhost:8000")
	cors.Cors.AddOrigin("http://localhost:8080")
	cors.Cors.AddOrigin("https://compassionate-wescoff-a0cb89.netlify.com")

	routers.StartRouter()
}