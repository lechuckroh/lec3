# Requirements
* [Go](https://golang.org/)
* [gb](https://getgb.io)

# Build
To build all:
```bash
$ gb build -ldflags="-s -w" all
```

To build specific executable:
```bash
$ gb build -ldflags="-s -w" [packageName]
```

Available package names are:
* `lec-conv`
* `lec-ip`

After successful build, you can find executable files under `./bin` directory.
