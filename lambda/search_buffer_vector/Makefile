.PHONY:

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -tags lambda.norpc -o bootstrap main.go

zip:
	zip kenja_vector.zip bootstrap