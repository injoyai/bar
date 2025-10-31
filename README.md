### 实现控制台上的进度条

### 如何使用
* 普通模式
```go
import (
    "github.com/injoyai/bar"
)

func main() {
    b := bar.New(
        bar.WithTotal(100),
        bar.WithFlush(),
    )
    defer b.Close()
    for i := 0; i < 100; i++ {
        b.Add(1)
        b.Flush()
        time.Sleep(time.Millisecond * 500)
    }
}
```
* 协程模式
```go
import (
    "github.com/injoyai/bar"
)

func main() {
    // 能自动增长,和并发数量控制
    b := bar.NewCoroutine(100,2)
    for i := 0; i < 100; i++ {
        b.Go(func () {
            <time.After(time.Second * 1)
        })
    }
    b.Wait()
}
```