package sqlstore

import (
	"database/sql"
	"fmt"
	"inventory/internal/app/model"
	"inventory/internal/app/store"

	"github.com/google/uuid"
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
	m.ID = uuid.New().String()
	return r.store.db.QueryRow(
		fmt.Sprintf("INSERT INTO %s (id, name, count, cost) VALUES ($1, $2, $3, $4) RETURNING id", table),
		m.ID,
		m.Name,
		m.Count,
		m.Cost,
	).Scan(&m.ID)
}

// Update обновляет существующую запись в базе данных.
func (r *Repository) Update(m *model.Model) error {
	m.ID = uuid.New().String()
	return r.store.db.QueryRow(
		fmt.Sprintf("UPDATE %s SET name = $2, count = $3, cost = $4 WHERE id = $1 RETURNING id", table),
		m.ID,
		m.Name,
		m.Count,
		m.Cost,
	).Scan(&m.ID)
}

// Delete удаляет запись из базы данных по идентификатору.
func (r *Repository) Delete(id string) (*model.Model, error) {
	p := &model.Model{}
	if err := r.store.db.QueryRow(
		fmt.Sprintf("SELECT id, name, count, cost FROM %s WHERE id = $1", table),
		id,
	).Scan(
		&p.ID,
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
		fmt.Sprintf("DELETE FROM %s WHERE id = $1", table),
		id,
	); err != nil {
		return nil, err
	}

	return p, nil
}

// FindOne находит запись в базе данных по идентификатору.
func (r *Repository) FindOne(id string) (*model.Model, error) {
	p := &model.Model{}
	if err := r.store.db.QueryRow(
		fmt.Sprintf("SELECT id, name, count, cost FROM %s WHERE id = $1", table),
		id,
	).Scan(
		&p.ID,
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

	rows, err := r.store.db.Query("SELECT id, name, count, cost FROM inventory")
	if err != nil {
		return nil, store.ErrRecordNotFound
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&m.ID, &m.Name, &m.Count, &m.Cost); err != nil {
			return nil, store.ErrRecordNotFound
		}

		records[m.ID] = &model.Model{
			ID:    m.ID,
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
