package db

type User struct {
	Model
	Name string
}

//TableName.
func (User) TableName() string {
	return "user"
}
