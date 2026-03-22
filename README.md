
# RUN (en web-app)
go run ./cmd/web/

# Test
go test -v ./...

# Coverage - Ver en terminal
go test -cover ./...

# Coverage - Crear archivo HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Coverage - Crear archivo (sin abrir)
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html


# instalar para mejora en visualizacion de test 
go install gotest.tools/gotestsum@latest

### Uso en Makefile
```makefile
test_handlers:
	cd web-app && gotestsum --format testname -- ./cmd/web/ -run "Test_application_handlers|..."
```

### Formatos disponibles

| Flag | Descripción |
|------|-------------|
| `testname` | Un test por línea con ✓ o ✗ — recomendado |
| `dots` | Minimalista, un punto por test |
| `pkgname` | Agrupa resultados por paquete |
| `standard-verbose` | Similar a `-v` pero con colores |