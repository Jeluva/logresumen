package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	top := flag.Int("top", 5, "cantidad de filas en cada ranking")
	asJSON := flag.Bool("json", false, "salida en JSON")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "uso: logresumen [-top N] [-json] [archivo]   (sin archivo lee stdin)")
		flag.PrintDefaults()
	}
	flag.Parse()

	var in io.Reader = os.Stdin
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		defer f.Close()
		in = f
	}
	s, err := Read(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error leyendo:", err)
		os.Exit(1)
	}
	if *asJSON {
		out := map[string]any{
			"total": s.Total, "invalidas": s.Invalid, "bytes": s.Bytes,
			"estados": s.ByStatusClass, "top_paths": Top(s.Paths, *top),
			"top_ips": Top(s.IPs, *top), "top_errores": Top(s.Errors, *top),
			"por_hora": s.PerHour,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(out)
		return
	}
	printReport(s, *top)
}

func printReport(s *Stats, n int) {
	fmt.Printf("Requests: %d (lineas invalidas: %d)\nTransferido: %.1f KB\n", s.Total, s.Invalid, float64(s.Bytes)/1024)
	classes := make([]string, 0)
	for k := range s.ByStatusClass {
		classes = append(classes, k)
	}
	sort.Strings(classes)
	fmt.Print("\nEstados:")
	for _, c := range classes {
		fmt.Printf("  %s=%d", c, s.ByStatusClass[c])
	}
	fmt.Println()
	section("Paths mas pedidos", Top(s.Paths, n))
	section("IPs mas activas", Top(s.IPs, n))
	section("Paths con mas errores (4xx/5xx)", Top(s.Errors, n))
	hours := Top(s.PerHour, 1)
	if len(hours) > 0 {
		fmt.Printf("\nHora pico: %s con %d requests\n", hours[0].Key, hours[0].N)
	}
}

func section(title string, rows []KV) {
	fmt.Printf("\n%s:\n", title)
	if len(rows) == 0 {
		fmt.Println("  (ninguno)")
	}
	for _, r := range rows {
		fmt.Printf("  %6d  %s\n", r.N, r.Key)
	}
}
