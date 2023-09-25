package imports

import "gorm.io/gorm"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Cursor  uint        `json:"cursor,omitempty"`
}

type Model struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
