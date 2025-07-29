package entities

type UserFriend struct {
	UserId   uint64 `gorm:"column:user_id"`
	FriendId uint64 `gorm:"column:friend_id"`
}

func (UserFriend) TableName() string {
	return "user_friends"
}
