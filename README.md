# Basic Go-Postgres App

This basic app is for demonstration purposes. It is a single file connector between http and a postgres DB

### Installation

Mike-demo requires [Golang](https://golang.org/) v1.7+ to run.

### Glide
```sh
$ curl https://glide.sh/get | sh
```

```sh
$ postgres -D /usr/local/var/postgres
$ git clone git@github.com:sunilgopinath/mike-demo.git
$ cd mike-demo
$ glide update
$ mike-demo go run cmd/demo/main.go
# Creating table
# Inserting values
```
