package protocol

//GossProtocol 传输协议号.
type GossProtocol int

const (
	WriteFileProrocol = 10000 //写文件.
	ReadFileProrocol  = 20000 //读文件.
	NodeAddProtocol   = 30000 //新增节点.
	NodelDelProtocol  = 40000 //删除节点.
)
