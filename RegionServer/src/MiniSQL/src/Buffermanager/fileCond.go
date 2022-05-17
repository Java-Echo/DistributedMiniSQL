package Buffermanager

type FileCond struct {
	DelFirst  FileAddr // 第一条被删除的记录地址(链表向)
	DelLast   FileAddr //最后一条被删除的记录地址(链表向)
	NewInsert FileAddr // 文件末尾可插入新数据的地址
	TotalPage int      // 当前文件中总共的页数
}
