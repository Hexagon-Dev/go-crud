package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Hexagon-Dev/go-crud/common"
	"github.com/Hexagon-Dev/go-crud/controllers"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
)

type Route struct {
	Pattern     string
	HandlerFunc func(w http.ResponseWriter, r *http.Request) any
}

func main() {
	db, err := sql.Open("sqlite", "go.sqlite")
	if err != nil {
		log.Panicln(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Failed to close database connection.")
		}
	}(db)

	routes := []Route{
		{"GET /product/{id}", controllers.GetProduct(db)},
		{"POST /product", controllers.CreateProduct(db)},
		{"PUT /product/{id}", controllers.UpdateProduct(db)},
		{"DELETE /product/{id}", controllers.DeleteProduct(db)},
	}

	for _, route := range routes {
		http.HandleFunc(route.Pattern, toJson(func(w http.ResponseWriter, r *http.Request) any {
			return route.HandlerFunc(w, r)
		}))
	}

	log.Fatal(http.ListenAndServe(":80", nil))
}

func toJson(f func(w http.ResponseWriter, r *http.Request) any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		obj := f(w, r)

		encoded, err := json.Marshal(obj)
		if err != nil {
			http.Error(w, "Failed to encode JSON.", http.StatusInternalServerError)
		}

		switch errObj := obj.(type) {
		case common.HttpError:
			http.Error(w, string(encoded), errObj.StatusCode)
		default:
			_, _ = w.Write(encoded)
		}
	}
}
