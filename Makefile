build:
	go build -o gity main.go
clean-init:
	rm -r .gity
	./gity init