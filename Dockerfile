# Specifies a parent image
FROM golang:1.21 as install-stage

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

# Installs Go dependencies
RUN go mod download

FROM install-stage as test-stage
RUN go test ./...

FROM install-stage as run-stage
# Builds your app with optional configuration
RUN make build

# Tells Docker which network port your container listens on
EXPOSE 3000

#ENTRYPOINT ["./bin/app"]
# Specifies the executable command that runs when the container starts
CMD ["./bin/app"]
