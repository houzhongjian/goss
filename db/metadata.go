package db

//Metadata 元数据.
type Metadata struct {
	Model
	Name      string
	Type      string
	Size      int64
	Hash      string
	StoreNode string
}

//TableName .
func (Metadata) TableName() string {
	return "metadata"
}

//Create 创建.
func (m *Metadata) Create() error {
	return Db.Create(&m).Error
}
