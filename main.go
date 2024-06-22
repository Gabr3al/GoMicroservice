package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/gabr3al/GoMicroservice/application"
)

func main() {
	config := application.LoadConfig()

	app := application.New(config.RedisAddress)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx, uint(config.ServerPort))
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
