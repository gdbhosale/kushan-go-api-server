package main

import (
	"goat/internal"
	"goat/internal/http"
	"goat/internal/pgx"
	"strconv"

	"context"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

// Init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		internal.Error("Main::init", "No .env file found", err)
	}
}

// @title			GOAT Admin
// @version		1.0.0
// @description	This is a API Documentation for GOAT Admin which defines all its API's for Angular Frontend.
// @description
// @description This API provides endpoints to manage users, RBAC, customers and many other standard Modules of CRM (Customer Relationship Management) system. It allows users to create, retrieve, update, and delete various records associated with customers.
// @description
// @description It needs JWT Token based authentication for user access, Modular RBAC, validation of input data, and error handling to ensure data integrity and security.
// @description
// @description The API is designed to be RESTful, using standard HTTP methods (GET, POST, PUT, DELETE) for CRUD operations on resources. Responses are formatted as JSON.
// @termsOfService
// @contact.name				GOAT Support
// @contact.url				https://goatadmin/support
// @contact.email				support@goatadmin.com
// @license.name				MIT
// @license.url				https://opensource.org/license/mit
// @host						localhost:8083
// @BasePath					/
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Enter the token received from Signin API with the `Bearer ` prefix, e.g. "Bearer ABcDE12345"
func main() {

	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Instantiate a new type to represent our application.
	// This type lets us shared setup code with our end-to-end tests.
	main := NewMain()

	// Run program
	if err := main.Run(ctx); err != nil {
		internal.Error("Main", "error while running Main API", err)
		main.Close()
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := main.Close(); err != nil {
		internal.Error("Main", "Error while closing Main API", err) // fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Main struct {
	// SQLite database
	DB *sqlx.DB

	// Port
	Port int

	// HTTP server for handling HTTP communication.
	// SQLite services are attached to it before running.
	HTTPServer *http.Server

	// Services
	UserService internal.UserService
}

func NewMain() *Main {
	// Get DB Details from env

	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_username := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	// Initiate DB
	db, err := sqlx.Open("postgres", "postgres://"+db_username+":"+db_password+"@"+db_host+":"+db_port+"/"+db_name+"?sslmode=disable")
	if err != nil {
		panic(err)
	}

	// Postgres Driver
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	runMigrationsAndSeeds := false

	if runMigrationsAndSeeds {
		// Migrations Setup
		migrationsDir := filepath.Join("internal/pgx/migrations")
		migrations, err := migrate.NewWithDatabaseInstance(
			"file://"+migrationsDir,
			"postgres", driver)
		if err != nil {
			panic(err)
		}

		// Run Migrations
		migrations.Down()
		migrations.Up()

		// Data Seeding
		err = pgx.Seed(db, "internal/pgx/seeders/users.sql")
		if err != nil {
			return nil
		}
	}

	// Port. To be taken from config
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		internal.Error("Main::NewMain", "Env doesn't have PORT", err)
		port = 8083 // Assign default port
	}

	// Create Main Object
	return &Main{
		DB:         db,
		HTTPServer: http.NewServer(port),
		Port:       port,
	}
}

// Run executes the main program
func (main *Main) Run(ctx context.Context) (err error) {

	// Initiate Services
	userService := pgx.NewUserService(main.DB)

	// Attach user service to Main for testing.
	main.UserService = userService

	// Attach underlying services to the HTTP server.
	main.HTTPServer.UserService = userService

	// Start Server
	go func() { main.HTTPServer.ListenAndServe() }()

	return nil
}

// Close gracefully stops the program.
func (main *Main) Close() error {
	internal.Warn("Main::Close", "Closing Server...")
	if main.HTTPServer != nil {
		if err := main.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if main.DB != nil {
		if err := main.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}
