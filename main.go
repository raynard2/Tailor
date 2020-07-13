package main

import (
	"Mlops/routes"

)


func main() {
	e := routes.New()

	e.Logger.Fatal(e.Start(":1111"))

}
