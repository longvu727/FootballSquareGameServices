FROM golang:1.22.4-alpine3.20 AS build
RUN apk add --no-cache git \
                openssh-client \
                ca-certificates

ENV GOPRIVATE="github.com/longvu727/FootballSquaresLibs"
RUN git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"

RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh && \
    ssh-keyscan gitlab.com > /root/.ssh/known_hosts &&\
    chmod 644 /root/.ssh/known_hosts && touch /root/.ssh/config \
    && echo "StrictHostKeyChecking no" > /root/.ssh/config

COPY env/.ssh/id_* /root/.ssh/

ENV GOOS=linux GOARCH=amd64

WORKDIR /api

COPY go.mod ./

RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags "-s -w" -o api main.go

####Debug####
#RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
#RUN go build -gcflags "all=-N -l" -o api main.go


FROM build AS runtime

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

CMD ["./api"]

####Debug####
#CMD [ "/go/bin/dlv", "--listen=:2101", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "./api" ]
