# Imgr-rszr

Image resizer command line tool. I shall use this for my blog images before uploading it to my server.

### Usage
```
go get github.com/dcefram/imgr-rszr

# cd to imgr-rszr

go install
```

```
img-rszr -i "path/to/images" -o "path/to/output/folder" -height 720
```

`-i` can take a path to a folder filled with png and jpegs, or the exact path of a single image file.

omitting `-i` would default to the current directory.