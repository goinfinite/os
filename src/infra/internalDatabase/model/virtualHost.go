package dbModel

type VirtualHost struct {
	ID             uint   `gorm:"primarykey"`
	Hostname       string `gorm:"not null"`
	Type           string `gorm:"not null"`
	RootDirectory  string `gorm:"not null"`
	ParentHostname *string
	Mappings       []Mapping
}
