package main

import (
	"log"
	"os"

	"github.com/NikiTesla/lamoda_test/pkg/environment"
	"github.com/NikiTesla/lamoda_test/pkg/jsonrpc"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("can't load env variables, err:", err.Error())
	}

	configFile := os.Getenv("CONFIGFILE")
	env, err := environment.NewEnvironment(configFile)
	if err != nil {
		log.Fatal("can't load environment, err:", err.Error())
	}

	server := jsonrpc.NewServer()

	if err = server.Run(env.Config.Port); err != nil {
		log.Fatal("error occured while running server, error:", err.Error())
	}
}
