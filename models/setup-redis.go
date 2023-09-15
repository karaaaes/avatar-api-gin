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

func UpdateRedisData(updatedAvatar Avatar, ID int64) error {
	updatedAvatar.ID = ID
	client := ConnectRedis()
	key := "avatar_list"

	// Mengambil semua data dari Redis
	existingData, err := client.Get(key).Result()
	if err != nil && err != redis.Nil {
		// Return error jika terjadi kesalahan selain key yang tidak ditemukan
		return err
	}

	// Mengonversi data JSON yang ada menjadi slice dari Avatar
	var avatars []Avatar
	if existingData != "" {
		if err := json.Unmarshal([]byte(existingData), &avatars); err != nil {
			return err
		}
	}

	// Mencari indeks data yang sesuai dengan ID yang diberikan
	var index = -1
	for i, avatar := range avatars {
		if avatar.ID == ID {
			index = i
			break
		}
	}

	// Jika data ditemukan, perbarui data tersebut
	if index != -1 {
		avatars[index] = updatedAvatar
	} else {
		// Jika data tidak ditemukan, tambahkan data baru ke dalam slice
		avatars = append(avatars, updatedAvatar)
	}

	// Mengonversi slice avatar yang telah diperbarui menjadi string JSON
	updatedDataJSON, err := json.Marshal(avatars)
	if err != nil {
		return err
	}

	// Menyimpan data yang telah diperbarui ke Redis
	err = client.Set(key, updatedDataJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func AppendToAvatarList(avatar Avatar, ID int64) error {
	client := ConnectRedis()
	key := "avatar_list"

	// Mengambil semua data dari Redis
	existingData, err := client.Get(key).Result()
	if err != nil && err != redis.Nil {
		// Return error jika terjadi kesalahan selain key yang tidak ditemukan
		return err
	}

	// Mengonversi data JSON yang ada menjadi slice dari avatar
	var avatars []Avatar
	if existingData != "" {
		if err := json.Unmarshal([]byte(existingData), &avatars); err != nil {
			return err
		}
	}

	// Menambahkan avatar baru ke dalam slice
	avatars = append(avatars, avatar)

	// Mengonversi slice avatar yang telah diperbarui menjadi string
	updatedDataJSON, err := json.Marshal(avatars)
	if err != nil {
		return err
	}

	// Menyimpan data yang telah diperbarui ke Redis
	err = client.Set(key, updatedDataJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func DeleteFromAvatarList(ID int64) error {
	client := ConnectRedis()
	key := "avatar_list"

	// Mengambil semua data dari Redis
	existingData, err := client.Get(key).Result()
	if err != nil && err != redis.Nil {
		// Return error jika terjadi kesalahan selain key yang tidak ditemukan
		return err
	}

	// Mengonversi data JSON yang ada menjadi slice dari avatar
	var avatars []Avatar
	if existingData != "" {
		if err := json.Unmarshal([]byte(existingData), &avatars); err != nil {
			return err
		}
	}

	// Mencari indeks data yang sesuai dengan ID yang akan dihapus
	var index = -1
	for i, avatar := range avatars {
		if avatar.ID == ID {
			index = i
			break
		}
	}

	// Jika data ditemukan, hapus data tersebut dari slice
	if index != -1 {
		avatars = append(avatars[:index], avatars[index+1:]...)
	}

	// Mengonversi slice avatar yang telah diperbarui menjadi string
	updatedDataJSON, err := json.Marshal(avatars)
	if err != nil {
		return err
	}

	// Menyimpan data yang telah diperbarui ke Redis
	err = client.Set(key, updatedDataJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
