package db

//Disk 硬盘.
type Disk struct {
	Model
	IP       string //所属ip.
	Key      string //硬盘标识.
	DiskPath string //硬盘存储路径.
	Size     int    //硬盘大小单位G.
}

//TableName .
func (Disk) TableName() string {
	return "disk"
}

//Write 写入硬盘.
func (Disk) Write(b []byte) error {
	return nil
}
