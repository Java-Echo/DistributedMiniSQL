package rpc

import (
	"log"
	"testing"
)

func TestReportServiceClient_ReportTable(t *testing.T) {
	client, err := DialReportService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	request := []LocalTable{
		{"aab", "127.0.0.1"},
		{"bbc", "127.0.0.1"},
		{"cca", "127.0.0.1"},
	}
	err = client.ReportTable(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(reply)
	t.Error("终止")
}
