package entities

import (
	"time"

	"gorm.io/gorm"
)

type ScheduleFile struct {
	BaseModel
	BuildingNo string `gorm:"column:building_no;not null" json:"buildingNo"`
	FileStream []byte `gorm:"column:file_stream;type:varbinary(max)" json:"fileStream"`
}

func (m *ScheduleFile) BeforeCreate(*gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *ScheduleFile) BeforeUpdate(*gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
