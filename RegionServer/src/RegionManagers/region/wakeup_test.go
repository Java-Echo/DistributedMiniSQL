package regionWorker

import (
	"fmt"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"testing"
)

func TestMain(m *testing.M) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	fmt.Println("初始化完成")
	m.Run()
}

func Test_findLocalTable(t *testing.T) {
	// fmt.Println(filepath.Abs(config.Configs.Minisql_table_store))
	files := findLocalTable(config.Configs.Minisql_table_store)
	for _, fileName := range files {
		fmt.Println(fileName)
	}
	t.Error("终止")
}
