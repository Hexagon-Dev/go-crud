package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Hexagon-Dev/go-crud/common"
	"net/http"
)

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
	BarCode     string `json:"bar_code"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
}

func GetProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request) any {
	return func(w http.ResponseWriter, r *http.Request) any {
		id := r.PathValue("id")

		row := db.QueryRow("SELECT * FROM products WHERE id = ?", id)

		var product Product

		err := row.Scan(&product.Id, &product.Name, &product.IsAvailable, &product.BarCode, &product.Category, &product.CreatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return common.HttpError{Message: err.Error(), StatusCode: http.StatusNotFound}
			} else {
				return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
			}
		}

		return product
	}
}

func CreateProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request) any {
	return func(w http.ResponseWriter, r *http.Request) any {
		var product Product

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
		}

		result, err := db.Exec(
			"INSERT INTO products (name, is_available, bar_code, category, created_at) VALUES (?, ?, ?, ?, ?)",
			product.Name,
			product.IsAvailable,
			product.BarCode,
			product.Category,
			product.CreatedAt,
		)
		if err != nil {
			return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
		}

		id, err := result.LastInsertId()
		if err != nil {
			return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
		}

		product.Id = int(id)

		return product
	}
}
