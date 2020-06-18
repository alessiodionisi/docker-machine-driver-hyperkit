# HyperKit driver for Docker Machine

Simple and clean HyperKit driver for Docker Machine

## Requirements

This driver needs the [HyperKit](http://github.com/moby/hyperkit) and isoinfo binary, you can install it from Homebrew:

```
$ brew install hyperkit cdrtools
```

## Installation

### Binaries (WIP)

Binaries are available on [every release](https://github.com/adnsio/docker-machine-driver-hyperkit/releases).

### Homebrew (WIP)

```
$ brew install adnsio/tap/docker-machine-driver-hyperkit
```

### Build from sources

This command will build the driver and link the binary under `/usr/local/bin/docker-machine-driver-hyperkit`.

```
$ make build-and-link
```

## Options

| Flag name | Environment variable | Type | Default |
|-----------|----------------------|------|---------|
| `--hyperkit-binary` | `HYPERKIT_BINARY` | string | hyperkit |
| `--hyperkit-boot2docker-url` | `HYPERKIT_BOOT2DOCKER_URL` | string | latest version |
| `--hyperkit-cpus` | `HYPERKIT_CPUS` | int | 1 |
| `--hyperkit-disk-size` | `HYPERKIT_DISK_SIZE` | int | 20000 |
| `--hyperkit-isoinfo-binary` | `HYPERKIT_ISOINFO_BINARY` | string | isoinfo |
| `--hyperkit-memory-size` | `HYPERKIT_MEMORY_SIZE` | int | 1024 |

## License

This project is made with ♥ by [contributors](https://github.com/adnsio/docker-machine-driver-hyperkit/graphs/contributors) and it's released under the [MIT license](https://github.com/adnsio/docker-machine-driver-hyperkit/blob/master/LICENSE).
