FROM golang:1.17-alpine

#Set app directory
WORKDIR /usr/app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy source files
COPY *.go ./

# Build
RUN go build -o /diag-upload-service-go

EXPOSE 8000

# RUN
CMD [ "/diag-upload-service-go" ]
