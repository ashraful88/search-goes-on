FROM golang:1.12-alpine3.10 AS goesbuilder

RUN apk add --no-cache curl git
RUN mkdir -p /etc/secret/
COPY  .env /etc/secret/.env

WORKDIR /usr/local/goes
COPY . .

ENV TZ Asia/Kuala_Lumpur
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN curl https://glide.sh/get | sh
RUN glide install
RUN cd api && CGO_ENABLED=0 go build -a -o goapp


# Final stage
FROM alpine:3.10

RUN mkdir -p /etc/secret/
COPY  .env /etc/secret/.env

ENV TZ Asia/Kuala_Lumpur
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk --no-cache add ca-certificates

WORKDIR /usr/local/goes
COPY . .
COPY --from=goesbuilder /usr/local/goes/goapp /usr/local/goes/
RUN chmod +x /usr/local/goes/goapp
ENTRYPOINT ["./goapp"]