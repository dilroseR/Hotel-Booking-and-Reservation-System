package main

import (
	"encoding/gob"
	"fmt"
	"hotelManagement/internal/config"
	"hotelManagement/internal/driver"
	"hotelManagement/internal/handlers"
	"hotelManagement/internal/helpers"
	"hotelManagement/internal/models"
	"hotelManagement/internal/render"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"

	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User {})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	 
	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=rose password=")
	if err!=nil{
		log.Fatal("Cannot connect to database...")
	}
	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err!=nil{
		log.Fatal("cannot create template cache ")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer (&app)
	helpers.NewHelpers(&app)

	return db, nil
}
