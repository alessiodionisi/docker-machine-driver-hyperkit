.PHONY: build
build:
	go build -o ./build/docker-machine-driver-hyperkit ./cmd/main.go
	sudo chown root:wheel ./build/docker-machine-driver-hyperkit
	sudo chmod u+s ./build/docker-machine-driver-hyperkit

link:
	rm -f /usr/local/bin/docker-machine-driver-hyperkit
	ln -s $(PWD)/build/docker-machine-driver-hyperkit /usr/local/bin/docker-machine-driver-hyperkit

build-and-link: build link
