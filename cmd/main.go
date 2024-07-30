package main

import (
	"github.com/Aytya/projects-manager-HL"
	"github.com/Aytya/projects-manager-HL/internal/handler"
	"github.com/Aytya/projects-manager-HL/internal/repository"
	"github.com/Aytya/projects-manager-HL/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

// @title Projects-Manager
// @ version 1.22.4
// @description API Server for Projects-Manager Application

// @host localhost:8080
func main() {
	if err := initConfig(); err != nil {
		log.Fatal("Error at initializing config", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error at loading .env file", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Database: viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})

	if err != nil {
		log.Fatal("Failed at initializing db", err)
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	server := new(projects_manager.Server)
	if err := server.Run(viper.GetString("8080"), handlers.InitRoutes()); err != nil {
		log.Fatal("Error occured while running http server")
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
