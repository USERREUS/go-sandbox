package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost" // Адрес хоста PostgreSQL
	port     = "5433"      // Порт PostgreSQL
	user     = "postgres"  // Имя пользователя PostgreSQL
	password = "password"  // Пароль пользователя
	dbname   = "postgres"  // Имя вашей базы данных
	sslmode  = "disable"   // Отключение SSL (важно для контейнера без SSL)
)

type Product struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
}

var db *sql.DB

func main() {
	// Инициализация базы данных
	initDB()

	// Создание маршрутизатора gorilla/mux
	router := mux.NewRouter()

	// Определение обработчиков для эндпоинтов
	router.HandleFunc("/products/{id}", getOne).Methods("GET")
	router.HandleFunc("/products", getAll).Methods("GET")
	router.HandleFunc("/products/{id}", deleteOne).Methods("DELETE")
	router.HandleFunc("/products", create).Methods("POST")

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
}

// @Summary Получить продукт по ID
// @Description Получить информацию о продукте по его коду
// @ID get-product-by-id
// @Produce json
// @Param id path string true "Код продукта"
// @Success 200 {object} Product
// @Router /products/{id} [get]
func getOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productCode := params["id"]

	// Запрос к базе данных для получения информации о продукте
	var product Product
	err := db.QueryRow("SELECT code, name, weight, description FROM product WHERE code = $1", productCode).Scan(&product.Code, &product.Name, &product.Weight, &product.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, product)
}

// @Summary Получить все продукты
// @Description Получить список всех продуктов
// @ID get-all-products
// @Produce json
// @Success 200 {array} Product
// @Router /products [get]
func getAll(w http.ResponseWriter, r *http.Request) {
	// Запрос к базе данных для получения всех продуктов
	rows, err := db.Query("SELECT code, name, weight, description FROM product")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	var products []Product

	// Итерация по результатам запроса и добавление продуктов в слайс
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Code, &product.Name, &product.Weight, &product.Description); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		products = append(products, product)
	}

	respondWithJSON(w, http.StatusOK, products)
}

// @Summary Удалить продукт по ID
// @Description Удалить продукт по его коду
// @ID delete-product
// @Param id path string true "Код продукта"
// @Produce json
// @Success 200 {string} string
// @Router /products/{id} [delete]
func deleteOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productCode := params["id"]

	// Удаление продукта из базы данных
	_, err := db.Exec("DELETE FROM product WHERE code = $1", productCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted"})
}

// @Summary Создать новый продукт
// @Description Создать новый продукт
// @ID create-product
// @Accept json
// @Param product body Product true "Данные нового продукта"
// @Produce json
// @Success 201 {object} Product
// @Router /products [post]
func create(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newProduct); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Генерация случайного кода продукта
	productCode := uuid.New().String()

	// Вставка нового продукта в базу данных
	sqlStatement := `
        INSERT INTO product (code, name, weight, description)
        VALUES ($1, $2, $3, $4)
    `
	_, err := db.Exec(sqlStatement, productCode, newProduct.Name, newProduct.Weight, newProduct.Description)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	respondWithJSON(w, http.StatusCreated, newProduct)
}

//Don't work
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}

//Don't work
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
