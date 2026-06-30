# 認証とワークスペース設定

Yubi App の Web UI は、Google sign-in に Auth.js / NextAuth を使います。ローカル評価用として、`DEFAULT_USER_ID` による開発用 fallback も残しています。

## 認証モデル

Google sign-in 後、frontend は `POST /api/auth/google/session` を通じて対応する backend user を作成または取得します。その後、frontend の server-side API client は backend に以下の headers を送ります。

- `X-User-ID`: Auth.js session、または開発用の `active_user_id` / `DEFAULT_USER_ID` fallback から解決されます。
- `X-Organization-ID`: Auth.js session の active organization、または開発用の `active_organization_id` cookie から解決されます。

backend は、その user が存在し、active organization に対する `organization_membership` を持っているか確認します。active organization が未指定の場合は、その user の最初の membership を使います。

provisioning endpoint は、Yubi user が存在する前に呼ばれるため、通常の user 認証の外側に登録されています。本番環境では frontend / backend の両方で `AUTH_INTERNAL_API_SECRET` を設定して保護してください。

## 必須のローカル設定

1. 環境変数ファイルをコピーします。

```bash
cp backend/.env.example backend/.env
cp frontend/.env.sample frontend/.env
```

2. サービスを起動し、DB migration と seed を実行します。

```bash
make up PLATFORM=arm64
make migrate
make seed
```

3. `frontend/.env` が seed 済みの default admin user を指していることを確認します。

```dotenv
DEFAULT_USER_ID=69fad3df-d73f-45e1-9fb4-df52bd4857b0
```

seed data は、この user、sample organization、および両者をつなぐ admin `organization_membership` を作成します。

## Google OAuth 設定

Google OAuth client を作成し、callback URL に以下を設定します。

```text
http://localhost:3000/web/api/auth/callback/google
```

デプロイ環境では、origin を公開frontend originに置き換え、path は `/web/api/auth/callback/google` のままにします。

frontend の環境変数:

```dotenv
AUTH_SECRET=<random session secret>
AUTH_URL=http://localhost:3000/web
AUTH_GOOGLE_ID=<google oauth client id>
AUTH_GOOGLE_SECRET=<google oauth client secret>
AUTH_INTERNAL_API_SECRET=<shared frontend/backend internal secret>
```

backend 側にも同じ値を設定します。

```dotenv
AUTH_INTERNAL_API_SECRET=<same shared secret>
```

ローカル開発では `AUTH_INTERNAL_API_SECRET` は両方とも空でも動作します。ただし、backend が信頼された server network の外から到達できる環境では空にしないでください。

## Dashboard が 403 になる場合

dashboard の 403 は、多くの場合「user header は認識できたが、その user が organization に対して認可されていない」状態です。

以下を確認してください。

- 現在の schema/migration に対して `make seed` を実行済みである。
- `frontend/.env` の `DEFAULT_USER_ID` が `69fad3df-d73f-45e1-9fb4-df52bd4857b0`、またはDBに存在する別userである。
- その user に `organization_membership` が1件以上ある。
- browser に古い `active_organization_id` cookie が残っている場合は削除するか、user が所属する organization に切り替える。

ローカル環境を作り直す場合:

```bash
make reset
make up PLATFORM=arm64
make migrate
make seed
```

その後、`http://localhost:3000/web` を開いてください。

## 開発用 fallback

Google session がない場合でも、ローカル開発では `DEFAULT_USER_ID` にfallbackできます。Google sign-in を必須にしたいデプロイ環境では、frontend 環境変数から `DEFAULT_USER_ID` を外してください。
