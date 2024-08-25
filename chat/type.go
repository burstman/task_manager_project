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
	ProjectID   uint      `gorm:"primaryKey;column:project_id"`
	Name        string    `gorm:"unique;not null"`
	TypeProject string    `gorm:"column:type_project;not null"`
	Description string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	CreatedBy   uint      `gorm:"column:created_by"`
}
