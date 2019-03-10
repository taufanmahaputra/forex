FROM golang:latest AS build

# Setting up working directory
WORKDIR $GOPATH/src/github.com/taufanmahaputra/forex
COPY . .

# Install dependencies
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

# Build inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /forex ./cmd/forex/

# Single layer image
FROM scratch
COPY --from=build /forex .
CMD ["/forex"]