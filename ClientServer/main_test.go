package main

import "testing"

func Test_parser(t *testing.T) {
	parser("select * from aaa;")
	t.Error("终止")
}
