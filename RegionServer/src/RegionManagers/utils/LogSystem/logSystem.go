package logsystem

// 通用日志的构造函数
func LogStart() chan<- LogObj {
	input := make(chan LogObj)
	// ToDo:这里需要开启一个协程用以处理input中的日志，在这个协程中将会调用以下的不同处理方法
	go func() {
		for {
			obj := <-input
			// ToDo:通用日志仅仅用来写日志内容，需要判断日志的类型，对于不同类型的日志，需要写在不同的文件里，同时也需要有自己的自定义格式
			println("得到了一个新的日志,它的内容为:" + obj.content)
		}
	}()
	return input
}
