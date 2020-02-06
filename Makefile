BINARY=invoice

build:
	go build -o ${BINARY} main.go

test:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker:
	docker build -t go-clean-arch .

run:
	docker-compose up -d

stop:
	docker-compose down

.PHONY: clean install test build docker run stop vendor