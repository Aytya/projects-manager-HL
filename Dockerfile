FROM golang:1.22.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go get -u github.com/spf13/viper
RUN go get -u github.com/joho/godotenv
RUN go get -u github.com/gin-gonic/gin

COPY . .

RUN go build -o projects-manager ./cmd/main.go

CMD ["./projects-manager"]