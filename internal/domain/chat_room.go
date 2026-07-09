package domain

type ChatRoom struct {
	ID    uint `gorm:"primaryKey"`
	Title string
	Type  ChatRoomType `gorm:"type:enum('private', 'group');default:'private'"`
}

type ChatRoomType string

const (
	TypePrivate ChatRoomType = "private"
	TypeGroup   ChatRoomType = "group"
)
