package config

import (
	"io/ioutil"
	"path"
	"runtime"
	"sync"

	yaml "gopkg.in/yaml.v3"
)

var Configs Config
var once sync.Once

func BuildConfig() {
	once.Do(func() {
		//从外部的conf.yaml文件读取数据
		_, filename, _, _ := runtime.Caller(0)
		filePath := path.Join(path.Dir(filename), "./conf.yaml")
		data, er := ioutil.ReadFile(filePath)
		if er != nil {
			print("yamlFile.Get err")
		}
		//使用yaml包，把读取到的data格式化后解析到config实例中
		err := yaml.Unmarshal(data, &Configs)
		if err != nil {
			panic("decode error")
		}
	})
}
