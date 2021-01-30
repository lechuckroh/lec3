# LEC3

## Build
```bash
# build binary
$ make build

# build linux binary
$ make build-linux

# build windows binary
$ make build-windows

# build docker image
$ make image
```

## Run
### Image Processing
#### Batch mode
Apply image filtering to all image files in specified directories.

```bash
$ ./lec ip -cfg config/ip-batch.yml

# override source directory and destination directory
$ ./lec ip -cfg config/ip-batch.yml \
    -src ./images/src \
    -dest ./images/dest

# use watch mode
$ ./lec ip -cfg config/ip-batch.yml -watch
```

#### Watch mode
Apply image filtering to all image files using watch mode.
Any new/updated files will be processed after specified delays.

```bash
$ ./lec ip -cfg config/ip-watch.yml

# override source directory and destination directory
$ ./lec ip -cfg config/ip-watch.yml \
    -src ./images/src \
    -dest ./images/dest
```

### Convert Format
converts images to other formats like pdf and zip.

```bash
$ ./lec conv -cfg config/conv-kindle3-pdf.yml

# override source file and destination file/directory
$ ./lec conv -cfg config/conv-kindle3-pdf.yml \
    -src ./images/foo.zip \
    -dest ./images/foo-kindle3.pdf
```
