package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func ConnectRedis() *redis.Client {
	// Mengecek apakah klien Redis sudah diinisialisasi sebelumnya
	if redisClient != nil {
		return redisClient
	}

	// Membuat klien Redis jika belum ada
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Alamat dan port Redis
		Password: "",               // Kata sandi Redis (kosong jika tidak ada)
		DB:       0,                // Database Redis (biasanya 0)
	})

	// Menguji koneksi ke Redis
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong) // Output: PONG

	// Menyimpan klien Redis ke variabel global
	redisClient = client

	return client
}

func AppendToAvatarList(avatar Avatar, ID int64) error {
	client := ConnectRedis()
	key := "avatar_list"

	// Mengonversi data JSON baru menjadi string
	avatarJSON, err := json.Marshal(avatar)
	if err != nil {
		return err
	}

	err = client.Append(key, ","+string(avatarJSON)).Err()
	if err != nil {
		return err
	}

	return nil
}
