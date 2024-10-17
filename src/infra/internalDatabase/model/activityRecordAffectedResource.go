package dbModel

type ActivityRecordAffectedResource struct {
	ID                       uint64 `gorm:"primarykey"`
	SystemResourceIdentifier string `gorm:"not null"`
	ActivityRecordID         uint64 `gorm:"not null"`
}

func (ActivityRecordAffectedResource) TableName() string {
	return "activity_records_affected_resources"
}
