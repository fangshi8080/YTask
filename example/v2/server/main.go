package main

import (
	"context"
	workers2 "github.com/gojuukaze/YTask/example/v2/server/workers"
	"github.com/gojuukaze/YTask/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// For the server, you do not need to set up the poolSize
	// Server端无需设置poolSize，
	broker := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	backend := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(60*5),
		ytask.Config.ResultExpires(60*5),
	)

	ser.Add("group1", "add", workers2.Add)
	ser.Add("group1", "add_sub", workers2.AddSub)
	ser.Add("group1", "retry", workers2.Retry)
	ser.Add("group1", "add_user", workers2.AppendUser)

	ser.Run("group1", 3)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ser.Shutdown(context.Background())

}
