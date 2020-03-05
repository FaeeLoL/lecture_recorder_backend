package main

import "github.com/faeelol/lecture_recorder_backend/routes"

func main() {
	router := routes.InitRoutes()
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
