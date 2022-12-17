package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"necross.it/backend/auth"
	"necross.it/backend/auth/user"
	"necross.it/backend/database"
	"necross.it/backend/settings"
	"necross.it/backend/util"
	"necross.it/backend/verify"
	"os"
	"time"
)

type App struct {
	DB         database.Server
	Settings   settings.Settings
	User       user.Service
	Auth       auth.Service
	Verify     verify.Service
	Permission user.Service
	General    util.Service
}

func main() {
	start := time.Now()
	log.Println("Starting application...")
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseConnection := database.InitializeDBConnection()

	err = database.CreateTables(databaseConnection)
	if err != nil {
		log.Panic(err.Error())
	}

	generalService := util.Service{Server: databaseConnection}
	userService := user.Service{Server: databaseConnection}
	verifyService := verify.Service{
		Server: databaseConnection,
	}
	permissionService := user.Service{Server: databaseConnection}
	authService := auth.Service{
		Server:  databaseConnection,
		User:    userService,
		Verify:  verifyService,
		General: generalService,
	}

	appSettings := settings.Settings{}

	app := App{
		DB:         databaseConnection,
		Settings:   appSettings,
		User:       userService,
		Auth:       authService,
		General:    generalService,
		Permission: permissionService,
		Verify:     verifyService,
	}

	server := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 500,
	})
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	api := server.Group("/api")        // /api
	v1 := api.Group("/v1")             // /api/v1
	userGroup := v1.Group("user")      // /api/v1/user
	verifyGroup := v1.Group("/verify") // /api/v1/verify

	app.Auth.RegisterRoutes(v1)            // /api/v1
	app.Auth.RegisterUserRoutes(userGroup) // /api/v1/user
	app.Verify.RegisterRoutes(verifyGroup) // /api/v1/verify

	elapsed := time.Since(start)
	log.Printf("Application started in %s", elapsed)

	backendPort := os.Getenv("BACKEND_PORT")
	log.Fatal(server.Listen(":" + backendPort))
}
