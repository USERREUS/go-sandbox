package sqlstore

import (
	"database/sql"
	"fmt"
	"inventory/internal/app/model"
	"inventory/internal/app/store"
)

const (
	table = "inventory"
)

// Repository представляет собой реализацию интерфейса store.Repository для работы с базой данных SQL.
type Repository struct {
	store *Store
}

// Create добавляет новую запись в базу данных.
func (r *Repository) Create(m *model.Model) error {
	return r.store.db.QueryRow(
		fmt.Sprintf("INSERT INTO %s (code, name, count, cost) VALUES ($1, $2, $3, $4) RETURNING code", table),
		m.Code,
		m.Name,
		m.Count,
		m.Cost,
	).Scan(&m.Code)
}

// Update обновляет существующую запись в базе данных.
func (r *Repository) Update(m *model.Model) error {
	return r.store.db.QueryRow(
		fmt.Sprintf("UPDATE %s SET name = $2, count = $3, cost = $4 WHERE code = $1 RETURNING code", table),
		m.Code,
		m.Name,
		m.Count,
		m.Cost,
	).Scan(&m.Code)
}

// Delete удаляет запись из базы данных по идентификатору.
func (r *Repository) Delete(code string) (*model.Model, error) {
	p := &model.Model{}
	if err := r.store.db.QueryRow(
		fmt.Sprintf("SELECT code, name, count, cost FROM %s WHERE code = $1", table),
		code,
	).Scan(
		&p.Code,
		&p.Name,
		&p.Count,
		&p.Cost,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	if _, err := r.store.db.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE code = $1", table),
		code,
	); err != nil {
		return nil, err
	}

	return p, nil
}

// FindOne находит запись в базе данных по идентификатору.
func (r *Repository) FindOne(code string) (*model.Model, error) {
	p := &model.Model{}
	if err := r.store.db.QueryRow(
		fmt.Sprintf("SELECT code, name, count, cost FROM %s WHERE code = $1", table),
		code,
	).Scan(
		&p.Code,
		&p.Name,
		&p.Count,
		&p.Cost,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return p, nil
}

// FindAll возвращает все записи из базы данных.
func (r *Repository) FindAll() (map[string]*model.Model, error) {
	records := make(map[string]*model.Model)
	m := &model.Model{}

	rows, err := r.store.db.Query("SELECT code, name, count, cost FROM inventory")
	if err != nil {
		return nil, store.ErrRecordNotFound
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&m.Code, &m.Name, &m.Count, &m.Cost); err != nil {
			return nil, store.ErrRecordNotFound
		}

		records[m.Code] = &model.Model{
			Code:  m.Code,
			Name:  m.Name,
			Count: m.Count,
			Cost:  m.Cost,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, store.ErrRecordNotFound
	}

	if len(records) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return records, nil
}
