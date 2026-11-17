package main

import (
	"regexp"
	"strconv"
	"time"
)

// Entry es una linea de log en formato "combined" de nginx/apache.
type Entry struct {
	IP     string
	Time   time.Time
	Method string
	Path   string
	Status int
	Bytes  int64
}

var lineRe = regexp.MustCompile(`^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+)[^"]*" (\d{3}) (\d+|-)`)

// ParseLine devuelve false si la linea no tiene el formato esperado.
func ParseLine(s string) (Entry, bool) {
	m := lineRe.FindStringSubmatch(s)
	if m == nil {
		return Entry{}, false
	}
	t, err := time.Parse("02/Jan/2006:15:04:05 -0700", m[2])
	if err != nil {
		return Entry{}, false
	}
	status, _ := strconv.Atoi(m[5])
	var b int64
	if m[6] != "-" {
		b, _ = strconv.ParseInt(m[6], 10, 64)
	}
	return Entry{IP: m[1], Time: t, Method: m[3], Path: m[4], Status: status, Bytes: b}, true
}
