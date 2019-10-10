package protocol

type NodeInfo struct {
	Types    string `json:"types"`
	CpuNum   string `json:"cpu_num"`
	MemSize  string `json:"mem_size"`
	SourceIP string `json:"source_ip"`
	Name     string `json:"name"`
}
