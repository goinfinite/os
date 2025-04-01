package dbModel

type ScheduledTaskTag struct {
	ID              uint64 `gorm:"primaryKey"`
	Tag             string `gorm:"not null"`
	ScheduledTaskID uint64 `gorm:"not null"`
}

func (ScheduledTaskTag) TableName() string {
	return "scheduled_tasks_tags"
}
