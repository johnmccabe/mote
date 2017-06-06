FROM arm32v6/golang:1.8-alpine
RUN apk update
RUN apk add git gcc musl-dev linux-headers
RUN go get github.com/johnmccabe/mote
RUN go get github.com/Masterminds/glide
WORKDIR /go/src/github.com/johnmccabe/mote/
COPY rainbow.go ./example/rainbow/
RUN glide install
RUN go build example/rainbow/rainbow.go

FROM arm32v6/alpine
LABEL org.label-schema.docker.cmd "docker run -ti --device=/dev/ttyACM0 johnmccabe/mote-rainbow"
LABEL org.label-schema.description "Run the Pimoroni Mote Rainbow example on your Raspberry Pi (arm32v6)"
LABEL org.label-schema.url "https://github.com/johnmccabe/mote"
COPY --from=0 /go/src/github.com/johnmccabe/mote/rainbow /rainbow
CMD ["/rainbow"]