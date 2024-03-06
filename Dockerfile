FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /mbserver ./server/server.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder mbserver ./

EXPOSE 8090

CMD ["./mbserver"]