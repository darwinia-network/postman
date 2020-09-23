# Postman

[![Shadow][workflow-badge]][github]
[![LICENSE](https://img.shields.io/github/license/darwinia-network/postman)](https://choosealicense.com/licenses/gpl-3.0/)

Check the pending headers and emit mails to mailing list.

## Usage

```
# Run your prometheus alert manager
$ docker run --rm -it -p 9093:9093 prom/alertmanager:v0.21.0

# Download postman
$ go get -u darwinia-network/postman

# Run postman
$ postman
```

## ENVIROMENTS

```golang

var (
	SECONDS    = 10
	PROMETHEUS = "http://0.0.0.0:9093/api/v2"
	SHADOW     = "https://testnet.shadow.darwinia.network"
	ENDPOINT   = "wss://crab.darwinia.network"
)
```

## LICENSE

GPL-3.0


[github]: https://github.com/darwinia-network/postman
[workflow-badge]: https://github.com/darwinia-network/postman/workflows/postman/badge.svg
