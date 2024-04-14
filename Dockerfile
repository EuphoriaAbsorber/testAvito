FROM golang:alpine
LABEL maintainer="Art"
WORKDIR /app
#COPY . .
COPY go.mod go.sum ./
#RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY ./ ./
RUN swag init -g delivery/delivery.go
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]