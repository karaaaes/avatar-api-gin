package models

type Avatar struct {
	ID             int64  `gorm:"primaryKey" json:"id"`
	AvatarName     string `gorm:"type:varchar(255)" json:"avatar_name"`
	AvatarImage    string `gorm:"type:varchar(255)" json:"avatar_image"`
	AvatarUsername string `gorm:"type:varchar(255)" json:"avatar_username"`
	AvatarPassword string `gorm:"type:varchar(255)" json:"avatar_password"`
	AvatarEmail    string `gorm:"type:varchar(255)" json:"avatar_email"`
}
