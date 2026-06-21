# ユーザーガイド

## Yubi App とは

Yubi App は **ロボット遠隔操作データ収集** を管理するための platform です。物体の picking や部品組み立てなど、robot が実行する tasks を体系的に収集・追跡・review したい場合に使います。

基本的な流れは次のとおりです。

1. **Task を定義する** — 「ボトルを持ち上げ、皿の上に置く」
2. **Episode を割り当てる** — 「Robot-01 が Location A でこの task を実行する」
3. **Robot が実行する** — robot が task を取得し、実行して結果を返す
4. **結果を review する** — 何が起きたか、どれくらい時間がかかったか、成功したかを確認する

これらは Web interface で管理し、robots は REST API 経由で通信します。

## はじめに

### 初回セットアップ

application を起動したら (`make up && make migrate && make seed`)、次を開きます。

```text
http://localhost:3000/web
```

seed data により、すぐに使える環境が作成されます。

- 1 organization、1 site、1 location
- 1 admin user
- 1 robot
- 3 subtasks を持つ 1 task
- 複数状態の sample episodes

default Admin user として login しています。別 user に切り替えるには、右上の user avatar をクリックし、**"Switch Account"** を選択してください。

> **Warning**: teleoperation session や episode recording などの active operation 中に user を切り替えないでください。以後の API request の user context が変わり、進行中 session の data と不整合になる可能性があります。

### Data Model の理解

system 内の data はすべて **Organization** 配下に整理されます。

```text
Organization
  └── Site (例: "Tokyo Office")
       └── Location (例: "Lab Room A")
            └── Robot (例: "Robot-01")

Organization
  └── Task (例: "Pick and place bottle")
       └── Task Version (v1.0.0)
            ├── Subtask 1: "Detect the object"
            ├── Subtask 2: "Pick up the object"
            └── Subtask 3: "Place on the plate"

Episode = "Robot-01 performs Pick-and-place v1.0.0 at Lab Room A"
```

**Episode** は中心概念で、robot による task の 1 回の実行を表します。episode は次の状態を遷移します。

```text
Ready → Recording → Completed
                  → Cancelled
```

## Tutorial: 初めてのデータ収集

### Step 1: Workspace を準備する

data を収集する前に、robots が作業する場所を用意します。

1. **Locations** (`/web/locations`) に移動する
2. 新しい site が必要なら作成する (例: "My Lab")
3. site 内に location を作成する (例: "Workbench A")

> seed data には site と location が含まれるため、試すだけならこの step は skip できます。

### Step 2: Robot を登録する

1. **Robots** (`/web/robots`) に移動する
2. **"Create Robot"** をクリックする
3. 次の fields を入力する
   - **Location**: robot が稼働する場所
   - **Name**: 人間が読める名前 (例: "Robot-01")
   - **Robot Type**: model/type (例: "UR10e")
   - **Robot Config** (任意): camera や connection settings の JSON

```json
{
  "host": "192.168.1.100",
  "port": 9090,
  "cameras": [
    {"namespace": "camera_0", "name": "Front Camera"}
  ]
}
```

作成直後の robot status は **Offline** です。Robot API から heartbeat を送信すると **Online** になります。

### Step 3: Task を定義する

task は robot が何をするかを表します。各 task は ordered list の subtasks を持ちます。

1. **Tasks** (`/web/tasks`) に移動する
2. **"Create Task"** をクリックする
3. 次の fields を入力する
   - **Name**: 例 "Pick up bottle and place on plate"
   - **Description**: 詳細な instructions
   - **Manual URL**: reference document への link
   - **Priority / Difficulty**: 整理・filter 用
   - **Robot Type**: この task を実行できる robot model
4. task 作成後、subtasks を持つ **Task Version** を追加する
   - Subtask 1: "Detect the bottle"
   - Subtask 2: "Pick up the bottle"
   - Subtask 3: "Place on the plate"
5. episode で使えるように task version を **Approve** する

### Step 4: Episode を作成する

episode は実際の割り当てです。「この robot がこの task を実行する」を表します。

1. **Episodes** (`/web/episodes`) に移動する
2. **"Create Episode"** をクリックする
3. 次を選択する
   - **Task**: 作成した task
   - **Robot**: 実行する robot
   - **Location**: 実行場所
4. **"Create"** をクリックする

episode は **Ready** status になり、robot が取得するのを待ちます。

> **Tip**: create dialog の **count** を 2 以上にすると、同じ task/robot の episode をまとめて作成できます。

### Step 5: Robot が Episode を実行する

この step は robot 側で REST API を使って行います。詳細は [Robot API ガイド](robot-api-guide.md) を参照してください。

概要は次のとおりです。

1. 割り当て済み episodes を取得する (`GET /api/robot/episodes`)
2. episode を開始する (`POST /api/robot/episodes/{id}/start`)
3. 各 subtask を実行する (execution 作成 → start → finish → complete)
4. episode を終了する (`POST /api/robot/episodes/{id}/finish`)

robot が作業している間、episode status は **Ready** → **Recording** → **Completed** と変化します。

> **物理 robot がない場合**: [Robot API ガイド](robot-api-guide.md) の curl commands で robot を simulate できます。まず `PUT /api/robot/status` で heartbeat を送り、robot を online にしてから episode execution flow を実行してください。

### Step 6: 結果を確認する

1. **Episodes** (`/web/episodes`) に移動する
2. status (例: "Completed") で filter する
3. episode をクリックして詳細を見る
   - timing (開始時刻、各 subtask の所要時間)
   - subtask execution history (attempts, successes, skips)
   - recording data (robot が recordings を upload している場合)

## Real-time Monitoring

### Robot Status

**Robots** page (`/web/robots`) では全 robots の live status を確認できます。

| Status | 意味 |
|--------|------|
| **Online** | robot が接続済みで task を受け付け可能 |
| **Busy** | robot が episode を実行中 |
| **Offline** | robot が heartbeat を送信していない |
| **Faulted** | robot が error を報告した |
| **Maintenance** | robot が maintenance 中 |

status update は Server-Sent Events (SSE) により real-time に反映されます。page refresh は不要です。

### Teleoperation Console

**Online** の robots は teleoperation console から real-time に操作できます。

1. robot list から online robot をクリックする
2. **"Start Teleoperation"** をクリックする
3. episode を選択または作成する
4. console で subtask execution を monitor/control する

console では次を確認できます。

- live camera feeds (`robot_config` の camera settings が必要)
- subtask control panel (start, complete, skip)
- real-time robot sensor data
- episode progress

> **Prerequisite**: robot が `PUT /api/robot/status` で heartbeat を送信している必要があります。物理 robot がない場合は [Update Robot Status](robot-api-guide.md#update-robot-status) を参照してください。

### Dashboard

**Dashboard** (`/web/dashboard`) は fleet-level overview を提供します。

- 今日/今週完了した episodes 数
- 時系列の collection trends
- active robots と現在 status

## API Key Management (Admin Only)

API keys を使うと、robots は手動 header 設定なしで platform に authenticate できます。Admin users は key の発行・確認・revoke ができます。

### API Key の発行

1. **API Keys** (`/web/api-keys`) に移動する
2. **"Create API Key"** をクリックする
3. 名前を入力し、key を紐づける robot を選択する
4. **"Create"** をクリックする
5. **raw key をすぐに copy する** — 一度しか表示されず、後から取得できません

key は特定 robot と作成 user に紐づきます。robot がこの key で authenticate すると、system は robot と user の両方を自動識別します。

### API Key の revoke

1. **API Keys** (`/web/api-keys`) に移動する
2. list から key を探す
3. **"Revoke"** をクリックして確認する

revoked key は即座に使えなくなります。この操作は取り消せません。

### Key Lifecycle

- key には任意の expiration date を設定できます
- `Last Used` column で最後に authentication に使われた時刻を確認できます
- revoked keys は list に残ります ("Include revoked" toggle で表示)

## User Roles

system には 5 つの role があります。

| Role | できること |
|------|------------|
| **Admin** | users, robots, tasks, episodes, sites, locations, API keys を含む全操作 |
| **Data Engineer** | tasks, robots, episodes を管理。locations は read-only |
| **Manager** | Data Engineer と同等 |
| **Operator** | episodes の作成・更新、robots の操作。tasks/locations は read-only |
| **Viewer** | すべて read-only |

### Users の管理

1. **Users** (`/web/users`) に移動する
2. **"Create User"** で user を作成する (email, name, role)
3. access scope として locations/sites を割り当てる
4. CSV upload で bulk import する (columns: `email`, `display_name`, `role`)

## Data Import

### CSV から Tasks を import

1. Tasks に移動し、**"Import"** をクリックする
2. columns `name`, `description`, `priority`, `difficulty` などを持つ CSV を upload する
3. validation results (valid rows, duplicates, errors) を確認する
4. confirm して import する

### CSV から Users を import

1. Users に移動し、**"Import"** をクリックする
2. columns `email`, `display_name`, `role` を持つ CSV を upload する
3. 内容を確認して confirm する

## Data Export

### Episode Export

1. **Episodes** に移動する
2. filters (task, robot, date range, status) を適用する
3. **"Export"** をクリックして CSV を download する

### Operator Yield Report

1. **Dashboard** → Reports に移動する
2. date range を選択する
3. operator productivity data を CSV として export する

## Tips

- **Bookmarkable filters**: filter state は URL に保存されます。URL を copy して team members と共有できます。
- **Real-time updates**: robot status と episode progress は自動更新されます。refresh は不要です。
- **Bulk episode creation**: episode create dialog の count を 2 以上にすると、同じ task の繰り返し episode を一括作成できます。
- **Task versioning**: subtasks を更新する場合、既存 version を編集せず新しい task version を作成してください。過去 episodes で実行された内容の履歴を保持できます。
