package types

import (
    "gorm.io/gorm"
)

// Server struct
type Server struct {
    Addr string
}

// Book model
type Books struct {
    gorm.Model
    Title       string  `json:"title" validate:"required,min=2,max=100"`
    Description string  `json:"description" validate:"max=1000"`
    Poster      string  `json:"poster" validate:"required,url"`
    Pages       uint    `json:"pages" validate:"required,gt=0"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Author      string  `json:"author" validate:"required,min=2,max=100"`
    Link        string  `json:"link" validate:"url"`
    DownloadLink string  `json:"download_link" validate:"url"`
}

