package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"hotelManagement/internal/config"
	"hotelManagement/internal/handlers"
	"hotelManagement/internal/models"
	"hotelManagement/internal/render"

	"github.com/alexedwards/scs/v2"

	"time"
)

const portNumber = ":1234"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))
	//_ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
