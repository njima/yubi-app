# Public SO101 Collector

The public SO101 collector is a guest-accessible workflow at:

```text
/public/so101
```

It is intentionally outside the authenticated app layout. Opening this route does not require Google sign-in and does not create a backend user, organization, robot, task, or episode.

## Current MVP Flow

The current screen is a frontend-only shell with a mocked local agent:

1. Connect local agent.
2. Run motor check.
3. Run calibration.
4. Select a built-in generic task.
5. Start and stop a local recording.
6. Download a local JSON manifest.

The downloaded manifest is a placeholder for the future local dataset artifact produced by the SO101 bridge or `yubi-agent`.

## Data Boundary

Guest state stays in the browser and local agent. The backend is not called in the guest-only workflow.

When a guest signs in with Google in a later slice, the app can upload the completed local artifact into the authenticated user's active organization. Until upload is explicitly implemented, guest downloads are local-only.

## Local Agent Contract Direction

The mocked workflow will be replaced by a local HTTP/WebSocket bridge exposing:

- health/version
- device discovery
- motor check status
- calibration status/result
- teleoperation status
- recording status
- local artifact manifest
- upload status
