package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type Request struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// при старте сервера проверяем есть ли запросы в кэше, если есть - выполняем
// сервер бежит, выполняет запросы
// когда сервер падает, замыкаем контекст, роллбекаем транзакции из постгрес, записываем в кэш
// восстанавливаем сервер

func main() {
	// Создание клиента Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Проверка соединения с Redis
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	// Получение запроса
	request := Request{
		ID:   "1",
		Data: "example data",
	}

	// Сериализация запроса в JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	// Сохранение запроса в Redis
	err = client.Set(request.ID, requestJSON, 0).Err()
	if err != nil {
		panic(err)
	}

	// Проверка наличия сохраненного запроса
	val, err := client.Get(request.ID).Result()
	if err == redis.Nil {
		fmt.Println("Запрос не найден")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("Найден сохраненный запрос:", val)

		// Восстановление состояния и обработка запроса

		// Удаление запроса из Redis после успешной обработки
		err = client.Del(request.ID).Err()
		if err != nil {
			panic(err)
		}
	}
}
