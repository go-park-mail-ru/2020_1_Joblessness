package main

import (
	"./haha/server"
)

func main() {

	haha := server.NewServer()
	haha.AddOrigin("https://91.210.170.6:8000")
	haha.AddOrigin("http://91.210.170.6:8000")
	haha.AddOrigin("http://localhost:8000")
	haha.AddOrigin("http://localhost:8080")
	haha.AddOrigin("https://compassionate-wescoff-a0cb89.netlify.com")

	haha.StartRouter()
}