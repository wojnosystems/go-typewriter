# Overview

re-usable autotyper for presentations.

# Build

```shell script
docker build -t tw:latest .
```

# Run

```shell script
docker run --rm -ti --user $(id -u):$(id -g) -v $(pwd)/goprogram:/goprogram -v $(pwd)/typewriter:/typewriter go-typewriter:latest -- /goprogram/main.go
```
