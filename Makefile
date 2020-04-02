vendor:
	go mod vendor

.PHONY: build
build: vendor
	go build -o build/docker-machine-driver-hyperkit
	sudo chown root:wheel build/docker-machine-driver-hyperkit
	sudo chmod u+s build/docker-machine-driver-hyperkit

link:
	ln -s $(PWD)/build/docker-machine-driver-hyperkit /usr/local/bin/docker-machine-driver-hyperkit
