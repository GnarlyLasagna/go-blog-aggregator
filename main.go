package main

import (
    "fmt"
    "log"
    "os"
    "net/http"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"
    "github.com/GnarlyLasagna/go-blog-aggregator/internal/database"
)


func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port not found")
	}

    dbURL := os.Getenv("DB_URL")
	if DbURL == "" {
		log.Fatal("Database URL is not found in the environment")
	}

    conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
    v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)	

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("server starting on %v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port:", portString)

}
