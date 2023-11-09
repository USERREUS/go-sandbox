package mongostore

import (
	"order-service/internal/app/model"
	"order-service/internal/app/store"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository представляет собой реализацию интерфейса store.Repository для работы с MongoDB.
type Repository struct {
	store *Store
}

// Create создает новый заказ с указанными товарами в MongoDB.
func (r *Repository) Create(items []*model.ProductItem) error {
	order := model.Model{
		Date:         time.Now().String(),
		ProductItems: items,
	}

	_, err := r.store.collection.InsertOne(r.store.context, order)
	if err != nil {
		return err
	}

	return nil
}

// FindOne находит заказ в MongoDB по его идентификатору.
func (r *Repository) FindOne(id string) (*model.Model, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var result model.Model
	err = r.store.collection.FindOne(r.store.context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete удаляет заказ из MongoDB по его идентификатору.
func (r *Repository) Delete(id string) (*model.Model, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var result model.Model
	err = r.store.collection.FindOne(r.store.context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	_, err = r.store.collection.DeleteOne(r.store.context, filter)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// FindAll возвращает все заказы из MongoDB.
func (r *Repository) FindAll() (map[string]*model.Model, error) {
	records := make(map[string]*model.Model)

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
		records[data.OrderCode] = &model.Model{
			Date:         data.Date,
			ProductItems: data.ProductItems,
		}
	}

	if len(records) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return records, nil
}
