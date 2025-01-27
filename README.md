# ofxmerge
 
**Don't use -- status is still experimental.**

`ofxmerge` is a small tool that merges 2 or more OFX files.

To install:
```shell
$ go build -o ./ofxmerge .
```

To merge with output of new OFX file being sent to standard out:
```shell
$ ./ofxmerge file1.ofx file2.ofx
```

