package main

import "github.com/SoroushBeigi/workout-go/internal/app"

func main() {
	app, err := app.NewApplication()
	if err!=nil{
		panic(err)
	}
	app.Logger.Println("running")
}