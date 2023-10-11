package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teatou/tolerant/internal/config"
	"github.com/teatou/tolerant/internal/handlers/add"
	"github.com/teatou/tolerant/internal/storage/cache"
	"github.com/teatou/tolerant/internal/storage/relational"
	"github.com/teatou/tolerant/pkg/mylogger"
)

const configEnv = "CONFIG"

func main() {
	val, ok := os.LookupEnv(configEnv)
	if !ok {
		val = "configs/dev.yaml"
	}

	cfg, err := config.LoadConfig(val)
	if err != nil {
		panic("uploading config error")
	}

	logger, err := mylogger.NewZapLogger(cfg.Logger.Level)
	if err != nil {
		panic("making mylogger error")
	}
	defer logger.Sync()

	// init storage
	storage, err := relational.New(cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DbName)
	if err != nil {
		panic(err)
	}

	// init cache
	redis, err := cache.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	// router
	router := gin.Default()

	router.POST("/add", add.New(storage, redis, logger))
	router.POST("/transfer")
	router.POST("/withdraw")

	// server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		// Получаем все ключи из Redis
		keys, err := redis.Keys("*").Result()
		if err != nil {
			panic(err)
		}

		// Проходим по каждому ключу и получаем значение
		for _, key := range keys {
			value, err := redis.Get(key).Result()
			if err != nil {
				panic(err)
			}
			fmt.Printf("Ключ: %s, Значение: %s\n", key, value)
		}
		log.Println("starting server...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
