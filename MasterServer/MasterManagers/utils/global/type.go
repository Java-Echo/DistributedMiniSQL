package gloabl

type RegionMeta struct {
	IP   string
	Port string
}

// 记录了master的全局数据结构
var TableMap map[string]RegionMeta
