package main

import (
	//"Mlops/db"
	"Mlops/db"
	"Mlops/routes"

)


func main() {
	db.DatabaseInt()
	e := routes.New()

	e.Logger.Fatal(e.Start(":8080"))

}
