package regionWorker

/*--------------master节点的操作--------------*/
// ToDo:查询语句
func MasterSQLSelect(SQL string) (string, bool) {
	return "", true
}

// ToDo:对部分表项的修改
func MasterSQLChange(SQL string) bool {
	return true
}

// ToDo:对数据表的增加或者删除
func MasterSQLTableChange(SQL string) bool {
	return true
}

/*--------------常规的SQL操作--------------*/
