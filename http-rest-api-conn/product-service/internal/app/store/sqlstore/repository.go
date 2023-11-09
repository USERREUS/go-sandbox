package sqlstore

import (
	"database/sql"
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"

	"github.com/google/uuid"
)

type ProductRepository struct {
	store *Store
}

func (r *ProductRepository) Create(p *model.Product) error {
	p.ID = uuid.New().String()
	return r.store.db.QueryRow(
		"INSERT INTO products (id, name, weight, description) VALUES ($1, $2, $3, $4) RETURNING id",
		p.ID,
		p.Name,
		p.Weight,
		p.Description,
	).Scan(&p.ID)
}

func (r *ProductRepository) FindOne(id string) (*model.Product, error) {
	p := &model.Product{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, weight, description FROM products WHERE id = $1",
		id,
	).Scan(
		&p.ID,
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

func (r *ProductRepository) Delete(id string) (*model.Product, error) {
	p := &model.Product{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, weight, description FROM products WHERE id = $1",
		id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Weight,
		&p.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	_, err := r.store.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *ProductRepository) FindAll() (map[string]*model.Product, error) {
	records := make(map[string]*model.Product)
	p := &model.Product{}

	rows, err := r.store.db.Query("SELECT id, name, weight, description FROM products")
	if err != nil {
		return nil, store.ErrRecordNotFound
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&p.ID, &p.Name, &p.Weight, &p.Description); err != nil {
			return nil, store.ErrRecordNotFound
		}

		records[p.ID] = p
	}

	if err := rows.Err(); err != nil {
		return nil, store.ErrRecordNotFound
	}

	if len(records) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return records, nil
}
