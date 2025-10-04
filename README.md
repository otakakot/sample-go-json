# Go JSON Parsing Performance Benchmark

HTTPリクエストにおける`json.Unmarshal`、`json.NewDecoder`、`io.Pipe`のパフォーマンス比較ベンチマーク

## 📊 ベンチマーク結果

### テスト環境
- **OS**: macOS (darwin)
- **CPU**: Apple M1 (arm64)
- **Go Version**: 1.25.1

### 🔹 小さなJSONペイロード（約300バイト）

```
BenchmarkJSONUnmarshal-8            137703    7778 ns/op    8751 B/op     83 allocs/op
BenchmarkJSONDecode-8               146958    7877 ns/op    8943 B/op     85 allocs/op
BenchmarkJSONPipe-8                 104872   11334 ns/op    9409 B/op     90 allocs/op
```

| メソッド | 実行時間 | メモリ割り当て | アロケーション回数 |
|---------|---------|--------------|-----------------|
| `json.Unmarshal` ✅ | 7,778 ns/op | 8,751 B/op | 83 allocs/op |
| `json.NewDecoder` | 7,877 ns/op (+1.3%) | 8,943 B/op (+2.2%) | 85 allocs/op (+2) |
| `io.Pipe` | 11,334 ns/op (+45.7%) ❌ | 9,409 B/op (+7.5%) | 90 allocs/op (+7) |

**結論**: 小さなペイロードでは`io.Pipe`は約46%遅く、オーバーヘッドが大きい。

---

### 🔹 大きなJSONペイロード（約208KB、500件のレコード）

```
BenchmarkJSONUnmarshal_LargePayload-8     243    4889584 ns/op    2677326 B/op    45564 allocs/op
BenchmarkJSONDecode_LargePayload-8        250    4771826 ns/op    2275603 B/op    45554 allocs/op
BenchmarkJSONPipe_LargePayload-8          248    4795272 ns/op    2296054 B/op    45561 allocs/op
```

| メソッド | 実行時間 | メモリ割り当て | アロケーション回数 |
|---------|---------|--------------|-----------------|
| `json.NewDecoder` ✅ | 4.77 ms/op | 2.17 MB/op | 45,554 allocs/op |
| `io.Pipe` | 4.80 ms/op (+0.5%) | 2.19 MB/op (+0.9%) | 45,561 allocs/op (+7) |
| `json.Unmarshal` | 4.89 ms/op (+2.5%) | 2.55 MB/op (+17.5%) ❌ | 45,564 allocs/op (+10) |

**結論**: 大きなペイロードでは`json.NewDecoder`と`io.Pipe`がほぼ同等のパフォーマンス。`json.Unmarshal`は最もメモリを消費。

---

## 🎯 重要な発見

### パフォーマンス特性

1. **小さなペイロード（<1KB）**
   - `json.Unmarshal`が最速
   - `io.Pipe`は約46%遅く、goroutineのオーバーヘッドが顕著
   - `json.NewDecoder`が速度とメモリのバランスが良い

2. **大きなペイロード（>100KB）**
   - **速度**: `json.NewDecoder`と`io.Pipe`がほぼ同等で最速
   - **メモリ**: `json.NewDecoder`と`io.Pipe`が約17%効率的 ✅
   - **アロケーション**: ほぼ同等

### 各メソッドの特徴

#### `json.Unmarshal` (io.ReadAll + Unmarshal)
```go
body, err := io.ReadAll(r.Body)  // 全体をメモリにコピー
json.Unmarshal(body, &req)        // さらにパース
```
- **メリット**: 小さなペイロードで最速、ボディを複数回読める
- **デメリット**: 大きなペイロードでメモリ消費が大きい（約17%増）

#### `json.NewDecoder` (ストリーミングデコード)
```go
json.NewDecoder(r.Body).Decode(&req)  // 直接パース
```
- **メリット**: メモリ効率的、コードがシンプル、大きなペイロードで高速
- **デメリット**: 小さなペイロードでわずかに遅い（1.3%）

#### `io.Pipe` (パイプ + goroutine)
```go
pr, pw := io.Pipe()
go func() {
    io.Copy(pw, r.Body)
    pw.Close()
}()
json.NewDecoder(pr).Decode(&req)
```
- **メリット**: 大きなペイロードでメモリ効率的、非同期処理が可能
- **デメリット**: 小さなペイロードでgoroutineのオーバーヘッドが大きい（46%遅い）、複雑性が増す

---

## 📌 推奨事項

| シナリオ | 推奨メソッド | 理由 |
|---------|------------|------|
| **小さなペイロード（<1KB）** | `json.Unmarshal` または `json.NewDecoder` | ほぼ同等、`io.Pipe`は避ける |
| **中〜大きなペイロード（>100KB）** | `json.NewDecoder` ✅ | メモリ効率的でシンプル |
| **超高速処理が必須（小）** | `json.Unmarshal` | 小さなペイロードで最速 |
| **メモリが制約条件** | `json.NewDecoder` または `io.Pipe` ✅ | 約17%メモリ効率的 |
| **ストリーミング処理** | `json.NewDecoder` ✅ | シンプルで効率的 |
| **非常に大きなペイロード（>1MB）** | `json.NewDecoder` ✅ | メモリ効率が重要 |
| **複雑な非同期処理が必要** | `io.Pipe` | 柔軟性が高い |

---

## 💡 実践的なアドバイス

### `json.Unmarshal`を選ぶべき場合
- ペイロードサイズが小さい（<1KB）
- 最高速度が必要
- ボディを複数回読む必要がある
- バリデーションのために生データが必要

### `json.NewDecoder`を選ぶべき場合
- ペイロードサイズが大きい（>100KB）
- メモリ使用量を最小化したい
- ストリーミング処理が必要
- 高負荷環境でスケーラビリティが重要
- 一度だけパースすれば十分
- **最もバランスが取れた選択肢** ✅

### `io.Pipe`を選ぶべき場合
- 非同期処理やパイプライン処理が必要
- データの変換や加工を行いながらパース
- 大きなペイロード（>100KB）で柔軟な処理が必要
- ただし、通常のHTTPハンドラーでは`json.NewDecoder`で十分

### 一般的な推奨
**本番環境では`json.NewDecoder`を推奨**
- メモリ効率の向上がスケーラビリティに貢献
- 大きなペイロードでも高速
- コードがシンプル（`io.ReadAll`不要）
- `io.Pipe`は特殊なケースのみ使用

---

## 🚀 ベンチマークの実行方法

### すべてのベンチマークを実行
```bash
go test -bench=. -benchmem -benchtime=2s
```

### 小さなペイロードのみ
```bash
go test -bench='BenchmarkJSONUnmarshal-8|BenchmarkJSONDecode-8|BenchmarkJSONPipe-8' -benchmem
```

### 大きなペイロードのみ
```bash
go test -bench='LargePayload' -benchmem
```

### 詳細な統計情報付き
```bash
go test -bench=. -benchmem -benchtime=3s -count=5
```

---

## 📝 実装例

### json.Unmarshal の実装
```go
func HandleJSONMarshal(w http.ResponseWriter, r *http.Request) {
    var req map[string]any

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }

    if err := json.Unmarshal(body, &req); err != nil {
        http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(req)
}
```

### json.NewDecoder の実装
```go
func HandleJSONDecode(w http.ResponseWriter, r *http.Request) {
    var req map[string]any

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(req)
}
```

### io.Pipe の実装
```go
func HandleJSONPipe(w http.ResponseWriter, r *http.Request) {
    var req map[string]any

    // io.Pipeを使ってストリーミング処理
    pr, pw := io.Pipe()

    // goroutineでボディをパイプに書き込む
    go func() {
        defer pw.Close()
        _, err := io.Copy(pw, r.Body)
        if err != nil {
            pw.CloseWithError(err)
        }
    }()

    // パイプからJSONをデコード
    if err := json.NewDecoder(pr).Decode(&req); err != nil {
        http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(req)
}
```
