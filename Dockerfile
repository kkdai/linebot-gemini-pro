# 使用官方 Golang 鏡像作為基礎鏡像
FROM golang:1.21.5

# 在容器內設置工作目錄
WORKDIR /app

# 將本地的文件複製到容器的工作目錄
COPY . .

# 下載所有依賴項
RUN go mod download

# 將本地的源碼文件複製到容器的工作目錄
# COPY . .

# 編譯應用程式
RUN go build -o main .

# 指定容器啟動時執行的命令
CMD ["/app/main"]
