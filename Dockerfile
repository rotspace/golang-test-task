FROM golang:1.9.1-alpine3.6

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /go/src/github.com/rotspace/tagc
COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
