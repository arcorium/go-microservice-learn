package db

import (
	"Microservices/event/model"
	"errors"
	"log"
)

type DBType string

const (
	// Database constant
	DB_NAME       string = "microservice"
	DB_MONGODB    DBType = "MONGODB"
	DB_MYSQL      DBType = "MYSQL"
	DB_POSTGRESQL DBType = "POSTGRESQL"

	// Collection constant
	COLL_EVENT   string = "event"
	COLL_USER    string = "user"
	COLL_BOOKING string = "booking"
)

type DatabaseService interface {
	AddEvent(*model.Event) (any, error)
	FindEventById(any) (*model.Event, error)
	FindEventByName(string) (*model.Event, error)
	FindAllEvents() ([]model.Event, error)
	FindAllAvailableEvents() ([]model.Event, error)
}

func NewDatabaseService(dbType_ DBType, dbConnection_ string) (DatabaseService, error) {
	log.Println("Using database service:", dbType_)

	switch dbType_ {
	case DB_MONGODB:
		return NewMongoDBService(dbConnection_)
	default:
		log.Println("Database service", dbType_, "is not supported!")
	}

	return nil, errors.New("failed to create database service")
}
