# Deployment

Create `deploy/.env` using `deploy/.env.example` as a template. Do not overwrite an existing `.env`.

- `HONEYPOT_DOMAIN`: the hostname Caddy serves, without a scheme or path.
- `REVEALD_WEBHOOK_URL`: the complete log ingestion webhook URL.

From the repository root:

```sh
chmod 600 deploy/.env
docker compose --env-file deploy/.env -f deploy/docker-compose.yml up -d --build
```

Compose requires both values. It passes the hostname only to Caddy and the webhook URL only to reveald. Caddy expands `{$HONEYPOT_DOMAIN}` when loading its configuration.

The optional host setup script, `deploy/script.sh`, requires `DEPLOY_USER` in its process environment. This selects the account added to the Docker group. It does not read `.env`. Supply the account explicitly rather than exporting the webhook credential to the installer.

Keep `.env` and `webhook.id` private. Git ignores both, and the Docker build excludes `deploy/` and `.git/`. Avoid sharing expanded Compose configuration, container inspection output, or logs containing these values.

Environment variables keep these values out of source control. They do not hide a public hostname from DNS, certificate-transparency logs, TLS clients, or people with access to the deployment host.
