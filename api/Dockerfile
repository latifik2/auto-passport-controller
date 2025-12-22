FROM golang:1.25.4-alpine as builder

WORKDIR /app

COPY src/ .

RUN go build -o auto-passport-api .


FROM alpine:3.22.2

WORKDIR /app
RUN adduser -D auto-passport

USER auto-passport
COPY --from=builder --chown=auto-passport /app/auto-passport-api .

CMD ["./auto-passport-api"]