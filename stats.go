package main

import (
	"bufio"
	"io"
	"sort"
	"strings"
)

type Stats struct {
	Total, Invalid int
	Bytes          int64
	ByStatusClass  map[string]int // "2xx", "4xx", ...
	Paths          map[string]int
	IPs            map[string]int
	Errors         map[string]int // paths con 4xx/5xx
	PerHour        map[string]int // "2006-01-02 15h"
}

func NewStats() *Stats {
	return &Stats{
		ByStatusClass: map[string]int{}, Paths: map[string]int{},
		IPs: map[string]int{}, Errors: map[string]int{}, PerHour: map[string]int{},
	}
}

func (s *Stats) Add(e Entry) {
	s.Total++
	s.Bytes += e.Bytes
	s.ByStatusClass[string(rune('0'+e.Status/100))+"xx"]++
	// Sin query string, para agrupar /a?x=1 con /a.
	p := e.Path
	if i := strings.IndexByte(p, '?'); i >= 0 {
		p = p[:i]
	}
	s.Paths[p]++
	s.IPs[e.IP]++
	if e.Status >= 400 {
		s.Errors[p]++
	}
	s.PerHour[e.Time.Format("2006-01-02 15")+"h"]++
}

// Read procesa todas las lineas de r.
func Read(r io.Reader) (*Stats, error) {
	s := NewStats()
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if e, ok := ParseLine(line); ok {
			s.Add(e)
		} else {
			s.Invalid++
		}
	}
	return s, sc.Err()
}

type KV struct {
	Key string
	N   int
}

// Top devuelve las n claves mas frecuentes (desempate alfabetico).
func Top(m map[string]int, n int) []KV {
	out := make([]KV, 0, len(m))
	for k, v := range m {
		out = append(out, KV{k, v})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].N != out[j].N {
			return out[i].N > out[j].N
		}
		return out[i].Key < out[j].Key
	})
	if n > 0 && len(out) > n {
		out = out[:n]
	}
	return out
}
