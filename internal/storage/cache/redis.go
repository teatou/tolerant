package cache

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/teatou/tolerant/internal/storage"
)

type Cache struct {
	client *redis.Client
}

func New(addr, password string, db int) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("cache client connection error: %w", err)
	}

	return &Cache{
		client: client,
	}, nil
}

func (c *Cache) Cache(t storage.Transaction) (int, error) {
	// создать uid

	// создать json

	// сохранить кеш
	// err := c.client.Set(request.ID, requestJSON, 0).Err()
	// if err != nil {
	// 	panic(err)
	// }
	return 0, nil
}

func GetValueFromTransaction(operation, sum, to, from int) ([]byte, error) {
	transaction := storage.Transaction{
		Operation: operation,
		Sum:       sum,
		To:        to,
		From:      from,
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	return transactionJSON, nil
}

func GetTransactionFromValue(value string) (storage.Transaction, error) {
	var transaction storage.Transaction
	err := json.Unmarshal([]byte(value), &transaction)
	if err != nil {
		return storage.Transaction{}, err
	}

	return transaction, nil
}
