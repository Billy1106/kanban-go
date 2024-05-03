FROM golang:1.22-alpine

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .

# air for hot reload
RUN go mod download && go install github.com/cosmtrek/air@v1.51.0

COPY . .

CMD ["air"]