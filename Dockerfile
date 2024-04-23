# Используем официальный образ Golang в качестве базы
# Я удалил сертификаты с проекта, чтобы не брали :)
FROM golang:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /go/src/app

# Копируем все файлы вашего проекта в текущую директорию контейнера
COPY . .

COPY git.apps.okd.sebbia.org.crt /usr/local/share/ca-certificates/git.apps.okd.sebbia.org.crt
COPY sebbia.orgIntermediateCA.crt /usr/local/share/ca-certificates/sebbia.orgIntermediateCA.crt
COPY sebbia.orgRootCA.crt /usr/local/share/ca-certificates/sebbia.orgRootCA.crt
RUN update-ca-certificates

# Устанавливаем указанные зависимости проекта
RUN go get -d -v \
    github.com/go-resty/resty/v2@v2.11.0 \
    golang.org/x/net@v0.17.0 \
    gopkg.in/yaml.v2@v2.4.0

# Собираем ваше приложение
RUN go build -o main .

# Запускаем приложение при старте контейнера
CMD ["./main"]