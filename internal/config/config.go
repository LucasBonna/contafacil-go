package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Db_url                 string
	Port                   string
	RabbitMQUrl            string
	StorageEndpoint        string
	StorageRegion          string
	StorageAccessKeyId     string
	StorageAccessKeySecret string
	StorageAccountId       string
	StorageBucketName      string
	TSUsername             string
	TSPassword             string
	TSBaseUrl              string
	FrontEndUrl            string
)

func LoadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	Db_url = os.Getenv("DB_URL")
	Port = os.Getenv("PORT")
	RabbitMQUrl = os.Getenv("RABBITMQ_URL")
	StorageEndpoint = os.Getenv("STORAGE_ENDPOINT")
	StorageRegion = os.Getenv("STORAGE_REGION")
	StorageAccessKeyId = os.Getenv("STORAGE_ACCESS_KEY_ID")
	StorageAccessKeySecret = os.Getenv("STORAGE_ACCESS_KEY_SECRET")
	StorageAccountId = os.Getenv("STORAGE_ACCOUNT_ID")
	StorageBucketName = os.Getenv("STORAGE_BUCKET_NAME")
	TSUsername = os.Getenv("TS_USERNAME")
	TSPassword = os.Getenv("TS_PASSWORD")
	TSBaseUrl = os.Getenv("TS_BASE_URL")
	FrontEndUrl = os.Getenv("FRONTEND_URL")
}
