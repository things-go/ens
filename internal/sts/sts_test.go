package sts

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/plugin/soft_delete"
)

// Announce 公告-面向所有人的消息
type Announce struct {
	Id        int64                 `gorm:"column:id;not null;autoIncrement:true;primaryKey" json:"id,omitempty"`
	Title     string                `gorm:"column:title;type:varchar(255);not null;comment:标题" json:"title,omitempty"`                        // 标题
	Content   string                `gorm:"column:content;type:varchar(2048);not null;comment:内容" json:"content,omitempty"`                   // 内容
	Priority  uint                  `gorm:"column:priority;type:int(10) unsigned;not null;default:255;comment:优先级" json:"priority,omitempty"` // 优先级
	Visible   bool                  `gorm:"column:visible;type:tinyint(1) unsigned;not null;default:0;comment:是否显示" json:"visible,omitempty"` // 是否显示
	CreatedAt time.Time             `gorm:"column:created_at;type:datetime;not null;comment:发布时间" json:"created_at,omitempty"`                // 发布时间
	UpdatedAt time.Time             `gorm:"column:updated_at;type:datetime;not null" json:"updated_at,omitempty"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;type:bigint(20);not null;default:0" json:"deleted_at,omitempty"`
}

func Test(t *testing.T) {
	v := Announce{}
	fmt.Println(Parse(v))
}
