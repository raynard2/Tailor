package main

import (
	//"Mlops/db"
	"Mlops/routes"

)


func main() {
routes.Dbinit()
	e := routes.New()

	e.Logger.Fatal(e.Start(":1111"))

}
