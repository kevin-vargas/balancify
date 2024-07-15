FROM node:18 as build
WORKDIR /app
COPY app/package*.json .
RUN npm install
# TODO: remove dockerignore and improve copy
COPY app .
RUN npm run build

FROM golang:1.20-alpine AS builder

USER root 

WORKDIR /app

COPY --from=build /app/dist /app/app/dist

COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 go build -o /bin/app

FROM alpine:3.20

COPY --from=builder /bin/app /bin/app

CMD ["/bin/app"]
