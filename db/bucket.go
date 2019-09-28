package db

type Bucket struct {
	Model
	Name   string
	UserID int
}

func (Bucket) TableName() string {
	return "bucket"
}
