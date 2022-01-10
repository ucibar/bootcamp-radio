package main

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/uCibar/bootcamp-radio/api/handler"
	"github.com/uCibar/bootcamp-radio/config"
	"github.com/uCibar/bootcamp-radio/repository"
	"github.com/uCibar/bootcamp-radio/service/auth"
	"github.com/uCibar/bootcamp-radio/service/user"
	"log"
	"net/http"
	"os"
)

func main() {
	conf := config.FromENV()

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	httpLogger := log.New(os.Stdout, "HTTP", log.Ldate|log.Ltime)

	userRepo := repository.NewUserPostgresRepository(db)
	userService := user.NewService(userRepo)

	authService := auth.NewService(userService)

	broadcastRepository := repository.NewSessionRepository()

	authHandler := handler.NewAuthHandler(authService, httpLogger)

	broadcastHandler := handler.NewBroadcastHandler(broadcastRepository, httpLogger)

	router := httprouter.New()

	router.POST("/api/auth/login", authHandler.Login)
	router.POST("/api/auth/register", authHandler.Register)

	router.GET("/broadcast/list", authHandler.Middleware(broadcastHandler.List))
	router.GET("/broadcast/info/:broadcast_id", authHandler.Middleware(broadcastHandler.Info))
	router.POST("/broadcast/create", authHandler.Middleware(broadcastHandler.Create))
	router.POST("/broadcast/join", authHandler.Middleware(broadcastHandler.Join))

	router.ServeFiles("/public/*filepath", http.Dir("public"))

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	fmt.Println("server is running on port 8085")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
