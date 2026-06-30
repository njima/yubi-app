# 認証とワークスペース設定

Yubi App の Web UI は、現時点では開発用の簡易認証を使っています。Google OAuth は今後の想定に入っていますが、ローカルアクセスに必要な設定としてはまだ組み込まれていません。

## 現在のローカル認証モデル

frontend の server-side API client は、backend に以下の headers を送ります。

- `X-User-ID`: `active_user_id` cookie、または `DEFAULT_USER_ID` から解決されます。
- `X-Organization-ID`: `active_organization_id` cookie がある場合に送られます。

backend は、その user が存在し、active organization に対する `organization_membership` を持っているか確認します。active organization が未指定の場合は、その user の最初の membership を使います。

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

## Google OAuth の状態

現時点では、Google OAuth はアプリに設定されていません。今回のブランチでは、`google_sub` と personal workspace provisioning により Google 認証ユーザー向けの backend model は準備していますが、frontend はまだ上記の開発用 session を使っています。

Google OAuth を実装する場合、本番設定として少なくとも以下が必要になります。

- Google OAuth client ID
- Google OAuth client secret
- OAuth redirect/callback URL
- session signing/encryption secret
- 許可するdomain、またはuser admission policy

frontend の認証 layer が追加されるまでは、ローカルアクセスに Google login は不要です。
