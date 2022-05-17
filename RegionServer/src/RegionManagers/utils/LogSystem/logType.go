package logsystem

/*
*	几种日志的基本类型
 */

// ======== 常规日志 =========
type NormalLog struct {
	LogObj
}

func NewNormalLog(content string) NormalLog {
	res := NormalLog{}
	res.content = content
	return res
}

func (nl *NormalLog) LogGen(ch chan<- LogObj) error {
	ch <- nl.LogObj
	return nil
}

// ======== SQL日志 =========
