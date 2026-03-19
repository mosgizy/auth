package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
	"github.com/mosgizy5/auth-app/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Printf("Port not in environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if portString == "" {
		log.Printf("DB URL not in environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Printf("Error connecting to the database: %v",err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use((cors.Handler((cors.Options{
		AllowedOrigins: []string{"https://*","http://*"},
		AllowedMethods: []string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))))

	router.Use(httprate.Limit(
	10,             
	10*time.Second, 
	httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/register", apiCfg.handlerCreateUser)
	v1Router.Post("/login", apiCfg.handlerLogin)

	//protected routes goes here
	v1Router.Group(func(v1Router chi.Router){
		v1Router.Use(apiCfg.authMiddleWare)
		v1Router.Post("/forget-password",apiCfg.handlerForgetPassword)
		v1Router.Post("/reset-password",apiCfg.handlerResetPassword)
	})
	
		

	router.Mount("/v1",v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":"+portString,
	}

	log.Printf("Server running on port: %s", portString)

	log.Fatal(srv.ListenAndServe())
} 