# Deploy and rollback

## Deployment order

1. Back up MySQL and the current Compose environment.
2. Rotate any exposed R2 credential and load the new values only on the host.
3. Apply `server/assets/migrations/20260919_managed_assets.sql`.
4. Build and deploy backend, then Admin, then frontend with
   `ASSET_MANIFEST_ENABLED=false`.
5. Run the migration dry-run and apply flow in `asset-migration.md`.
6. In Admin, upload, preview, and publish `home.xray`, `resume.model`, and all
   eight sticker resources.
7. Verify `/api/front/assets?scope=home` and `scope=resume`, CORS, immutable
   caching, and the browser interactions.
8. Set `ASSET_MANIFEST_ENABLED=true` and restart only the backend container.
9. Observe upload errors, object 4xx/5xx, R2 requests, and traffic for at least
   24 hours before removing source fallbacks from future images.

Validation commands:

```bash
cd server && go test ./...
cd front && pnpm lint && pnpm build
cd admin && pnpm lint && pnpm build
docker compose -f deploy/start/docker-compose.yml config
```

Mobile acceptance must show no initial After or `.glb` request. Desktop must
show the static Before image first and load enhancements only after visibility
and idle scheduling. Test manifest 500, R2 404/CORS rejection, and WebGL failure
without losing readable content.

## Resource rollback

Open Admin → Resource Management and choose **Roll back to this version** on
one of the two archived versions. Confirm the old/new version IDs. This reuses
the publication transaction, changes only `active_version_id`, and does not
overwrite objects or purge Cloudflare.

## Application rollback

1. Set `ASSET_MANIFEST_ENABLED=false` if manifests or R2 delivery are faulty.
   The frontend keeps the small generated Before/model/sticker fallbacks.
2. Restore the previous backend/Admin/frontend image set.
3. Do not reverse the additive database migration and do not delete R2 objects.
4. Diagnose and republish or roll back the affected managed asset independently.

If a publication fails before commit, the previous active pointer remains
unchanged. If R2 is temporarily unavailable, Admin publication should stop;
public pages retain text and repository fallbacks.
