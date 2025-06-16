package properties

import (
	"os"
	"strconv"
	"time"
)

const (
	ApplicationRepositoryName = "go-api"
	CacheDbSecretKey          = "_DB_SECRET"
	CacheExpiration           = time.Hour * 24 * 16
)

func GetNotesQueueURL() string {
	return os.Getenv("NOTES_QUEUE_URL")
}

func GetDatabaseSslMode() string {
	return os.Getenv("DATABASE_SSL_MODE")
}

func GetEnv() string {
	return os.Getenv("ENV")
}

func GetRegion() string {
	return os.Getenv("AWS_REGION")
}

func GetDatabaseTimeout() string {
	return os.Getenv("DATABASE_TIMEOUT")
}

func GetDatabaseMaxOpenConnections() int {
	v, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_OPEN_CONNECTIONS"))
	return v
}
func GetDatabaseMaxIdleConnections() int {
	v, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_IDLE_CONNECTIONS"))
	return v
}

func GetSecretDatabase() string {
	return os.Getenv("SECRET_DATABASE")
}

func GetCacheSetHost() string {
	return os.Getenv("CACHE_SET_HOST")
}

func GetCacheGetHost() string {
	return os.Getenv("CACHE_GET_HOST")
}

func GetCachePrefix() string {
	return os.Getenv("CACHE_PREFIX")
}

func GetCacheMaxRetries() int {
	v, _ := strconv.Atoi(os.Getenv("CACHE_MAX_RETRIES"))
	return v
}

func GetCacheTimeout() time.Duration {
	v, _ := strconv.Atoi(os.Getenv("CACHE_TIMEOUT"))
	return time.Second * time.Duration(v)
}
