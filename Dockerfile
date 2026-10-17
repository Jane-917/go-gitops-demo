# ==========================================
# Stage 1: Build Stage (編譯階段)
# ==========================================
FROM golang:1.23-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 先複製 go.mod 與 go.sum，利用 Docker 快取機制優化後續建置速度
COPY go.mod go.sum ./
RUN go mod download

# 複製其餘專案原始碼
COPY . .

# 編譯 Go 程式碼為獨立二進位檔 (關閉 CGO 以確保能在極簡環境執行)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# ==========================================
# Stage 2: Final Stage (運行階段)
# ==========================================
FROM alpine:latest

# 安裝 ca-certificates (確保微服務能進行 HTTPS 外部呼叫) 與 tzdata (時區設定)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 從 Builder 階段複製編譯好的二進位檔
COPY --from=builder /app/main .

# 複製必要的設定檔 (.env)
COPY --from=builder /app/.env .

# 暴露服務使用的埠號
EXPOSE 8080

# 啟動容器時執行的命令
CMD ["./main"]