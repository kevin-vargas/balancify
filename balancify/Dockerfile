FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 go build -o /bin/app

FROM alpine:3.20

COPY --from=builder /bin/app /bin/app

CMD ["/bin/app"]