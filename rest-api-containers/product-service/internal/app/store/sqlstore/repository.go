package sqlstore

import (
	"database/sql"
	"product-service/internal/app/model"
	"product-service/internal/app/store"

	"github.com/google/uuid"
)

// ProductRepository представляет собой репозиторий продуктов для взаимодействия с базой данных.
type ProductRepository struct {
	store *Store
}

// Create создает новую запись о продукте в базе данных.
func (r *ProductRepository) Create(p *model.Product) error {
	p.Code = uuid.New().String()
	return r.store.db.QueryRow(
		"INSERT INTO products (code, name, weight, description) VALUES ($1, $2, $3, $4) RETURNING code",
		p.Code,
		p.Name,
		p.Weight,
		p.Description,
	).Scan(&p.Code)
}

// FindOne возвращает продукт из базы данных по его уникальному коду.
func (r *ProductRepository) FindOne(code string) (*model.Product, error) {
	p := &model.Product{}
	if err := r.store.db.QueryRow(
		"SELECT code, name, weight, description FROM products WHERE code = $1",
		code,
	).Scan(
		&p.Code,
		&p.Name,
		&p.Weight,
		&p.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return p, nil
}

// Delete удаляет продукт из базы данных по его уникальному коду.
func (r *ProductRepository) Delete(code string) (*model.Product, error) {
	p := &model.Product{}
	if err := r.store.db.QueryRow(
		"SELECT code, name, weight, description FROM products WHERE code = $1",
		code,
	).Scan(
		&p.Code,
		&p.Name,
		&p.Weight,
		&p.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	_, err := r.store.db.Exec("DELETE FROM products WHERE code = $1", code)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// FindAll возвращает все продукты из базы данных.
func (r *ProductRepository) FindAll() ([]*model.Product, error) {
	var records []*model.Product
	p := &model.Product{}

	rows, err := r.store.db.Query("SELECT code, name, weight, description FROM products")
	if err != nil {
		return nil, store.ErrRecordNotFound
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&p.Code, &p.Name, &p.Weight, &p.Description); err != nil {
			return nil, store.ErrRecordNotFound
		}

		records = append(records, &model.Product{
			Code:        p.Code,
			Name:        p.Name,
			Weight:      p.Weight,
			Description: p.Description,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, store.ErrRecordNotFound
	}

	if len(records) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return records, nil
}
