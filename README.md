# logresumen

CLI chica en Go que lee un log de acceso en formato *combined* (nginx/apache) y imprime un resumen.
Es un **proyecto de práctica**: sin dependencias externas, pensado para aprender parsing y mapas en Go.

## Qué hace
- Cuenta requests, líneas inválidas y bytes transferidos.
- Agrupa por clase de estado (2xx, 3xx, 4xx, 5xx).
- Rankings de paths más pedidos (sin query string), IPs más activas y paths con más errores.
- Hora pico. Salida en texto o en JSON (`-json`).

## Cómo correrlo
Requiere Go 1.23+.

```bash
go run . ejemplo.log            # desde archivo
cat ejemplo.log | go run . -top 3   # desde stdin
go run . -json ejemplo.log      # salida JSON
go test ./...
```

Flags: `-top N` (filas por ranking, default 5), `-json`.

## Limitaciones
- Solo el formato combined; las líneas que no coinciden se cuentan como inválidas.
- Lee todo en memoria en mapas (no está pensado para logs gigantes con millones de paths distintos).
