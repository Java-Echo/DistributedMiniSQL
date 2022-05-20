package regionWorker

// 查询语句
func SQL_select(SQL string) (string, bool) {
	return "", true
}

// 对部分表项的修改
func SQL_change(SQL string) bool {
	return true
}

// 对数据表的增加或者删除
func SQL_tableChange(SQL string) bool {
	return true
}
