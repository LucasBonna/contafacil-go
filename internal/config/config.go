package config

import (
	"log"

	"github.com/sakirsensoy/genv"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
)

type envs struct {
	DB_Host                string
	DB_Port                string
	DB_Name                string
	DB_User                string
	DB_Password            string
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
	RedisAddr              string
	Type                   string
	JWTSecret              string
}

var Env = &envs{
	DB_Host:                genv.Key("DB_HOST").String(),
	DB_Port:                genv.Key("DB_PORT").String(),
	DB_Name:                genv.Key("DB_NAME").String(),
	DB_User:                genv.Key("DB_USER").String(),
	DB_Password:            genv.Key("DB_PASSWORD").String(),
	Port:                   genv.Key("PORT").Default("8000").String(),
	RabbitMQUrl:            genv.Key("RABBITMQ_URL").String(),
	StorageEndpoint:        genv.Key("STORAGE_ENDPOINT").String(),
	StorageRegion:          genv.Key("STORAGE_REGION").String(),
	StorageAccessKeyId:     genv.Key("STORAGE_ACCESS_KEY_ID").String(),
	StorageAccessKeySecret: genv.Key("STORAGE_ACCESS_KEY_SECRET").String(),
	StorageAccountId:       genv.Key("STORAGE_ACCOUNT_ID").String(),
	StorageBucketName:      genv.Key("STORAGE_BUCKET_NAME").Default("CONTAFACIL_DEV").String(),
	TSUsername:             genv.Key("TS_USERNAME").String(),
	TSPassword:             genv.Key("TS_PASSWORD").String(),
	TSBaseUrl:              genv.Key("TS_BASE_URL").String(),
	FrontEndUrl:            genv.Key("FRONTEND_URL").Default("http://localhost:3000").String(),
	RedisAddr:              genv.Key("REDIS_ADDR").String(),
	Type:                   genv.Key("TYPE").String(),
	JWTSecret:              genv.Key("JWTSecret").String(),
}

func InitEnvs() {
	log.Println("Initializing Environment Variables")
}
