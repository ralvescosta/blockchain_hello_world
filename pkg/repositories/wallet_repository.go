package repositories

import "github.com/go-redis/redis/v8"

type WallletRepository struct {
	db *redis.Client
}

func NewWallletRepository(db *redis.Client) *WallletRepository {
	return &WallletRepository{db}
}
