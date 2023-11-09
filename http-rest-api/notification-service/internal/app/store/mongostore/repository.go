package mongostore

import (
	"notification-service/internal/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	table = "inventory"
)

type Notification struct {
	MsgType     string `bson:"message_type"`
	Description string `bson:"description"`
	Data        string `bson:"data,omitempty"`
}

type Repository struct {
	store *Store
}

func (r *Repository) Create(m *model.Model) error {
	temp := Notification{
		MsgType:     m.MsgType,
		Data:        time.Now().String(),
		Description: m.Description,
	}

	_, err := r.store.collection.InsertOne(r.store.context, temp)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(MsgType string) error {
	filter := bson.M{"message_type": MsgType}
	_, err := r.store.collection.DeleteMany(r.store.context, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindMany(MsgType string) ([]*model.Model, error) {
	var records []*model.Model

	filter := bson.M{"message_type": MsgType}
	cursor, err := r.store.collection.Find(r.store.context, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.store.context)

	var results []Notification
	if err := cursor.All(r.store.context, &results); err != nil {
		return nil, err
	}

	for _, data := range results {
		records = append(records, &model.Model{
			MsgType:     data.MsgType,
			Description: data.Description,
			Data:        data.Data,
		})
	}

	return records, nil
}

func (r *Repository) FindAll() ([]*model.Model, error) {
	var records []*model.Model

	filter := bson.M{}
	cursor, err := r.store.collection.Find(r.store.context, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.store.context)

	var results []Notification
	if err := cursor.All(r.store.context, &results); err != nil {
		return nil, err
	}

	for _, data := range results {
		records = append(records, &model.Model{
			MsgType:     data.MsgType,
			Description: data.Description,
			Data:        data.Data,
		})
	}

	return records, nil
}
