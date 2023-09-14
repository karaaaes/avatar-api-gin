package avatarController

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"avatar-api-gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var avatars []models.Avatar
	models.DB.Find(&avatars)
	c.JSON(http.StatusOK, gin.H{"avatars": avatars})
}

func Show(c *gin.Context) {
	var avatars models.Avatar
	id := c.Param("id")

	if err := models.DB.First(&avatars, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"avatar": avatars})
}

func Create(c *gin.Context) {
	var avatar models.Avatar

	if err := c.ShouldBindJSON(&avatar); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Encrypt nama images
	encryptedImageName := generateUniqueCode(avatar.AvatarImage)
	avatar.AvatarImage = encryptedImageName

	// Memeriksa apakah avatar_username sudah ada dalam database
	if err := models.DB.Where("avatar_username = ?", avatar.AvatarUsername).First(&models.Avatar{}).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Sudah terdapat avatar dengan username tersebut"})
		return
	}

	models.DB.Create(&avatar)

	// Menambahkan data JSON baru ke dalam "avatar_list" di Redis
	if err := models.AppendToAvatarList(avatar, avatar.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal menambahkan data ke Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar": avatar})
}

func generateUniqueCode(input string) string {
	// Pisahkan input menjadi dua bagian: "images/" dan "abc"
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return "" // Input tidak sesuai format yang diharapkan
	}

	// Ambil string setelah "images/"
	valueToHashBeforeJpg := parts[1]
	partsJpg := strings.Split(valueToHashBeforeJpg, ".")

	// Ambil string setelah ".jpg"
	valueToHash := partsJpg[0]

	// Buat objek hash MD5
	hasher := md5.New()

	// Konversi string valueToHash menjadi byte array dan hash
	hasher.Write([]byte(valueToHash))

	// Dapatkan hasil hash dalam bentuk byte
	hashedBytes := hasher.Sum(nil)

	// Konversi byte hash menjadi string heksadesimal
	hashedString := hex.EncodeToString(hashedBytes)

	finalString := parts[0] + "/" + hashedString + "." + partsJpg[1]

	return finalString
}

func Update(c *gin.Context) {
	var avatar models.Avatar
	id := c.Param("id")

	if err := c.ShouldBindJSON(&avatar); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if models.DB.Model(&avatar).Where("id = ?", id).Updates(&avatar).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat update avatar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil di Update"})
}

func Delete(c *gin.Context) {
	var avatar models.Avatar
	idParam := c.Param("id")

	// Konversi parameter id menjadi int64
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	// Cari avatar dengan ID yang sesuai
	result := models.DB.Where("id = ?", id).First(&avatar)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Avatar tidak ditemukan"})
		return
	}

	// Hapus avatar dari database
	models.DB.Delete(&avatar)
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

func Random(c *gin.Context) {
	var avatar models.Avatar
	// Mengambil semua ID dari database
	var avatarIDs []int64
	if err := models.DB.Model(&avatar).Pluck("id", &avatarIDs).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Mengambil ID secara acak
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(avatarIDs))
	randomID := avatarIDs[randomIndex]

	// Mengambil data Avatar berdasarkan ID yang diambil secara acak
	if err := models.DB.First(&avatar, randomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"avatar": avatar})
}
