package cache

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Cache struct {
	conn *redis.Client
	db   *gorm.DB
}

func NewCache(conn *redis.Client, db *gorm.DB) *Cache {
	cache := &Cache{conn: conn, db: db}

	// 初始化数据
	cache.initI18()

	return cache
}
