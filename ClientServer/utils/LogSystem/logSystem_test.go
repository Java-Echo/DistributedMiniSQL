package logsystem

import (
	"testing"
)

func TestLogStart(t *testing.T) {
	input := LogStart()
	log1 := NewNormalLog("测试日志1")
	log2 := NewNormalLog("测试日志2")
	log1.LogGen(input)
	log2.LogGen(input)
	t.Error("终止测试")
}
