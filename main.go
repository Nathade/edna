package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/justinas/alice" //A middleware chaining library.

	"github.com/rs/cors"
	"github.com/skratchdot/open-golang/open"
)

const (
	//Cost is the, well, cost of the bcrypt encryption used for storing user
	//passwords in the database. It determines the amount of processing power to
	// be used while hashing and saalting the password. The higher, the cost,
	//the more secure the password hash, and also the more cpu cycles used for
	//password related processes like comparing hasshes during authentication
	//or even hashing a new password.
	Cost int = 5
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	//REDISADDR, REDISPW, MONGOSERVER, MONGODB, Public, Private, RootURL, AWSBucket := checks()

	config := generateConfig()
	defer config.MongoSession.Close()
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)
	router := NewRouter()

	//router.Post("/api/v0.1/auth", commonHandlers.ThenFunc(appC.authHandler))
	router.Post("/api/auth/login", commonHandlers.Append(dbsetter).ThenFunc(config.LoginPost))

	router.Post("/api/staff", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.createUserHandler))
	router.Put("/api/staff", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.updateUserHandler))
	router.Get("/api/staff", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.getUsersHandler))

	router.Post("/api/class", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.createClassHandler))
	router.Get("/api/class", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.getClassesHandler))
	router.Put("/api/class", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.putClassHandler))

	router.Post("/api/subject", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.createSubjectHandler))
	router.Get("/api/subject", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.getSubjectsHandler))
	router.Put("/api/subject", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.putSubjectHandler))

	router.Get("/api/me", commonHandlers.Append(dbsetter, config.frontAuthHandler).ThenFunc(config.getMeHandler))
	router.Post("/register.html", commonHandlers.ThenFunc(config.NewSchool))
	router.Get("/verify", commonHandlers.ThenFunc(config.VerifySchool))
	router.Get("/", commonHandlers.Append(dbsetter).ThenFunc(config.RootHandler))

	router.HandleMethodNotAllowed = false
	router.NotFound = http.FileServer(http.Dir("./static")).ServeHTTP

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Println("No Global port has been defined, using default port :8080")

		PORT = "8080"

	}

	handler := cors.New(cors.Options{
		//		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedOrigins: []string{"*"},

		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Auth-Token", "*"},
		Debug:            false,
	}).Handler(router)
	log.Println("serving ")

	open.Run("http://localhost:" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
