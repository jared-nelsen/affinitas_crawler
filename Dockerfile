FROM golang:alpine as builder

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM scratch

WORKDIR /app

COPY --from=builder /app/cabo_affinitas /usr/bin/

ENTRYPOINT ["cabo_affinitas"]