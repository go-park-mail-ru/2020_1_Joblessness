package main

var Cors = CorsHandler{
	allowedOrigins: []string{},
}

func main() {
	Cors.AddOrigin("https://91.210.170.6:8000")
	Cors.AddOrigin("http://91.210.170.6:8000")
	Cors.AddOrigin("http://localhost:8000")
	Cors.AddOrigin("http://localhost:8080")
	Cors.AddOrigin("https://compassionate-wescoff-a0cb89.netlify.com")



	StartRouter()
}