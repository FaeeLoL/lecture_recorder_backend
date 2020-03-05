package main

import (
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/routes"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	router := routes.InitRoutes()
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
