package application

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/papijo/lms/db/postgresdb"
	"github.com/papijo/lms/utils/config"
	"github.com/papijo/lms/utils/middlewares"
	"gorm.io/gorm"
)

// Starts the Application

func Start() (*echo.Echo, *gorm.DB) {

	//Load Environment Variables
	envErr := config.LoadEnvironmentVariables()
	if envErr != nil {
		log.Fatal("Error loading environment variable file: ", envErr)
	}

	//Database Operations
	//Connect to the Database (PostgreSQL)
	db, err := postgresdb.InitializeDB()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	//Migrate the Database Schema

	err = postgresdb.MigrateSchema(db)
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}

	//Initialise Echo Instance
	e := echo.New()

	//Home Page of Server
	e.GET("/", func(c echo.Context) error {
		response := echo.Map{
			"Application Name":     "Learning Management System",
			"Application Owner":    "Industrial Training Fund Model Skills Training Centre",
			"Application Version":  "1.0.0",
			"Application Engineer": "Ebhota Jonathan",
		}

		return c.JSON(http.StatusOK, response)
	})

	// Middlewares
	// Logger Middleware
	e.Use(middlewares.LoggerConfigMiddleware())

	//Create Server Instance
	s := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	e.HideBanner = true

	// Start Server Instance
	go func() {
		if err := e.StartServer(s); err != nil {
			log.Println(err.Error(), "Shutting down the Server")
		}
	}()

	log.Println("⚡️🚀 LMS Server Started")

	return e, db
}

// Gracefully shuts down the Application

func Stop(e *echo.Echo, db *gorm.DB) error {
	time.Sleep(1 * time.Second)
	log.Println("⚡️🚀 LMS Server - Stopping")
	time.Sleep(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	//Close the DB Connection

	defer func() {
		//Close the Database Connection
		if err := postgresdb.CloseDB(db); err != nil {
			log.Fatalf("Failed to close Database Connection: %s", err.Error())
		}
		log.Println("⚡️🚀 Database Connection Closed")
	}()

	//Shut down the database
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}
