# Public SO101 Collector

The public SO101 collector is a guest-accessible workflow at:

```text
/public/so101
```

It is intentionally outside the authenticated app layout. Opening this route does not require Google sign-in and does not create a backend user, organization, robot, task, or episode.

## Current MVP Flow

The current screen is a frontend-only shell that checks a local agent health endpoint before enabling guest SO101 workflow steps:

1. Connect local agent.
2. Run motor check.
3. Run calibration.
4. Select a built-in generic task.
5. Start and stop a local recording.
6. Download a local JSON manifest.

The local agent defaults to:

```text
http://127.0.0.1:32101
```

Override it with `NEXT_PUBLIC_SO101_LOCAL_AGENT_URL` if LeLab or `yubi-agent` exposes a different local bridge URL. The downloaded manifest is a placeholder for the future local dataset artifact produced by the SO101 bridge or `yubi-agent`.

## Data Boundary

Guest state stays in the browser and local agent. The backend is not called in the guest-only workflow.

When a guest signs in with Google in a later slice, the app can upload the completed local artifact into the authenticated user's active organization. Until upload is explicitly implemented, guest downloads are local-only.

## Local Agent Contract Direction

The first wired contract is:

```text
GET /health
```

Expected response shape:

```json
{
  "ok": true,
  "name": "lelab",
  "version": "0.1.0",
  "robotType": "so101",
  "capabilities": ["health"]
}
```

The remaining workflow will be backed by a local HTTP/WebSocket bridge exposing:

- health/version
- device discovery
- motor check status
- calibration status/result
- teleoperation status
- recording status
- local artifact manifest
- upload status
