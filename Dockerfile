# Builder
FROM golang:1.15.0-alpine3.12 as builder

RUN apk update && apk upgrade && \
    apk --update add git make curl

WORKDIR /app

COPY . .

RUN make engine

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata curl && \
    mkdir /app 

WORKDIR /app 

EXPOSE 8000

COPY --from=builder /app/engine /app

CMD /app/engine