FROM golang:1.6-onbuild

COPY . /go/src/app

RUN go get -d -v

RUN go install -v

# Set the PORT environment variable inside the container
ENV PORT 8080

EXPOSE 8080