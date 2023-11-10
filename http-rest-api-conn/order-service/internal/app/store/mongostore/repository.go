package mongostore

import (
	"errors"
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
func (r *Repository) Create(items []*model.ProductItem) (string, error) {
	// Создание объекта заказа на основе переданных товаров и текущей даты.
	order := model.Model{
		Date:         time.Now().String(),
		ProductItems: items,
	}

	// Вставка заказа в коллекцию MongoDB.
	res, err := r.store.collection.InsertOne(r.store.context, order)
	if err != nil {
		return "", err
	}

	// Получение идентификатора вставленного заказа.
	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("Convert Error: InsertedID is not an ObjectID")
	}

	ID := objectID.Hex()

	return ID, nil
}

// FindOne находит заказ в MongoDB по его идентификатору.
func (r *Repository) FindOne(id string) (*model.Model, error) {
	// Преобразование строки идентификатора в ObjectID.
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Формирование фильтра для поиска заказа по идентификатору.
	filter := bson.M{"_id": objectID}

	// Поиск заказа в коллекции MongoDB.
	var result model.Model
	err = r.store.collection.FindOne(r.store.context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete удаляет заказ из MongoDB по его идентификатору.
func (r *Repository) Delete(id string) (*model.Model, error) {
	// Преобразование строки идентификатора в ObjectID.
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Формирование фильтра для поиска заказа по идентификатору.
	filter := bson.M{"_id": objectID}

	// Поиск заказа в коллекции MongoDB.
	var result model.Model
	err = r.store.collection.FindOne(r.store.context, filter).Decode(&result)
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
func (r *Repository) FindAll() (map[string]*model.Model, error) {
	// Инициализация карты для хранения найденных заказов.
	records := make(map[string]*model.Model)

	// Фильтр для получения всех заказов.
	filter := bson.M{}
	cursor, err := r.store.collection.Find(r.store.context, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(r.store.context)

	// Считывание всех заказов из курсора.
	var results []model.Model
	if err := cursor.All(r.store.context, &results); err != nil {
		return nil, err
	}

	// Преобразование результатов в формат, подходящий для возврата клиенту.
	for _, data := range results {
		records[data.OrderCode] = &model.Model{
			Date:         data.Date,
			ProductItems: data.ProductItems,
		}
	}

	// Проверка наличия заказов.
	if len(records) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return records, nil
}
