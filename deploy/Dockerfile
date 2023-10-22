FROM golang:1.20-alpine as builder
ARG SWAGGER_HOST_PORT

RUN apk update && \
    apk add --no-cache make && \
    mkdir -p /app
WORKDIR /app

# for cache go mod depends
COPY go.mod .
COPY go.sum .
RUN go mod download


ENV SWAGGER_HOST_PORT=$SWAGGER_HOST_PORT
COPY . .
RUN make build && \
    sed -i "s/\"host\": \"localhost:8080\"/\"host\": \"${SWAGGER_HOST_PORT}\"/" docs/swagger.json

FROM alpine
EXPOSE 8080
RUN adduser -D -H -h /app effectivemobile && \
    mkdir -p /app/internal/migrations && \
    chown -R effectivemobile:effectivemobile /app
WORKDIR /app
USER effectivemobile



COPY --chown=effectivemobile --from=builder /app/bin/effectivemobile /app
COPY --chown=effectivemobile --from=builder /app/internal/migrations /app/internal/migrations
COPY --chown=effectivemobile --from=builder /app/docs/swagger.json /app/docs/


CMD ["/app/effectivemobile"]