package main

import (
	_cors "../utils/cors"
	_routers "../utils/routers"
)

func main() {
	_cors.Cors.AddOrigin("https://91.210.170.6:8000")
	_cors.Cors.AddOrigin("http://91.210.170.6:8000")
	_cors.Cors.AddOrigin("http://localhost:8000")
	_cors.Cors.AddOrigin("http://localhost:8080")
	_cors.Cors.AddOrigin("https://compassionate-wescoff-a0cb89.netlify.com")



	_routers.StartRouter()
}