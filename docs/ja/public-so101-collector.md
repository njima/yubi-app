# Public SO101 Collector

public SO101 collector は、guest がアクセスできる workflow です。

```text
/public/so101
```

この route は authenticated app layout の外側にあります。Google sign-in は不要で、backend user、organization、robot、task、episode は作成しません。

## 現在の MVP Flow

現在の画面は、guest SO101 workflow step を有効化する前に local agent health endpoint を確認する frontend-only shell です。

1. local agent に接続する。
2. motor check を実行する。
3. calibration を実行する。
4. built-in generic task を選択する。
5. local recording を開始・停止する。
6. local JSON manifest を download する。

local agent の default は以下です。

```text
http://127.0.0.1:32101
```

LeLab または `yubi-agent` が別の local bridge URL を公開する場合は、`NEXT_PUBLIC_SO101_LOCAL_AGENT_URL` で上書きしてください。download される manifest は、今後 SO101 bridge または `yubi-agent` が生成する local dataset artifact の placeholder です。

## Data Boundary

guest state は browser と local agent に留まります。guest-only workflow では backend を呼びません。

後続sliceで guest が Google sign-in した場合、完了済み local artifact を authenticated user の active organization に upload できるようにします。upload が明示的に実装されるまでは、guest download は local-only です。

## Local Agent Contract Direction

最初に接続する contract は以下です。

```text
GET /health
```

期待する response shape は以下です。

```json
{
  "ok": true,
  "name": "lelab",
  "version": "0.1.0",
  "robotType": "so101",
  "capabilities": ["health"]
}
```

残りの workflow は、以下を公開する local HTTP/WebSocket bridge で支える予定です。

- health/version
- device discovery
- motor check status
- calibration status/result
- teleoperation status
- recording status
- local artifact manifest
- upload status
