FROM golang:latest
WORKDIR /app
EXPOSE 80
EXPOSE 443
EXPOSE 8080

COPY . .
CMD ["go", "run", "main.go"]