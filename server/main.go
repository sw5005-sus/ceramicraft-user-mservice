package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/sw5005-sus/ceramicraft-user-mservice/common/utils"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/grpc"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/mq"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository"
)

var (
	sigCh = make(chan os.Signal, 1)
)

func main() {
	fmt.Println("Starting ceramicraft-user-mservice...")
	config.Init()
	log.InitLogger()
	log.Logger.Info("Logger initialized.")
	utils.InitJwtSecret()
	log.Logger.Info("JWT secret initialized.")
	repository.Init()
	log.Logger.Info("Database initialized.")
	mq.InitKafka()
	log.Logger.Info("Kafka initialized.")
	proxy.InitZitadel()
	log.Logger.Info("Zitadel proxy initialized.")
	go grpc.Init(sigCh)
	go http.Init(sigCh)
	// listen terminage signal
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh // Block until signal is received
	debug.PrintStack()
	log.Logger.Infof("Received signal: %v, shutting down", sig)
}
