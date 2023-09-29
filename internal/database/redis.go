package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nu12/audio-gonverter/internal/user"
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

func (r *RedisRepo) Save(u *user.User) error {
	userJson, err := json.Marshal(*u)
	if err != nil {
		return err
	}

	err = r.Client.Set(context.TODO(), u.UUID, string(userJson), 1*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo) Load(uuid string) (*user.User, error) {

	userJson, err := r.Client.Get(context.TODO(), uuid).Result()
	if err != nil {
		u := user.NewUser()
		return &u, err
	}

	var u user.User
	err = json.Unmarshal([]byte(userJson), &u)
	if err != nil {
		u := user.NewUser()
		return &u, err
	}
	return &u, nil
}

func (r *RedisRepo) Exist(uuid string) (bool, error) {
	_, err := r.Client.Get(context.TODO(), uuid).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
