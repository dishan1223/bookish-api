package main

import (
    "log"
    "net/http"
    "github.com/go-playground/validator/v10"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/dishan1223/bookish-api/types"
    "github.com/dishan1223/bookish-api/controller"
    "github.com/dishan1223/bookish-api/middleware"
)

var DB *gorm.DB
var validate *validator.Validate


// Entry point
func main() {
    validate = validator.New()

    var err error
    DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
    if err != nil {
        panic("Failed to connect database")
    }

    if err := DB.AutoMigrate(&types.Books{}); err != nil {
        panic("Failed to migrate table")
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

