package main

import (
    "log"
    "os"
    "net/http"
    "github.com/go-playground/validator/v10"
    "gorm.io/gorm"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "github.com/dishan1223/bookish-api/types"
    "github.com/dishan1223/bookish-api/controller"
    "github.com/dishan1223/bookish-api/middleware"
)


// Globals
var DB *gorm.DB
var validate *validator.Validate


// Entry point
func main() {
    validate = validator.New()

    // load environment variables
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // DB CONNECTION STRING
    dsn := os.Getenv("DB_CONNECTION_STRING")
    if dsn == "" {
        log.Fatal("DB_CONNECTION_STRING is not set")
    }

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) 
    if err != nil {
        panic("Failed to connect to database")
    }

    // migrate the schema
    if err := DB.AutoMigrate(&types.Books{}); err != nil {
       panic("Failed to migrate schema") 
    }

    // initializing the controller
    controller.Init(DB, validate)

    // routes handlers
    http.HandleFunc("/", middleware.Logger(controller.Home))
    http.HandleFunc("/api/v1/books", middleware.Logger(controller.GetBooks))
    http.HandleFunc("/api/v1/books/add", middleware.Logger(controller.AddBook))

    // initializing the server
    s := types.Server{"127.0.0.1:3000"}
    log.Println("Server is running on http://" + s.Addr)
    log.Fatal(http.ListenAndServe(s.Addr, nil))
}

