package models

import "time"

type Category struct {
	Id       uint
	Image    *string
	Name     string
	CreateAt time.Time
	UpdateAt time.Time
}
