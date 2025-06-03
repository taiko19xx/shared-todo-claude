FROM golang:1.23-alpine AS development

WORKDIR /app

# CGOを有効にするために必要なツールをインストール
RUN apk add --no-cache gcc musl-dev sqlite-dev

# CGOを有効にする環境変数を設定
ENV CGO_ENABLED=1

RUN go install github.com/air-verse/air@v1.61.7

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]