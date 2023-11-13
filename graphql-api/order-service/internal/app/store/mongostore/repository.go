package mongostore

import (
	"order-service/internal/app/model"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Repository представляет собой реализацию интерфейса store.Repository для работы с MongoDB.
type Repository struct {
	store *Store
}

// Create создает новый заказ с указанными товарами в MongoDB.
func (r *Repository) Create(items []*model.ProductItem) (*model.Model, error) {
	// Создание объекта заказа на основе переданных товаров и текущей даты.
	order := model.Model{
		OrderCode:    uuid.New().String(),
		Date:         time.Now().String(),
		ProductItems: items,
	}

	// Вставка заказа в коллекцию MongoDB.
	_, err := r.store.collection.InsertOne(r.store.context, order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// FindOne находит заказ в MongoDB по его идентификатору.
func (r *Repository) FindOne(code string) (*model.Model, error) {
	// Формирование фильтра для поиска заказа по идентификатору.
	filter := bson.M{"order_code": code}

	// Поиск заказа в коллекции MongoDB.
	var result model.Model
	err := r.store.collection.FindOne(r.store.context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete удаляет заказ из MongoDB по его идентификатору.
func (r *Repository) Delete(code string) (*model.Model, error) {
	// Формирование фильтра для поиска заказа по идентификатору.
	filter := bson.M{"order_code": code}

	// Поиск заказа в коллекции MongoDB.
	var result model.Model
	err := r.store.collection.FindOne(r.store.context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Удаление заказа из коллекции MongoDB.
	_, err = r.store.collection.DeleteOne(r.store.context, filter)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// FindAll возвращает все заказы из MongoDB.
func (r *Repository) FindAll() ([]*model.Model, error) {
	// Фильтр для получения всех заказов.
	filter := bson.M{}
	cursor, err := r.store.collection.Find(r.store.context, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(r.store.context)

	// Считывание всех заказов из курсора.
	var results []*model.Model
	if err := cursor.All(r.store.context, &results); err != nil {
		return nil, err
	}

	return results, nil
}
