# Етап побудови
# Етап збірки
FROM golang:1.22.6-alpine3.19 AS builder

# Копіюємо код в контейнер
COPY . /telegram-bot
WORKDIR /telegram-bot

# Завантажуємо модулі
RUN go mod download

# Створення директорії для зібраних файлів і збірка
RUN mkdir -p ./.bin && go build -o ./.bin/bot ./cmd/bot/main.go

# Етап фінального образу
FROM alpine:latest

WORKDIR /root/

# Копіюємо зібраний файл та конфігурації
COPY --from=builder /telegram-bot/.bin/bot .
COPY --from=builder /telegram-bot/configs configs/

EXPOSE 80

CMD ["./bot"]
