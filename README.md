# Imgr-rszr

Image resizer command line tool. I shall use this for my blog images before uploading it to my server.

### Usage
```bash
go get github.com/dcefram/imgr-rszr

# cd to imgr-rszr

go install

# If you have your /go/bin included in your path, you can do this:
img-rszr -i "path/to/images" -o "path/to/output/folder" -height 720

# If not...
$GOPATH/bin/img-rszr -i "path/to/image" -o "path/to/output" -width 1
```

`-i` can take a path to a folder filled with png and jpegs, or the exact path of a single image file.

omitting `-i` would default to the current directory.

**Flags**

|Flag|Description|Default Value|
|-----------|-----------|--------|
|`-i`|Path to the folder filled with images, or Path to the target image (JPEG and PNG files are only supported as of now)|`./`|
|`-o`|Path to the output folder|`./output`|
|`-height`|Define the height we should resize the image(s) to.|720|
|`-width`|Define the width we should resize the image(s) to.|Keep aspect ratio|

*\* If width is only specified, it would resize the image based on the specified width while maintain aspect ratio. Same goes if only height is specified. If none of the two is specified, then it would default to 720px height whilst keeping aspect ratio.*
