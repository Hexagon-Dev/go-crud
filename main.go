package main

import (
	"database/sql"
	"encoding/json"
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

	defer db.Close()

	if err != nil {
		log.Panicln(err)
	}

	routes := []Route{
		{"GET /product/{id}", controllers.GetProduct(db)},
		{"POST /product", controllers.CreateProduct(db)},
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
