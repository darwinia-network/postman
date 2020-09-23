all:
	rm -rf target
	go build -o target/postman -v postman.go
