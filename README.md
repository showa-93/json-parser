# jsonparser

レキサー、パーサーの復習でつくっただけ  
cmd/main.goを実行するとjsonを一度パースした結果を表示できる。

```bash
go run cmd/main.go test/test.json

2023/11/30 04:13:12 read file path:  test/test.json
{
    "test1": "test",
    "test2": [
        123.000000,
        123.000000,
        456.000000,
        12300.000000
    ],
    "test3": {
        "test3-1": true,
        "test3-2": null
    }
}
```

## 参考

- [RustでJSONパーサーをフルスクラッチで実装する #Rust - Qiita](https://qiita.com/togatoga/items/9d600e20325775f09547)
