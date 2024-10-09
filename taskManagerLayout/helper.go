package taskManagerLayout

import (
	"log"
	"taskManager/app/db"
	"time"
)

func GetAllProjectNames() []string {
	var projectNames []string
	result := db.Get().Model(&Project{}).Pluck("name", &projectNames)
	if result.Error != nil {
		log.Printf("Error fetching project names: %v", result.Error)
		return nil
	}
	return projectNames
}

type Project struct {
	ProjectID   uint      `gorm:"primaryKey;column:project_id" json:"project_id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	ProjectType string    `gorm:"column:project_type;not null" json:"project_type"`
	Description string    `gorm:"not null" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	CreatedBy   uint      `gorm:"column:created_by" json:"created_by"`
}
