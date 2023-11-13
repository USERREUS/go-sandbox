package mongostore

import (
	"notification-service/internal/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	table = "inventory"
)

// Repository предоставляет методы для взаимодействия с MongoDB.
type Repository struct {
	store *Store
}

// Create добавляет новую запись в коллекцию MongoDB.
func (r *Repository) Create(m *model.Model) error {
	temp := model.Model{
		MsgType:     m.MsgType,
		Date:        time.Now().String(),
		Description: m.Description,
	}

	_, err := r.store.collection.InsertOne(r.store.context, temp)
	if err != nil {
		return err
	}

	return nil
}

// Delete удаляет записи из коллекции MongoDB по указанному типу сообщения.
func (r *Repository) Delete(MsgType string) error {
	filter := bson.M{"message_type": MsgType}
	_, err := r.store.collection.DeleteMany(r.store.context, filter)
	if err != nil {
		return err
	}
	return nil
}

// FindMany находит все записи в коллекции MongoDB по указанному типу сообщения.
func (r *Repository) FindMany(MsgType string) ([]*model.Model, error) {
	var records []*model.Model

	filter := bson.M{"message_type": MsgType}
	cursor, err := r.store.collection.Find(r.store.context, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.store.context)

	var results []model.Model
	if err := cursor.All(r.store.context, &results); err != nil {
		return nil, err
	}

	for _, data := range results {
		records = append(records, &model.Model{
			MsgType:     data.MsgType,
			Description: data.Description,
			Date:        data.Date,
		})
	}

	return records, nil
}

// FindAll находит все записи в коллекции MongoDB.
func (r *Repository) FindAll() ([]*model.Model, error) {
	var records []*model.Model

	filter := bson.M{}
	cursor, err := r.store.collection.Find(r.store.context, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.store.context)

	var results []model.Model
	if err := cursor.All(r.store.context, &results); err != nil {
		return nil, err
	}

	for _, data := range results {
		records = append(records, &model.Model{
			MsgType:     data.MsgType,
			Description: data.Description,
			Date:        data.Date,
		})
	}

	return records, nil
}
