package main

import (
	"github.com/joho/godotenv"
	"github.com/razagr/pensionera/config"
	"github.com/razagr/pensionera/repository"
	"github.com/razagr/pensionera/service"
)

func main() {

	//Load env variables
	godotenv.Load()

	// Configure the environment variables
	window, symbols, finnHubAPIKey := config.NewConfig().Configuration()

	// get our storage repository to store the data in CSV format
	storage := repository.NewFileStorage()

	// create currency data provider
	repo := repository.NewFinnHubRepository(symbols, finnHubAPIKey)

	// inject dependencies and Run our service
	service.NewService(window, symbols, storage, repo).Run()
}
