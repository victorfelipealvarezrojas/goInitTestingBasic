
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