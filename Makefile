
# если в консоли IDE русские символы выводятся не читабельно, то выполнить команду
# $OutputEncoding = [console]::InputEncoding = [console]::OutputEncoding = New-Object System.Text.UTF8Encoding

# Makefile для управления проектом Go + React

.PHONY: help install build start stop clean test lint docker-up docker-down logs


# Показать справку
help:
	@echo Доступные команды:
	@echo   install      - Установить зависимости для frontend и backend
	@echo   build        - Собрать проект
	@echo   start        - Запустить в режиме разработки
	@echo   start-prod   - Запустить в production режиме
	@echo   stop         - Остановить все сервисы
	@echo   clean        - Очистить собранные файлы
	@echo   test         - Запустить тесты
	@echo   docker-up    - Запустить с помощью Docker Compose
	@echo   docker-down  - Остановить Docker контейнеры
	@echo   logs         - Показать логи Docker контейнеров
	@echo   migrations   - Запуск миграции базы данных


migrations-up:
	migrate -path ./internal/repository/migrations -database "postgres://postgres:password@localhost:5432/myapp?sslmode=disable" up

migrations-down:
	migrate -path ./internal/repository/migrations -database "postgres://postgres:password@localhost:5432/myapp?sslmode=disable" down 1

# Установка зависимостей
install:
	@echo Установка зависимостей backend...
	cd backend && go mod download
	@echo Установка зависимостей frontend...
	cd frontend && npm install

# Сборка проекта
build:
	@echo Сборка backend...
	cd backend && go build -o bin/server ./cmd/server
	@echo Сборка frontend...
	cd frontend && npm run build

# Запуск в режиме разработки
start:
	@echo Запуск в режиме разработки...
	@echo Backend будет доступен на http://localhost:8080
	@echo Frontend будет доступен на http://localhost:3000
	@make -j2 start-backend start-frontend

start-backend:
	cd backend && go run ./cmd/server/main.go

start-frontend:
	cd frontend && npm start

# Запуск в production режиме
start-prod: build
	@echo Запуск в production режиме...
	cd backend && ./bin/server &
	cd frontend && npx serve -s build -l 3000

# Остановка всех процессов
stop:
	@echo Остановка всех процессов...
	@pkill -f "go run.*main.go" || true
	@pkill -f "npm start" || true
	@pkill -f "serve.*build" || true

# Очистка
clean:
	@echo Очистка файлов сборки...
	rm -rf backend/bin/
	rm -rf frontend/build/
	rm -rf frontend/node_modules/
	cd backend && go clean

# Тесты
test:
	@echo Запуск тестов backend...
	cd backend && go test -v ./...
	@echo Запуск тестов frontend...
	cd frontend && npm test -- --coverage --watchAll=false


# Docker команды
docker-up:
	@echo Запуск Docker Compose...
	docker-compose up --build -d
	@echo Приложение доступно на:
	@echo   Frontend: http://localhost:3000
	@echo   Backend API: http://localhost:8080
	@echo   PostgreSQL: localhost:5432

docker-down:
	@echo Остановка Docker контейнеров...
	docker-compose down

docker-rebuild:
	@echo Пересборка и запуск Docker контейнеров...
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

logs:
	docker-compose logs -f

# Команды для разработки
dev-setup: install
	@echo Настройка окружения для разработки...
	@echo 1. Убедитесь, что у вас установлены Go (1.21+) и Node.js (18+)
	@echo 2. Запустите 'make start' для запуска в режиме разработки
	@echo 3. Или 'make docker-up' для запуска в Docker

# Проверка зависимостей
check-deps:
	@echo Проверка установленных зависимостей...
	@command -v go >/dev/null 2>&1 || { echo "Go не установлен!"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "Node.js не установлен!"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "Docker не установлен!"; exit 1; }
	@echo Все зависимости установлены!