APP_NAME=merpgoledger
APP_VERSION=1.0.0

build:
	go build -o $(APP_NAME)

clean:
	rm -f $(APP_NAME)

run:
	go run .

test:
	go test ./...

release: clean build
	tar -czvf $(APP_NAME)-$(APP_VERSION).tar.gz $(APP_NAME)