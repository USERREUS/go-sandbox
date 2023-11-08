package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

type Inventory struct {
	ProductCode     string  `json:"product_code"`
	ProductName     string  `json:"product_name"`
	QuantityInStock float64 `json:"quantity_in_stock"`
	CurrentPrice    string  `json:"current_price"`
}

var db *sql.DB

func main() {
	// Инициализация базы данных
	initDB()

	// Создание маршрутизатора gorilla/mux
	router := mux.NewRouter()

	// Определение обработчиков для эндпоинтов
	router.HandleFunc("/inventory/{id}", getOne).Methods("GET")
	router.HandleFunc("/inventory", getAll).Methods("GET")
	router.HandleFunc("/inventory/{id}", deleteOne).Methods("DELETE")
	router.HandleFunc("/inventory", create).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func initDB() {
	// Строка подключения к базе данных PostgreSQL
	connStr := "user=" + user + " password=" + password + " host=" + host + " port=" + port + " dbname=" + dbname + " sslmode=disable"
	var err error

	// Открываем соединение с базой данных
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Создаем таблицу Inventory, если она не существует
	createInventoryTable()
}

func createInventoryTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS inventory (
        product_code VARCHAR NOT NULL,
        product_name VARCHAR NOT NULL,
        quantity_in_stock NUMERIC NOT NULL,
        current_price VARCHAR NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}
}

func getOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productCode := params["id"]

	var inventory Inventory

	err := db.QueryRow("SELECT product_code, product_name, quantity_in_stock, current_price FROM inventory WHERE product_code = $1", productCode).
		Scan(&inventory.ProductCode, &inventory.ProductName, &inventory.QuantityInStock, &inventory.CurrentPrice)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	respondWithJSON(w, http.StatusOK, inventory)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT product_code, product_name, quantity_in_stock, current_price FROM inventory")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var inventoryList []Inventory

	for rows.Next() {
		var inventory Inventory
		err := rows.Scan(&inventory.ProductCode, &inventory.ProductName, &inventory.QuantityInStock, &inventory.CurrentPrice)
		if err != nil {
			log.Fatal(err)
		}
		inventoryList = append(inventoryList, inventory)
	}

	respondWithJSON(w, http.StatusOK, inventoryList)
}

func deleteOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productCode := params["id"]

	_, err := db.Exec("DELETE FROM inventory WHERE product_code = $1", productCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func create(w http.ResponseWriter, r *http.Request) {
	var inventory Inventory
	err := json.NewDecoder(r.Body).Decode(&inventory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, err = db.Exec("INSERT INTO inventory (product_code, product_name, quantity_in_stock, current_price) VALUES ($1, $2, $3, $4)",
		inventory.ProductCode, inventory.ProductName, inventory.QuantityInStock, inventory.CurrentPrice)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
