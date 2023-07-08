package imports

import (
	"encoding/json"
	"log"
	"net/http"
)
import "gorm.io/gorm"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Model struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BuildErrorResponse Function Response Builder for errors that accepts a messages and a status
func BuildErrorResponse(handlerType string, w http.ResponseWriter, status_code int, err error) {
	log.Printf("ERROR: %s: %s", handlerType, err.Error())
	w.WriteHeader(status_code)

	response := Response{
		Message: err.Error(),
		Status:  "failed",
	}
	json.NewEncoder(w).Encode(response)
}
