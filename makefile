all:
	rm -rf target
	go build -o target/postman -v src/bin/main.go
