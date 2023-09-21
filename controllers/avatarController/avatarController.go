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

	// Mendapatkan protokol dari permintaan pengguna
	protocol := "http" // Default to http
	if c.Request.TLS != nil {
		protocol = "https"
	}

	// Mendapatkan base URL dari permintaan pengguna
	baseURL := protocol + "://" + c.Request.Host
	url := c.Request.URL.String()
	finalUrl := baseURL + "" + url

	// Mengambil status HTTP dinamis
	status := c.Writer.Status()

	// Membuat respons JSON yang sesuai dengan format yang diinginkan
	response := gin.H{
		"response": gin.H{
			"code":   status,
			"status": http.StatusText(status),
			"url":    finalUrl,
		},
		"data": avatars,
	}
	c.JSON(http.StatusOK, response)
}

func Show(c *gin.Context) {
	var avatars models.Avatar
	id := c.Param("id")

	// Mendapatkan protokol dari permintaan pengguna
	protocol := "http" // Default to http
	if c.Request.TLS != nil {
		protocol = "https"
	}

	// Mendapatkan base URL dari permintaan pengguna
	baseURL := protocol + "://" + c.Request.Host
	url := c.Request.URL.String()
	finalUrl := baseURL + "" + url

	if err := models.DB.First(&avatars, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Membuat respons JSON yang sesuai dengan format yang diinginkan
			response := gin.H{
				"response": gin.H{
					"code":   http.StatusNotFound,
					"status": http.StatusText(http.StatusNotFound),
					"url":    finalUrl,
				},
				"data": gin.H{
					"message": "Data tidak ditemukan",
				},
			}

			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return
		} else {
			// Membuat respons JSON yang sesuai dengan format yang diinginkan
			response := gin.H{
				"response": gin.H{
					"code":   http.StatusInternalServerError,
					"status": http.StatusText(http.StatusInternalServerError),
					"url":    finalUrl,
				},
				"data": gin.H{
					"message": err.Error(),
				},
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
	}

	// Mengambil status HTTP dinamis
	status := c.Writer.Status()

	// Membuat respons JSON yang sesuai dengan format yang diinginkan
	response := gin.H{
		"response": gin.H{
			"code":   status,
			"status": http.StatusText(status),
			"url":    finalUrl,
		},
		"data": avatars,
	}

	c.JSON(http.StatusOK, response)
}

func Create(c *gin.Context) {
	var avatar models.Avatar

	// Mendapatkan protokol dari permintaan pengguna
	protocol := "http" // Default to http
	if c.Request.TLS != nil {
		protocol = "https"
	}

	// Mendapatkan base URL dari permintaan pengguna
	baseURL := protocol + "://" + c.Request.Host
	url := c.Request.URL.String()
	finalUrl := baseURL + "" + url

	if err := c.ShouldBindJSON(&avatar); err != nil {
		response := gin.H{
			"response": gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"url":    finalUrl,
			},
			"data": gin.H{
				"message": err.Error(),
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Encrypt nama images
	encryptedImageName := generateUniqueCode(avatar.AvatarImage)
	avatar.AvatarImage = encryptedImageName

	// Memeriksa apakah avatar_username sudah ada dalam database
	if err := models.DB.Where("avatar_username = ?", avatar.AvatarUsername).First(&models.Avatar{}).Error; err == nil {
		response := gin.H{
			"response": gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"url":    finalUrl,
			},
			"data": gin.H{
				"message": "Sudah terdapat avatar dengan username tersebut",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	models.DB.Create(&avatar)

	// Menambahkan data JSON baru ke dalam "avatar_list" di Redis
	if err := models.AppendToAvatarList(avatar, avatar.ID); err != nil {
		response := gin.H{
			"response": gin.H{
				"code":   http.StatusInternalServerError,
				"status": http.StatusText(http.StatusInternalServerError),
				"url":    finalUrl,
			},
			"data": gin.H{
				"message": "Gagal menambahkan data ke Redis",
			},
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	// Membuat respons JSON yang sesuai dengan format yang diinginkan
	response := gin.H{
		"response": gin.H{
			"code":   http.StatusOK,
			"status": http.StatusText(http.StatusOK),
			"url":    finalUrl,
		},
		"data": avatar,
	}
	c.JSON(http.StatusOK, response)
}

func Update(c *gin.Context) {
	var avatar models.Avatar
	id := c.Param("id")

	// Mendapatkan protokol dari permintaan pengguna
	protocol := "http" // Default to http
	if c.Request.TLS != nil {
		protocol = "https"
	}

	// Mendapatkan base URL dari permintaan pengguna
	baseURL := protocol + "://" + c.Request.Host
	url := c.Request.URL.String()
	finalUrl := baseURL + "" + url

	if err := c.ShouldBindJSON(&avatar); err != nil {
		response := gin.H{
			"response": gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"url":    finalUrl,
			},
			"data": gin.H{
				"message": err.Error(),
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if models.DB.Model(&avatar).Where("id = ?", id).Updates(&avatar).RowsAffected == 0 {
		response := gin.H{
			"response": gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"url":    finalUrl,
			},
			"data": gin.H{
				"message": "Id tidak ditemukan. Tidak dapat update avatar",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	idUpdate, _ := strconv.ParseInt(id, 10, 64)

	if err := models.UpdateRedisData(avatar, idUpdate); err != nil {
		avatar.ID = idUpdate
		response := gin.H{
			"response": gin.H{
				"code":   http.StatusInternalServerError,
				"status": http.StatusText(http.StatusInternalServerError),
				"url":    finalUrl,
			},
			"data": gin.H{
				"message": "Gagal update data ke Redis",
			},
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	avatar.ID = idUpdate

	var avatarUpdate models.Avatar
	if err := models.DB.First(&avatarUpdate, idUpdate).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Membuat respons JSON yang sesuai dengan format yang diinginkan
			response := gin.H{
				"response": gin.H{
					"code":   http.StatusNotFound,
					"status": http.StatusText(http.StatusNotFound),
					"url":    finalUrl,
				},
				"data": gin.H{
					"message": "Data tidak ditemukan",
				},
			}

			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return
		} else {
			// Membuat respons JSON yang sesuai dengan format yang diinginkan
			response := gin.H{
				"response": gin.H{
					"code":   http.StatusInternalServerError,
					"status": http.StatusText(http.StatusInternalServerError),
					"url":    finalUrl,
				},
				"data": gin.H{
					"message": err.Error(),
				},
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
	}

	response := gin.H{
		"response": gin.H{
			"code":   http.StatusOK,
			"status": http.StatusText(http.StatusOK),
			"url":    finalUrl,
		},
		"message": "Data berhasil di update",
		"data":    avatarUpdate,
	}
	c.JSON(http.StatusOK, response)
}

func Delete(c *gin.Context) {
	var avatar models.Avatar
	idParam := c.Param("id")

	// Mendapatkan protokol dari permintaan pengguna
	protocol := "http" // Default to http
	if c.Request.TLS != nil {
		protocol = "https"
	}

	// Mendapatkan base URL dari permintaan pengguna
	baseURL := protocol + "://" + c.Request.Host
	url := c.Request.URL.String()
	finalUrl := baseURL + "" + url

	// Konversi parameter id menjadi int64
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		// Membuat respons JSON yang sesuai dengan format yang diinginkan
		response := gin.H{
			"response": gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"url":    finalUrl,
			},
			"message": err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Cari avatar dengan ID yang sesuai
	result := models.DB.Where("id = ?", id).First(&avatar)
	if result.Error != nil {
		// Membuat respons JSON yang sesuai dengan format yang diinginkan
		response := gin.H{
			"response": gin.H{
				"code":   http.StatusNotFound,
				"status": http.StatusText(http.StatusNotFound),
				"url":    finalUrl,
			},
			"message": "Id tidak ditemukan",
		}
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	// Hapus dari Redis
	idUpdate, _ := strconv.ParseInt(idParam, 10, 64)

	if err := models.DeleteFromAvatarList(idUpdate); err != nil {
		avatar.ID = idUpdate
		response := gin.H{
			"response": gin.H{
				"code":   http.StatusInternalServerError,
				"status": http.StatusText(http.StatusInternalServerError),
				"url":    finalUrl,
			},
			"data": gin.H{
				"message": "Data di Redis gagal dihapus.",
			},
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	// Hapus avatar dari database
	models.DB.Delete(&avatar)

	// Mengambil status HTTP dinamis
	status := c.Writer.Status()

	// Membuat respons JSON yang sesuai dengan format yang diinginkan
	response := gin.H{
		"response": gin.H{
			"code":   status,
			"status": http.StatusText(status),
			"url":    finalUrl,
		},
		"message": "Data berhasil dihapus",
	}
	c.JSON(http.StatusOK, response)
}

func Random(c *gin.Context) {
	var avatar models.Avatar

	// Mengambil status HTTP dinamis
	status := c.Writer.Status()

	// Mendapatkan protokol dari permintaan pengguna
	protocol := "http" // Default to http
	if c.Request.TLS != nil {
		protocol = "https"
	}

	// Mendapatkan base URL dari permintaan pengguna
	baseURL := protocol + "://" + c.Request.Host
	url := c.Request.URL.String()
	finalUrl := baseURL + "" + url

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

	// Membuat respons JSON yang sesuai dengan format yang diinginkan
	response := gin.H{
		"response": gin.H{
			"code":   status,
			"status": http.StatusText(status),
			"url":    finalUrl,
		},
		"data": avatar,
	}

	c.JSON(http.StatusOK, response)
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
