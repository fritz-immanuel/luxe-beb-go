FROM golang:1.21-alpine as builder

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

COPY --from=builder /go/bin /bin

RUN mkdir /filestore
RUN chmod -R 0777 /filestore

RUN mkdir /html
COPY /html /html

RUN mkdir /luxe-beb-data
ARG version
RUN echo "$version" >> /luxe-beb-data/.version

USER nobody:nobody
ENTRYPOINT ["/bin/luxe-beb-go"]