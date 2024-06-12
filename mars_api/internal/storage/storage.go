package storage

import (
	"github.com/google/mars_api/internal/storage/cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Storage struct {
	*cache.Cache
	db *gorm.DB
}

func NewStorage(redisConn *redis.Client, db *gorm.DB) *Storage {
	return &Storage{Cache: cache.NewCache(redisConn, db), db: db}
}
