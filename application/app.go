package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New(redisAddress string) *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr: redisAddress,
		}),
	}
	fmt.Print("\n")
	fmt.Println("Redis connection established")

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context, port uint) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis connection: ", err)
		}
	}()

	fmt.Println("Listening on port", port)

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
