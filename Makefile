.PHONY: fmt vet lint test test-race test-cover build all clean tools quality godoc

# Variables
BINARY_NAME=server
BUILD_DIR=bin
MAIN_PATH=./cmd/server

# Cible par défaut
all: fmt vet lint test build

# Formatage du code
fmt:
	@echo "Formatage du code..."
	go fmt ./...
	@echo "Formatage termine"

# Vérification statique
vet:
	@echo "Verification statique..."
	go vet ./...
	@echo "Verification statique terminee"

# Linting
lint:
	@echo "Linting..."
	@staticcheck ./... || (echo "staticcheck non installe, installation..." && go install honnef.co/go/tools/cmd/staticcheck@latest && staticcheck ./...)
	@echo "Linting termine"

# Tests
test:
	@echo "Execution des tests..."
	go test ./...
	@echo "Tests termines"

# Tests avec détection de race conditions
test-race:
	@echo "Tests avec detection de race conditions..."
	@go test -race ./... || echo "Tests race echouent (CGO requis sur Windows)"
	@echo "Tests race termines"

# Tests avec couverture
test-cover:
	@echo "Tests avec couverture..."
	go test -cover ./...
	@echo "Tests couverture termines"

# Compilation
build:
	@echo "Compilation..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Compilation terminee"

# Nettoyage
clean:
	@echo "Nettoyage..."
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
	@echo "Nettoyage termine"

# Installation des outils de développement
tools:
	@echo "Installation des outils de developpement..."
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install golang.org/x/tools/cmd/godoc@latest
	@echo "Outils installes"

# Vérification complète de la qualité
quality: fmt vet lint test test-race test-cover
	@echo "Verification de la qualite terminee"

# Lancement de la documentation
godoc:
	@echo "Lancement du serveur de documentation..."
	@echo "Documentation accessible sur: http://localhost:6060"
	@echo "Appuyez sur Ctrl+C pour arreter le serveur"
	@echo ""
	godoc -http=:6060
