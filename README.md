# diag-upload-service-go

## About
This is a golang implementation of the [diag-upload-service](https://github.com/criblio/diag-upload-service)

## Feature enhancements

* Accepts only *.tgz files

## Build and Run locally

```bash
git clone https://github.com/packetsnscripts/diag-upload-service-go.git
cd diag-upload-service-go
docker compose up --build -d

# Test All (Install go)
docker build . --tag diag-upload-service-go
go test -v

#Test Upload
touch upload.tgz

curl --location --request POST 'http://localhost:8000/upload' \
--form 'diag=@"upload.tgz"'

#Test Download
#Browse http://localhost:8000/download/upload.tgz. You should be prompted to save the file


#Teardown
docker compose down
```

## Architecture
Currently this is implemented in a single container with non-persistent storage.


## References
https://echo.labstack.com

https://docs.docker.com/language/golang