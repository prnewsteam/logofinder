## Build on Mac

```
brew install imagemagick@6
brew install librsvg
go mod download
export CGO_CFLAGS_ALLOW='-Xpreprocessor'
go build main.go
```

## Build Docker

```
docker build -t logofinder .
docker run -it -p 80:8099 logofinder
```

## Api

Request:
```
curl -i http://host/logo?domain=cometa-ck.news&width=500&height=500
```

Response:
```
streamed png file
```

## Tests

```

cd tests && ./test.sh

```

