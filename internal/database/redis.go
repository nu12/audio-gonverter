package database

import (
	"context"
	"encoding/json"

	"github.com/nu12/audio-gonverter/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func NewRedis(host, port, password string) *RedisRepo {
	return &RedisRepo{
		Client: redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: password,
			DB:       0,
		}),
	}
}

func (r *RedisRepo) Save(u *model.User) error {
	userJson, err := json.Marshal(*u)
	if err != nil {
		return err
	}

	err = r.Client.Set(context.TODO(), u.UUID, string(userJson), 0).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *RedisRepo) Load(uuid string) *model.User {

	userJson, err := r.Client.Get(context.TODO(), uuid).Result()
	if err != nil {
		u := model.NewUser()
		return &u
	}

	var user model.User
	err = json.Unmarshal([]byte(userJson), &user)
	return &user
}
