package protocol

//GossProtocol 传输协议号.
type GossProtocol int

const (
	WriteFileProrocol = 10000 //写文件.
	ReadFileProrocol  = 20000 //读文件.
	NodeAddProtocol   = 30000 //新增节点.

	ConnAuthProtocol       = 1000 //连接授权.
	ReportNodeInfoProtocol = 1001 //上报节点信息.
	MsgProtocol            = 1002 //发送消息.
	AddNodeProtocol        = 1003 //新增节点.
	RemoveNodeProtocol     = 1004 //删除节点.
)
