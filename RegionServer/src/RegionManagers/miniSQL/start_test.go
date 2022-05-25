package main

import "testing"

func Test_start(t *testing.T) {
	start()
	t.Error("终止")
}
