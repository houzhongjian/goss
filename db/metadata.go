package db

//Metadata 元数据.
type Metadata struct {
	Model
	Name      string   `gorm:"name"`
	Type      string   `gorm:"type"`
	Size      int      `gorm:"size"`
	Hash      string   `gorm:"hash"`
	Node      []string `gorm:"-"`
	StoreNode string   `gorm:"store_node"`
}

//TableName .
func (Metadata) TableName() string {
	return "metadata"
}

func (m *Metadata) Create() error {
	m.Name = "1.jpeg"
	m.Type = "image/png"
	m.Size = 123234
	m.Hash = "sdfasdasdfsdf"
	m.StoreNode = "127.0.0.1:9001,127.0.0.1:9002"
	return Db.Create(m).Error
}
