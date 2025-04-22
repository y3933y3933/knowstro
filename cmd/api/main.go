package main

const version = "1.0.0"

type application struct {
}

func main() {
	app := application{}

	router := app.routes()

	router.Run()

}
