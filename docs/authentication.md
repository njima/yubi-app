# Authentication and Workspace Setup

Yubi App uses Auth.js / NextAuth for Google sign-in in the web UI. A development fallback using `DEFAULT_USER_ID` is still available for local evaluation.

## Authentication model

After Google sign-in, the frontend provisions or loads the matching backend user through `POST /api/auth/google/session`. The frontend server-side API client then sends these headers to the backend:

- `X-User-ID`: resolved from the Auth.js session, or from the development `active_user_id` / `DEFAULT_USER_ID` fallback.
- `X-Organization-ID`: resolved from the Auth.js session active organization, or from the development `active_organization_id` cookie when present.

The backend then verifies that the user exists and has an `organization_membership` for the active organization. If no active organization is selected, the backend uses the first membership for that user.

The provisioning endpoint is intentionally registered outside normal user authentication because it is called before a Yubi user exists. In production, protect it with `AUTH_INTERNAL_API_SECRET` on both frontend and backend.

## Required local setup

1. Copy the environment files:

```bash
cp backend/.env.example backend/.env
cp frontend/.env.sample frontend/.env
```

2. Start the services and seed the database:

```bash
make up PLATFORM=arm64
make migrate
make seed
```

3. Confirm `frontend/.env` uses the seeded default admin user:

```dotenv
DEFAULT_USER_ID=69fad3df-d73f-45e1-9fb4-df52bd4857b0
```

The seed data creates this user, the sample organization, and an admin `organization_membership` between them.

## Google OAuth setup

Create a Google OAuth client and set the callback URL to:

```text
http://localhost:3000/web/api/auth/callback/google
```

For deployed environments, replace the origin with the public frontend origin and keep the `/web/api/auth/callback/google` path.

Set these frontend environment variables:

```dotenv
AUTH_SECRET=<random session secret>
AUTH_URL=http://localhost:3000/web
AUTH_GOOGLE_ID=<google oauth client id>
AUTH_GOOGLE_SECRET=<google oauth client secret>
AUTH_INTERNAL_API_SECRET=<shared frontend/backend internal secret>
```

Set the matching backend environment variable:

```dotenv
AUTH_INTERNAL_API_SECRET=<same shared secret>
```

For local development, `AUTH_INTERNAL_API_SECRET` may be empty on both services. Do not leave it empty when the backend is reachable from outside the trusted server network.

## Dashboard returns 403

A dashboard 403 usually means the backend could authenticate the user header, but could not authorize the user for an organization.

Check these points:

- `make seed` has been run after the current schema/migration.
- `frontend/.env` has `DEFAULT_USER_ID=69fad3df-d73f-45e1-9fb4-df52bd4857b0`, or another user that exists in the database.
- The user has at least one row in `organization_membership`.
- If the browser has an old `active_organization_id` cookie, clear it or switch back to an organization the user belongs to.

For local reset:

```bash
make reset
make up PLATFORM=arm64
make migrate
make seed
```

Then open `http://localhost:3000/web`.

## Development fallback

If no Google session exists, local development can still fall back to `DEFAULT_USER_ID`. Remove `DEFAULT_USER_ID` from deployed frontend environments if Google sign-in should be mandatory.
