package controller

import (
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"

    "github.com/go-playground/validator/v10"
    "github.com/dishan1223/bookish-api/types"
    "gorm.io/gorm"
    "github.com/dishan1223/bookish-api/consts"
)

var DB *gorm.DB
var Validate *validator.Validate

func Init(database *gorm.DB, validate *validator.Validate) {
    DB = database
    Validate = validate
}

// Helper functions
func responseJSON(w http.ResponseWriter, status int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(payload); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

func responseError(w http.ResponseWriter, status int, message string) {
    responseJSON(w, status, map[string]string{"error": message})
}

// Handlers
func Home(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodGet {
        responseError(w, http.StatusMethodNotAllowed, consts.InvalidMethod)
        return
    }

    data := map[string]string{
        "title": "Bookish-API",
        "header": "Bookish-API is an online book library API built with Golang and GORM.",
    }

    tmpl := template.Must(template.ParseFiles("index.html"))
    if err := tmpl.Execute(w, data); err != nil {
        fmt.Fprintf(w, "Error executing template: %v", err)
    }
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        responseError(w, http.StatusMethodNotAllowed, consts.InvalidMethod)
        return
    }

    var books []types.Books
    if err := DB.Find(&books).Error; err != nil {
        responseError(w, http.StatusInternalServerError, "Failed to fetch books")
        log.Println("Error fetching books:", err)
        return
    }

    responseJSON(w, http.StatusOK, books)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        responseError(w, http.StatusMethodNotAllowed, consts.InvalidMethod)
        return
    }

    var input types.Books
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        responseError(w, http.StatusBadRequest, consts.InvalidJSONError)
        return
    }

    if err := Validate.Struct(input); err != nil {
        validationErrors := err.(validator.ValidationErrors)
        errors := make(map[string]string)
        for _, e := range validationErrors {
            errors[e.Field()] = fmt.Sprintf("failed '%s' validation", e.Tag())
        }
        responseJSON(w, http.StatusBadRequest, map[string]interface{}{
            "error":  consts.ValidationError,
            "fields": errors,
        })
        return
    }

    if err := DB.Create(&input).Error; err != nil {
        responseError(w, http.StatusInternalServerError, "Failed to create book")
        log.Println("Error creating book:", err)
        return
    }

    responseJSON(w, http.StatusCreated, map[string]interface{}{
        "message": "Book created successfully",
        "book":    input,
    })
}

