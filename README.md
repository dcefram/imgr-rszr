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

Flags
```
-i # path to the folder filled with images (jpg or png only as of now)
-o # path to where we dump all proccessed images
-height # define the height we should resize the image to. Default is 720
-width #define the width we should resize the image to.
```

If you only specify `-width` or only `-height`, the tool would try to resize the image whilst 
keeping the original aspect ratio of the image.