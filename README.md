# go-units

Conveniently marshal/unmarshal units to/from json and yaml.

封装一些常见的单位, 实现json.Marshaler/json.Unmarshaler, yaml.Marshaler/yaml.Unmarshaler接口, 方便的将各种单位从配置中解析或格式化为字符串.

## Supported units

- [x] FileSize (B/KB/MB/GB/TB/PB)
- [ ] Duration (ns/us/ms/s/m/h/d)
- [ ] ...

## Import

`import "github.com/vj1024/go-units"`

## Quick start

```go
package main

import (
        "encoding/json"
        "fmt"
        "github.com/vj1024/go-units"
)

func main() {
        // Marshal json
        v := map[string]interface{}{
                "size1": 10 * units.KB,
                "size2": 10 * units.MB,
        }
        bs, _ := json.Marshal(&v)
        fmt.Printf("marshal result: %s\n", bs) // marshal result: {"size1":"10KB","size2":"10MB"}

        // Unmarshal json
        var st struct {
                Size1 units.FileSize `json:"size1"`
                Size2 units.FileSize `json:"size2"`
        }
        _ = json.Unmarshal(bs, &st)
        fmt.Printf("unmarshal result size1:%d, size2:%d\n", st.Size1, st.Size2) // unmarshal result size1:10240, size2:10485760
}
```
