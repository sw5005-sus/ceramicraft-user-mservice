package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	ceramicraftsecure "github.com/sw5005-sus/ceramicraft-secure"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/grpc"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/mq"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/redis"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/telemetry"
)

var (
	sigCh = make(chan os.Signal, 1)
)

func main() {
	fmt.Println("Starting ceramicraft-user-mservice...")
	config.Init()
	log.InitLogger()
	log.Logger.Info("Logger initialized.")
	ceramicraftsecure.Init()
	log.Logger.Infof("ceramicraft-secure init done.")
	repository.Init()
	log.Logger.Info("Database initialized.")
	redis.InitRedis()
	log.Logger.Info("Redis initialized.")
	mq.InitKafka()
	log.Logger.Info("Kafka initialized.")
	proxy.InitZitadel()
	log.Logger.Info("Zitadel proxy initialized.")
	shutdownTrace := telemetry.InitTracer()
	log.Logger.Info("Tracing initialized.")
	defer shutdownTrace()
	shutdownMetrics := telemetry.InitMetrics()
	log.Logger.Info("Metrics initialized.")
	defer shutdownMetrics()
	go grpc.Init(sigCh)
	go http.Init(sigCh)
	// listen terminage signal
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh // Block until signal is received
	debug.PrintStack()
	log.Logger.Infof("Received signal: %v, shutting down", sig)
}
