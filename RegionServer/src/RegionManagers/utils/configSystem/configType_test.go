package config

import "testing"

func Test_test(t *testing.T) {
	BuildConfig()
	t.Error("终止")
}
