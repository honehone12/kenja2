.PHONY:

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -tags "lambda.norpc avx" -o bootstrap main.go

zip:
	zip anime_search.zip bootstrap
