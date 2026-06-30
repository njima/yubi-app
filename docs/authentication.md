# Authentication and Workspace Setup

Yubi App currently uses development authentication for the web UI. Google OAuth is planned, but it is not required or wired for local access yet.

## Current local authentication model

The frontend server-side API client sends these headers to the backend:

- `X-User-ID`: resolved from the `active_user_id` cookie, or from `DEFAULT_USER_ID`.
- `X-Organization-ID`: resolved from the `active_organization_id` cookie when present.

The backend then verifies that the user exists and has an `organization_membership` for the active organization. If no active organization is selected, the backend uses the first membership for that user.

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

## Google OAuth status

Google OAuth is not configured in the current app. The current branch prepares the backend model for Google-authenticated users by adding `google_sub` and personal workspace provisioning, but the frontend still uses the development session described above.

When Google OAuth is implemented, production configuration should include at least:

- Google OAuth client ID
- Google OAuth client secret
- OAuth redirect/callback URL
- Session signing/encryption secret
- Allowed domains or user admission policy

Until that frontend auth layer is added, local access does not require Google login.
