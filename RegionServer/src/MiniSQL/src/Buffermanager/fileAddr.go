package Buffermanager

type FileAddr struct {
	filePageID uint32 // 文件页编号
	offset     uint32 // 页内偏移量
}

/*
* 比较两个“文件指针”的大小：等于返回0，大于返回1，小于返回-1
 */
func (fa *FileAddr) compare(fa_ *FileAddr) int {
	if fa.filePageID == fa_.filePageID && fa.offset == fa_.offset {
		return 0
	} else if fa.filePageID < fa_.filePageID || fa.offset < fa_.offset {
		return -1
	} else if fa.filePageID > fa_.filePageID || fa.offset > fa_.offset {
		return 1
	}
	// 这里是不可能到的吧
	return 233
}
