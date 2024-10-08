package chat

import (
	"time"
)

// type Chatter interface {
// 	SendMessage(message map[string]any) error
// 	ReceiveMessage() (map[string]any, error)
// }

// type rasaChatBot struct {
// 	serverUrl string
// }

// type Handler struct {
// 	bot Chatter
// }

type Project struct {
	ProjectID   uint      `gorm:"primaryKey;column:project_id" json:"project_id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	ProjectType string    `gorm:"column:project_type;not null" json:"project_type"`
	Description string    `gorm:"not null" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	CreatedBy   uint      `gorm:"column:created_by" json:"created_by"`
}
