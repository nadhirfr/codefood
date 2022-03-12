package helpers

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var REDIS *redis.Client

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName:   os.Getenv("MYSQL_DBNAME"),
	}
	return &dbConfig
}

func DbURL(dbConfig *DBConfig) mysql.Config {
	return mysql.Config{
		DSN: fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.DBName,
		), // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}
}

func RedisInit() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	network := os.Getenv("REDIS_NETWORK")
	if len(network) == 0 {
		network = "tcp"
	}

	password := os.Getenv("REDIS_PASSWORD")
	if len(password) == 0 {
		password = ""
	}

	REDIS = redis.NewClient(&redis.Options{
		Network:  network,
		Addr:     dsn, //redis port
		Password: password,
	})
	_, err := REDIS.Ping().Result()
	if err != nil {
		panic(err)
	}
}
