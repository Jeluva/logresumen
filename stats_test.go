package main

import (
	"strings"
	"testing"
)

const sample = `1.1.1.1 - - [10/Oct/2025:13:55:36 +0000] "GET /a?x=1 HTTP/1.1" 200 100 "-" "curl"
1.1.1.1 - - [10/Oct/2025:13:56:00 +0000] "GET /a HTTP/1.1" 404 - "-" "curl"
2.2.2.2 - - [10/Oct/2025:14:01:00 +0000] "POST /b HTTP/1.1" 500 50 "-" "curl"
basura
`

func TestRead(t *testing.T) {
	s, err := Read(strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	if s.Total != 3 || s.Invalid != 1 || s.Bytes != 150 {
		t.Fatalf("totales mal: %+v", s)
	}
	if s.Paths["/a"] != 2 || s.Errors["/a"] != 1 || s.ByStatusClass["5xx"] != 1 {
		t.Fatalf("conteos mal: %+v", s)
	}
	if top := Top(s.IPs, 1); top[0].Key != "1.1.1.1" {
		t.Fatalf("top ip mal: %v", top)
	}
}
