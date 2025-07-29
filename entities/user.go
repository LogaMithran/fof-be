package entities

import "time"

type User struct {
	Id      uint    `json:"id,omitempty" gorm:"column:id"`
	Name    string  `json:"name,omitempty" gorm:"column:name;index:idx_name,unique"`
	Email   string  `json:"email,omitempty" gorm:"column:email;index:idx_email,unique"`
	Friends []*User `json:"-" gorm:"many2many:user_friends"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "user"
}
