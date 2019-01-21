FROM golang:1.11.4-alpine3.8 as builder

WORKDIR /build

RUN adduser -S -u 10001 scratchuser && \
    apk add --no-cache git

COPY go.mod go.sum main.go ./
RUN go mod vendor

RUN CGO_ENABLED=0 go build -mod=vendor -o vault_exporter

FROM scratch

EXPOSE 9410

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/vault_exporter /

USER scratchuser

CMD ["/vault_exporter"]
