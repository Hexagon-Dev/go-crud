package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Hexagon-Dev/go-crud/common"
	"net/http"
	"strconv"
)

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
	BarCode     string `json:"bar_code"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
}

func IndexProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request) any {
	return func(w http.ResponseWriter, r *http.Request) any {
		rows, err := db.Query("SELECT * FROM products")
		if err != nil {
			return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
		}

		var products []Product

		for rows.Next() {
			var product Product

			err := rows.Scan(&product.Id, &product.Name, &product.IsAvailable, &product.BarCode, &product.Category, &product.CreatedAt)
			if err != nil {
				return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
			}

			products = append(products, product)
		}

		return products
	}
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

func UpdateProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request) any {
	return func(w http.ResponseWriter, r *http.Request) any {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return common.HttpError{Message: err.Error(), StatusCode: http.StatusUnprocessableEntity}
		}

		var product Product

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
		}

		_, err = db.Exec(
			"UPDATE products SET name = ?, is_available = ?, bar_code = ?, category = ?, created_at = ? WHERE id = ?",
			product.Name,
			product.IsAvailable,
			product.BarCode,
			product.Category,
			product.CreatedAt,
			id,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return common.HttpError{Message: err.Error(), StatusCode: http.StatusNotFound}
			} else {
				return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
			}
		}

		product.Id = id

		return product
	}
}

func DeleteProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request) any {
	return func(w http.ResponseWriter, r *http.Request) any {
		exec, err := db.Exec("DELETE FROM products WHERE id = ?", r.PathValue("id"))
		if err != nil {
			return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
		}

		rowsAffected, err := exec.RowsAffected()
		if err != nil {
			return common.HttpError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
		}

		if rowsAffected < 1 {
			return common.HttpError{Message: "Product not found.", StatusCode: http.StatusNotFound}
		}

		return nil
	}
}
