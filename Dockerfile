FROM golang:1.12-alpine3.10 AS gobuilder

RUN apk add --no-cache curl git

WORKDIR /usr/local/goes
COPY . .

RUN go mod vendor
RUN go build
RUN CGO_ENABLED=0 go build -a -o goapp

# Final stage
FROM alpine:3.10

RUN apk --no-cache add ca-certificates
WORKDIR /goes
COPY . .
COPY --from=gobuilder /usr/local/goes/goapp /goes/
RUN chmod +x /goes/goapp
WORKDIR /goes
ENTRYPOINT ["./goapp"]