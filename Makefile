.PHONY: help build run docker-build docker-up docker-down docker-logs migrate-up migrate-down

help: ## 顯示幫助訊息
	@echo "可用的指令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 編譯應用程式
	go build -o bin/main ./cmd/api

run: ## 執行應用程式
	go run cmd/api/main.go

docker-build: ## 建立 Docker 映像
	docker-compose build

docker-up: ## 啟動 Docker 容器
	docker-compose up -d

docker-down: ## 停止 Docker 容器
	docker-compose down

docker-logs: ## 查看 Docker 日誌
	docker-compose logs -f

docker-restart: ## 重啟 Docker 容器
	docker-compose restart

migrate-up: ## 執行資料庫遷移
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/golang_framework?sslmode=disable" -verbose up

migrate-down: ## 回滾資料庫遷移
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/golang_framework?sslmode=disable" -verbose down

test: ## 執行測試
	go test -v ./...

clean: ## 清理編譯產物
	rm -rf bin/