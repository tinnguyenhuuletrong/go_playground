package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"ttin.com/play2022/network"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	go network.StartSimpleHttpServer(ctx, ":8080")

	// wait for interupt
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	fmt.Println("Stoping...")
	cancel()

	time.Sleep(2 * time.Second)
}
