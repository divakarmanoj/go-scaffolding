package imports

import "time"
import "gorm.io/gorm"

type Response struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ID        uint   `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Model struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
