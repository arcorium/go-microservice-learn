package config

import (
	"Microservices/event/persistence/db"
	"encoding/json"
	"log"
	"os"
)

type ServiceConfig struct {
	Https              bool        `json:"https"`
	DatabaseType       db.DBType   `json:"databaseType,omitempty"`
	DatabaseConnection string      `json:"databaseConnection,omitempty"`
	Endpoint           string      `json:"endpoint,omitempty"`
	HttpsConfig        HttpsConfig `json:"httpsConfig,omitempty"`
}

type HttpsConfig struct {
	CertificatePath string `json:"certificatePath"`
	KeyPath         string `json:"keyPath"`
	Endpoint        string `json:"endpoint"`
}

func LoadConfiguration(filepath_ string) (ServiceConfig, error) {

	config := ServiceConfig{
		Https: false,
	}

	file, err := os.Open(filepath_)
	if err != nil {
		log.Println(err)
		return config, err
	}

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Println(err)
		return config, err
	}

	// Check empty configuration
	if len(config.DatabaseConnection) == 0 {
		config.DatabaseConnection = os.Getenv("MONGODB_URI")
	}
	if len(config.DatabaseType) == 0 {
		config.DatabaseType = db.DB_MONGODB
	}
	if len(config.Endpoint) == 0 {
		config.Endpoint = os.Getenv("IP_ADDRESS") + ":" + os.Getenv("PORT")
	}

	return config, err
}