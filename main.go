package main

var Cors = CorsHandler{
	allowedOrigins: []string{},
}

func main() {

	Cors.AddOrigin("http://localhost:8080")
	Cors.AddOrigin("https://compassionate-wescoff-a0cb89.netlify.com")

	StartRouter()
}