# Robot API ガイド

この document は、robots が Yubi App platform と通信し、teleoperation episodes を実行する方法を説明します。

## 概要

robots は **Robot Device API** (`/api/robot/*`) を通じて system とやり取りします。認証方式は 2 つあります。

```text
Robot
  │  X-API-Key header (推奨)
  │  — または —
  │  X-User-ID + X-Robot-ID headers (fallback)
  ▼
Backend API (port 8000)
  ├─ Auth middleware: API key または user/robot headers を検証
  ├─ Robot Device endpoints: /api/robot/*
  ▼
PostgreSQL
```

## 認証

### Method 1: API Key (推奨)

Web UI から発行した API key を使用します。production deployment ではこの方式を推奨します。

| Header | Required | 説明 |
|--------|----------|------|
| `X-API-Key` | Yes | admin UI で発行した API key |

backend は key を hash 化して database で lookup し、active かどうか (revoked でない、expired でない) を検証します。key に紐づく user と robot が request context に設定されます。

```bash
curl -H "X-API-Key: <your-api-key>" http://localhost:8000/api/robot/me
```

**API key の発行方法**

1. Admin user として login する
2. navigation の **API Keys** に移動する
3. **"Create API Key"** をクリックする
4. key を紐づける robot を選択する
5. raw key を copy する (一度しか表示されません)

### Method 2: X-User-ID + X-Robot-ID Headers (Fallback)

`X-API-Key` header がない場合、backend は header-based authentication に fallback します。

| Header | Required | 説明 |
|--------|----------|------|
| `X-User-ID` | Yes | robot を操作する user の UUID |
| `X-Robot-ID` | Yes | robot の UUID |

backend は user と robot が database に存在し、同じ organization に属していることを検証します。

```bash
curl -H "X-User-ID: <user-uuid>" -H "X-Robot-ID: <robot-uuid>" \
  http://localhost:8000/api/robot/me
```

### Error Responses

```json
{"error": "Invalid API key"}
{"error": "X-API-Key or X-User-ID header is required"}
{"error": "User not found"}
{"error": "Robot not found"}
{"error": "User and robot do not belong to the same organization"}
```

### Setup

この guide の examples は API key authentication を使います。

```bash
export BASE_URL="http://localhost:8000"
export API_KEY="<your-api-key>"
```

> **Tip**: API key は Web UI から発行できます (Admin → API Keys → Create)。代わりに `backend/internal/database/seeder/initial_data.sql` の UUID を使って header-based auth も利用できます。

## 前提条件

robot が episodes を実行する前に、Web UI で次を準備してください。

1. **robot を登録する** — `/web/robots` で robot entry を作成
2. **subtasks を持つ task を作成する** — `/web/tasks` で作業内容を定義
3. **episodes を作成する** — `/web/episodes` で task + robot + location を割り当て

手順は [ユーザーガイド](user-guide.md) を参照してください。

## Episode Execution Flow

robot が episode を実行する典型的な sequence です。

```text
[Preparation - Web UI]
1. episodes を作成する (status: Ready)

[Execution - Robot API]
2. 利用可能な episodes を一覧取得   GET  /api/robot/episodes
3. episode details を取得           GET  /api/robot/episodes/{id}
4. episode を開始                  POST /api/robot/episodes/{id}/start
5. 各 subtask について:
   a. execution を作成             POST /api/robot/episodes/{id}/subtasks/{subtaskId}/executions
   b. execution を開始             POST /api/robot/episodes/{id}/subtasks/{subtaskId}/executions/{execId}/start
   c. robot が task を実行し data を収集
   d. execution を終了             POST /api/robot/episodes/{id}/subtasks/{subtaskId}/executions/{execId}/finish
   e. subtask を complete          POST /api/robot/episodes/{id}/subtasks/{subtaskId}/complete
6. episode を終了                  POST /api/robot/episodes/{id}/finish

[Review - Web UI]
7. /web/episodes/{id} で結果を確認
```

## API Reference

### Get Robot Info

robot credentials が有効か確認し、現在の configuration を取得します。

```bash
curl -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/me"
```

**Response (200)** — `Robot` object を返します。

```json
{
  "id": "c2f8e62b-ea23-4a50-8660-d707e4d5c2bc",
  "name": "Sample Robot",
  "robot_type": "yubi-stationary",
  "status": 0,
  "organization_id": "7bfbe942-5fd6-4525-ac13-0356147c202b",
  "organization_name": "Sample Organization",
  "location_id": "91154897-df4b-4b39-8c4c-b48daf4a3b37",
  "location_name": "Sample Location",
  "leader_status": null,
  "battery_level": 85,
  "last_heartbeat_at": "2026-01-01T10:00:00Z",
  "offline_reason": null,
  "robot_config": {
    "host": "localhost",
    "port": 9090,
    "cameras": [{"namespace": "camera_0", "name": "Front Camera"}]
  },
  "active_episode_id": null,
  "active_user_id": null,
  "active_operator": null
}
```

その他の fields は `openapi.yaml` の `Robot` schema を参照してください。

### List Episodes

この robot に割り当てられた episodes を取得します。`Ready` status の episodes のみ返ります。

```bash
curl -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/episodes"
```

**Response (200)** — `Episode` object の array を返します。

```json
[
  {
    "id": "358637c3-bf03-4408-845f-f9b189ec767f",
    "user_id": "69fad3df-d73f-45e1-9fb4-df52bd4857b0",
    "robot_id": "c2f8e62b-ea23-4a50-8660-d707e4d5c2bc",
    "location_id": "91154897-df4b-4b39-8c4c-b48daf4a3b37",
    "task_id": "6013935a-ab9c-4bd8-b59d-49958f516d47",
    "task_name": "Sample Task",
    "task_version_id": "5437b101-6d9d-495f-a4e8-45420eb10d99",
    "status": 0,
    "created_at": "2026-01-01T09:00:00Z"
  }
]
```

`status` values: `0` = Ready, `1` = Recording, `2` = Cancel, `3` = Completed

### Get Episode Details

subtasks を含む full episode information を取得します。robot が実行すべき subtasks を把握するために使います。

```bash
curl -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}"
```

**Response (200)** — subtasks を含む `Episode` object を返します。

```json
{
  "id": "358637c3-bf03-4408-845f-f9b189ec767f",
  "task_name": "Sample Task",
  "status": 0,
  "created_at": "2026-01-01T09:00:00Z",
  "subtasks": [
    {
      "id": "b1c352f7-54c6-4d74-ba6d-aa1352dfaee0",
      "subtask_id": "7065d47f-8de7-4b6d-af34-4aa924dfa98e",
      "name": "Sample SubTask 1",
      "order_index": 0,
      "status": 0,
      "executions": []
    }
  ]
}
```

subtask operation で使うのは `subtask_id` ではなく subtask の `id` です。

### Start Episode

robot が episode 実行を開始することを system に通知します。episode は `Ready` から `Recording` へ遷移し、robot status は `Online` から `Busy` になります。

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"occurred_at": "2026-01-01T10:00:00Z"}' \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/start"
```

| Field | Required | Format | 説明 |
|-------|----------|--------|------|
| `occurred_at` | Yes | RFC 3339 | robot 上で実際に episode が開始した時刻 |

### Create Execution

subtask を実行する前に、robot は execution record を作成します。subtask は retry により複数 executions を持てます。

```bash
SUBTASK_ID="b1c352f7-54c6-4d74-ba6d-aa1352dfaee0"

curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/executions"
```

**Response (201)**

```json
{
  "execution_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

以降の start/finish/cancel calls で使うため、`execution_id` を保存してください。

### Start Execution

execution attempt の開始を記録します。

```bash
EXECUTION_ID="a1b2c3d4-e5f6-7890-abcd-ef1234567890"

curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"occurred_at": "2026-01-01T10:05:00Z"}' \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/executions/${EXECUTION_ID}/start"
```

| Field | Required | Format | 説明 |
|-------|----------|--------|------|
| `occurred_at` | Yes | RFC 3339 | execution が実際に開始した時刻 |

### Finish Execution

execution attempt が成功したことを記録します。

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"occurred_at": "2026-01-01T10:10:00Z"}' \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/executions/${EXECUTION_ID}/finish"
```

### Cancel Execution

進行中の execution を cancel します。障害物検知などで中断した場合に使います。retry する場合は新しい execution を作成してください。

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/executions/${EXECUTION_ID}/cancel"
```

### Complete Subtask

少なくとも 1 つの execution が finish した後、subtask を completed として mark します。

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/complete"
```

### Skip Subtask

現在の環境では不要などの理由で、subtask を実行せず skip します。

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/subtasks/${SUBTASK_ID}/skip"
```

### Finish Episode

episode を completed として mark します。robot status は `Busy` から `Online` に戻ります。全 subtasks を complete または skip した後に呼び出してください。

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"occurred_at": "2026-01-01T11:00:00Z"}' \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/finish"
```

| Field | Required | Format | 説明 |
|-------|----------|--------|------|
| `occurred_at` | Yes | RFC 3339 | episode が実際に終了した時刻 |

### Cancel Episode

進行中 episode を cancel します。hardware failure や safety stop などで task を完了できない場合に使います。robot status は `Busy` から `Online` に戻ります。

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  "${BASE_URL}/api/robot/episodes/${EPISODE_ID}/cancel"
```

### Repeat Last Episode

最後に completed になった episode と同じ task、robot、location で新しい episode を作成します。同じ task を繰り返し実行したい場合に便利です。

```bash
curl -X POST \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  "${BASE_URL}/api/robot/episodes/repeat-last"
```

**Response (201)** — 新しく作成された episode object (status: Ready)。実行開始には別途 Start Episode を呼び出します。

**Error (404)** — この robot に previous completed episode がありません。

### Update Robot Status

robot の現在 status を heartbeat として報告します。backend はこれを使って robot が online かどうかを判断します。定期的に送信してください (例: 10-30 秒ごと)。

request body は `RobotStatusUpdateRequest` で、`robot_type`、`reported_at`、battery/connection/metrics を含む `status` object を持ちます。

```bash
curl -X PUT \
  -H "X-API-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "robot_type": "yubi-stationary",
    "reported_at": "2026-01-01T10:00:00Z",
    "status": {
      "battery": {"pct": 85, "charging": false},
      "connection": {"quality_pct": 100},
      "uptime_sec": 3600.0,
      "metrics": []
    }
  }' \
  "${BASE_URL}/api/robot/status"
```

| Field | Required | 説明 |
|-------|----------|------|
| `robot_type` | Yes | robot type identifier (`yubi-stationary` or `yubi-portable`) |
| `reported_at` | Yes | status report の timestamp (RFC 3339) |
| `status` | Yes | `RobotStatusDetail` object |

**RobotStatusDetail fields**

| Field | Required | 説明 |
|-------|----------|------|
| `battery` | Yes | `{"pct": int, "charging": bool}` — battery percentage |
| `connection` | Yes | `{"quality_pct": int}` — connection quality (0-100) |
| `uptime_sec` | Yes | robot uptime in seconds |
| `metrics` | No | custom robot metrics の array |
| `gate_conditions` | No | gate condition status |

**RobotStatus values** (`Robot` object 内の表示)

| Value | Status |
|-------|--------|
| 0 | Online |
| 1 | Busy |
| 2 | Offline |
| 3 | Faulted |
| 4 | Maintenance |
| 5 | Ready |

> robot の status は heartbeat data に基づき backend が判断します。request body の `status` field は detailed status report であり、上記 enum value ではありません。

## Common Error Responses

すべての endpoints は次の errors を返す可能性があります。

| Status | 意味 | Example |
|--------|------|---------|
| 401 | authentication failed | missing header, user/robot not found, org mismatch |
| 403 | insufficient permissions | user role が operation を許可しない |
| 404 | resource not found | episode or subtask ID が存在しない |
| 409 | conflict | episode が期待 status ではない |
| 500 | internal server error | unexpected backend failure |

error response format:

```json
{
  "error": "Human-readable error message"
}
```

## Status Transitions

### Episode Status

```text
Ready (Web UI で作成)
    │ start
    ▼
Recording (robot executing)
    ├─ finish ──→ Completed
    └─ cancel ──→ Cancelled
```

### Subtask Status

```text
Ready
    ├─ create execution + start ──→ InProgress
    │                                  ├─ complete ──→ Completed
    │                                  └─ cancel ───→ Cancelled
    └─ skip ──→ Skipped
```

### Robot Status

```text
Offline (initial)
    │ heartbeat received
    ▼
Online (ready)
    │ episode start
    ▼
Busy (executing)
    │ episode finish/cancel
    ▼
Online
    │ heartbeat timeout
    ▼
Offline
```

## Real-time Updates (SSE)

platform は real-time monitoring 用に Server-Sent Events を提供します。主に Web UI で使われますが、robot-side tools から consume することもできます。

| Endpoint | 説明 |
|----------|------|
| `GET /api/robots/{robotId}/status/stream` | single robot status updates |
| `GET /api/robots/status/stream?robotIds=id1,id2` | multiple robot status updates |
| `GET /api/episodes/{episodeId}/stream` | episode progress updates |
| `GET /api/robots/{robotId}/teleop/stream` | combined teleop stream (status + episode + task) |

これらの endpoints には `X-User-ID` header が必要です (`X-Robot-ID` ではありません)。

> これらの SSE endpoints は server code で直接登録されており、OpenAPI spec には含まれていません。

## Full OpenAPI Specification

request/response schemas を含む完全な API reference は [`openapi/openapi.yaml`](../../openapi/openapi.yaml) を参照してください。
