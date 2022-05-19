package regionWorker

import (
	"fmt"
	"testing"
)

func Test_findMasterCopy(t *testing.T) {
	files := findSlaveCopy("...")
	for _, fileName := range files {
		fmt.Println(fileName)
	}
	t.Error("终止")
}
