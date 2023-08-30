FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags='-w -s' -o /crm ./cmd/api

CMD ["/crm"]
