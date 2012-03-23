# go-nbt

This is a Go package used for parsing the NBT files used throughout Minecraft. For now, it only supports reading NBT files (both gzipped and not). Eventually I will add writing capabilities.

## Usage

### Reading

```go
file, _ := ioutil.ReadFile("somefile.nbt")
z := bytes.NewReader(file)
data, _ := nbt.DecodeGzip(z)

var name string = data.Name()
var list []int32 = data.List("some list").Ints()
var compound *nbt.Compound = data.Compound("some compound")
```

See the test file (`nbt_test.go`) for more test cases.

## Suggestions, comments, hatemail

Contact moshee on Rizon or Freenode.
