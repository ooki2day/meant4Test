FROM golang:alpine
RUN apk add git

ADD . src/meant4Test
WORKDIR src/meant4Test

RUN go get
RUN go build
RUN chmod +x meant4Test

CMD ["./meant4Test"]
