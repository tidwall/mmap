# mmap

Load a huge file into a byte slice without actually reading any data.

## Installing

```
go get -u github.com/tidwall/mmap
```

## Using

Load a bigole file into a byte slice. This happens pretty much instantly even
if your file is many GBs.

```go
data, err := mmap.Open("my-big-file.txt")
if err != nil {
    panic(err)
}
```

Now you can read the `data` slice like any other Go slice.

Make sure to release the data when your done.

```go
mmap.Close(data)
```

Don't read the `data` after closing otherwise your f*cked.

That's all, bye now
