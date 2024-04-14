FROM golang:alpine
LABEL maintainer="Art"
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]