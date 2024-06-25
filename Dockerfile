FROM golang:1.22.4-alpine3.20 AS build

ENV GOOS=linux GOARCH=amd64

WORKDIR /api

COPY go.mod ./

RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags "-s -w" -o api main.go

FROM alpine:3.20 AS runtime

RUN apk add curl

WORKDIR /api

ENV USER=longvu727 UID=1000
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY --from=build --chown=${USER}:${USER} /api/ .

USER ${USER}:${USER}

EXPOSE 3000

CMD ["./api"]
