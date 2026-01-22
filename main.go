package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var product = []Product{
	{ID: 1, Name: "Indomie Rebus", Price: 3000, Stock: 10},
	{ID: 2, Name: "Mie Sedap Jumbo", Price: 4000, Stock: 20},
	{ID: 3, Name: "Nipis Madu", Price: 6000, Stock: 15},
	{ID: 4, Name: "Astor", Price: 2000, Stock: 5},
	{ID: 5, Name: "Le minerale", Price: 2500, Stock: 12},
}

var category = []Category{
	{ID: 1, Name: "Food", Description: "Ini makanan"},
	{ID: 2, Name: "Drink", Description: "Ini minuman"},
	{ID: 3, Name: "Snack", Description: "Ini snack"},
}

// get detail product /api/product/{id}
func getProductById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid product id",
		})
		return
	}

	// find the product using for loop
	for _, p := range product {
		if p.ID == id {
			json.NewEncoder(w).Encode(&p)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "product not found",
	})
}

// update product /api/product/{id}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid product id",
		})
		return
	}

	var updateProduct Product

	err = json.NewDecoder(r.Body).Decode(&updateProduct)

	// handle if request is invalid
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid request",
		})
		return
	}

	// for loop to find the product and update it with new body
	for i := range product {
		if product[i].ID == id {
			updateProduct.ID = id
			product[i] = updateProduct

			json.NewEncoder(w).Encode(updateProduct)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "product not found",
	})
}

// delete product /api/product/{id}
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid product id",
		})
		return
	}

	// find the product
	for i, p := range product {
		if p.ID == id {

			// create new slice with the target id removed
			product = append(product[:i], product[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "delete success",
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "product not found",
	})
}

// get detail category /api/categories/{id}
func getCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid category id",
		})
		return
	}

	// find the category using for loop
	for _, p := range category {
		if p.ID == id {
			json.NewEncoder(w).Encode(&p)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "category not found",
	})
}

// update category /api/categories/{id}
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid category id",
		})
		return
	}

	var updateCategory Category

	err = json.NewDecoder(r.Body).Decode(&updateCategory)

	// handle if request is invalid
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid request",
		})
		return
	}

	// for loop to find the product and update it with new body
	for i := range category {
		if category[i].ID == id {
			updateCategory.ID = id
			category[i] = updateCategory

			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "category not found",
	})
}

// delete category /api/categories/{id}
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid category id",
		})
		return
	}

	// find the category
	for i, p := range category {
		if p.ID == id {

			// create new slice with the target id removed
			category = append(category[:i], category[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "delete success",
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "category not found",
	})
}

func main() {
	// endpoint /api/product/{id}
	// GET detail product
	// PUT product
	// DELETE product
	http.HandleFunc("/api/product/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			getProductById(w, r)
		} else if r.Method == "PUT" {
			updateProduct(w, r)
		} else if r.Method == "DELETE" {
			deleteProduct(w, r)
		}
	})

	// endpoint /api/product
	// GET all product
	// POST create product
	http.HandleFunc("/api/product", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			json.NewEncoder(w).Encode(&product)
		} else if r.Method == "POST" {
			var newProduct Product

			err := json.NewDecoder(r.Body).Decode(&newProduct)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "invalid request",
				})
				return
			}

			newProduct.ID = len(product) + 1
			product = append(product, newProduct)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newProduct)
		}
	})

	// endpoint /api/categories/{id}
	// GET category
	// PUT category
	// DELETE category
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			getCategoryById(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	// endpoint /api/categories
	// GET all category
	// POST create endpoint
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			json.NewEncoder(w).Encode(&category)
		} else if r.Method == "POST" {
			var newCategory Category

			err := json.NewDecoder(r.Body).Decode(&newCategory)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "invalid request",
				})
				return
			}

			newCategory.ID = len(category) + 1
			category = append(category, newCategory)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newCategory)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "api running...",
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Hello world",
		})
	})

	fmt.Println("Server running in localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to run server")
	}
}
