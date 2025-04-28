package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin/binding"
	a "github.com/y3933y3933/knowstro/internal/app"
	"github.com/y3933y3933/knowstro/internal/routes"
)

func main() {
	// gin decoder config
	binding.EnableDecoderDisallowUnknownFields = true

	app, err := a.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	r := routes.SetupRoutes(app)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Info("starting server", "addr", srv.Addr, "env", app.Config.Env)

	err = srv.ListenAndServe()

	if err != nil {
		app.Logger.Error(err.Error())
		os.Exit(1)
	}

}
