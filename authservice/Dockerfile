FROM golang:1.6-onbuild

COPY . /go/src/app

RUN go get -d -v

RUN go install -v

# Set the PORT environment variable inside the container
ENV PORT 8080

# Expose port 8080 to the host so we can access our application
EXPOSE 8080