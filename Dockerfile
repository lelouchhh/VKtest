# Используем базовый образ для Go
FROM golang:latest

# Создадим директорию
RUN mkdir /VKtest

# Скопируем всё в директорию
ADD . /VKtest/

# Установим рабочей папкой директорию
WORKDIR /VKtest

# Получим зависимости, которые использовали в боте
RUN go get github.com/botanio/sdk/go

ENV TEST1=value1
ENV host=185.200.241.2
ENV port=5432
ENV user=vk
ENV password=vkdb
ENV db=vkdb
ENV ssl=disable
ENV key=5995894659:AAG82B6kbmD17TmmPcKT5Zzqz4S6LpQolYQ
# Соберём приложение
RUN go build -o main .

# Запустим приложение
CMD ["/VKtest/main"]