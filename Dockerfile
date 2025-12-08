# 第一階段：構建
FROM golang:1.24-alpine AS builder

# ... (安裝必要工具和 migrate 工具等保持不變) ...

# 設定工作目錄
WORKDIR /app

# 1. 【優化快取】只複製模組相關文件
# 這些文件（go.mod, go.sum）變動頻率低，可最大化利用快取
COPY go.mod go.sum ./

# 2. 下載依賴
# 只有當 go.mod/go.sum 變動時，才會重新執行此步驟
RUN go mod download

# 3. 複製所有其他專案原始碼（包含 docs/ 等內部套件）
# 確保所有內部套件在 go build 之前都存在於工作目錄中
COPY . .

# 4. 編譯應用程式 (現在所有依賴和內部套件都已到位)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# 第二階段：運行
FROM alpine:latest

# 安裝 ca-certificates（HTTPS 需要）
RUN apk --no-cache add ca-certificates tzdata

# 設定時區
ENV TZ=Asia/Taipei

# 建立非 root 使用者
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 設定工作目錄
WORKDIR /app

# 從 builder 階段複製編譯好的執行檔
COPY --from=builder /app/main .

# 複製 migrate 工具
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# 複製配置檔案
COPY --from=builder /app/config ./config

# 複製資料庫相關檔案
COPY --from=builder /app/db ./db

# 複製初始化腳本
COPY --from=builder /app/scripts ./scripts

# 改變檔案擁有者（除了 migrate）
RUN chown -R appuser:appuser /app && \
    chmod +x /usr/local/bin/migrate && \
    chmod +x /app/scripts/init-db.sh

# 切換到非 root 使用者
USER appuser

# 暴露端口
EXPOSE 8080

# 執行應用程式
CMD ["./main"]