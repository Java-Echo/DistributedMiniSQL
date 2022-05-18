package logsystem

/*
*	设计日志的基本格式，定义日志相关的接口
 */

// 最底层的日志结构体
type LogObj struct {
	Content string // 日志的内容
	LogType string // 日志的类型 ToDo:现阶段我们使用string，之后可以换成更合适的
	// ToDo:有需要的话加入更多的日志相关信息
}

// 可日志化的对象需要满足这一接口
type LogRealizable interface {
	LogGen(ch chan<- LogObj) error // 将自身转化成一个日志结构体，并传入日志channel当中
}

// ========对日志对象的不同处理方法========
// 接受日志结构体，将其写入本地的相关文件
func LogWrite(obj LogObj) error {
	//ToDo:测试阶段，暂时将其打印出来，尚未设计“写入的具体位置”
	//ToDo:需要在这里添加日志的“格式信息”，也就是一些括号啊，横线啊，让它看起来像日志
	println("[LogWrite测试]" + obj.Content)
	return nil
}
