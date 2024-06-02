# Copyright 2019 Core Services Team.

FROM golang:1.15-alpine as builder

RUN apk add --no-cache ca-certificates git

WORKDIR /luxe-beb-go
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go install .

FROM alpine:3.13
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache wkhtmltopdf xvfb ttf-dejavu ttf-droid ttf-freefont ttf-liberation 
    

RUN ln -s /usr/bin/wkhtmltopdf /usr/local/bin/wkhtmltopdf
RUN chmod +x /usr/local/bin/wkhtmltopdf

COPY --from=builder /go/bin /bin

RUN mkdir /filestore
RUN chmod -R 0777 /filestore

RUN mkdir /html
COPY /html /html

RUN mkdir /angke-data
ARG version
RUN echo "$version" >> /angke-data/.version

USER nobody:nobody
ENTRYPOINT ["/bin/luxe-beb-go"]