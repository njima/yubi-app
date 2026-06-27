# Backend Status Policy

このドキュメントは、backend の status のうち、domain state として永続化するものと、API 表示用に派生させるものを定義します。新しい status logic は `backend/internal/domain/model` に置き、usecase では transition rule を複製せず domain helper を呼び出してください。

## Robot

Robot status には 2 つの意味があります。

- **永続化する operation status**: `Ready`, `Busy`, `Faulted`, `Maintenance`
- **connection-only の表示 status**: `Online`, `Offline`

`Ready` が保存時のデフォルトです。`Online` と `Offline` は、保存された operation status が `Ready` のときに Redis heartbeat から派生させます。manual operation update としては使いません。legacy data に `Online` が保存されている場合も、表示時は heartbeat rule で解決します。

Teleoperation の transition は次の通りです。

```text
Ready + heartbeat alive -> Busy
Busy -> Ready
```

`Faulted` と `Maintenance` は manual operation state です。`Busy` は teleoperation lifecycle action で制御し、manual status update で上書きしません。

## Episode

Episode status は lifecycle が所有します。

- `Ready`: 作成済みで、まだ recording していない
- `Recording`: recording 中
- `Completed`: terminal かつ successful
- `Cancel`: terminal だが successful ではない

Lifecycle transition には `Start`, `Finish`, `Cancel` を使います。直接の status update は no-op のみ許可し、lifecycle state の直接変更で domain rule を迂回しないでください。

## Episode Subtask

`EpisodeSubTask.CollectionStatus` は collection workflow state を表します。

- `Ready` と `InProgress` は open state です。
- `Completed`, `Skipped`, `Cancelled` は terminal state です。
- `Completed` と `Skipped` は workflow-resolved state です。
- successful completion は `Completed` のみです。

Transition には `StartProgress`, `Complete`, `Skip`, `Cancel` を使います。

## Episode Subtask Execution

Execution status は個別 execution attempt の状態を表します。

- `Ready` と `Started` は open state です。
- `Finished` と `Cancelled` は terminal かつ workflow-resolved state です。
- successful completion は `Finished` のみです。

Status field を直接設定せず、execution lifecycle helper を使ってください。

## Task

Task progress status は duration から派生します。

- `Planning`: actual duration が 0。
- `Doing`: actual duration が 0 より大きく、target duration より小さい。
- `Completed`: actual duration が target duration 以上。

`Canceled` は別の明示的な state であり、duration から推論しません。
