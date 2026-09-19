# Asset migration

Back up MySQL before any apply run. The migration tool is dry-run by default,
streams input through a bounded temporary file, verifies the uploaded object
with `Head`, and only then updates a database reference. Failed or unreadable
legacy URLs remain unchanged and appear in the JSON report.

## Inventory and dry-run

Run from `server/` with production environment variables loaded:

```bash
go run ./tools/assetctl \
  -mode migrate \
  -config config.docker.yml \
  -roots ../front/public,public/uploaded \
  -report /tmp/asset-migration-dry-run.json
```

The inventory covers repository/static files plus known URL-bearing fields in
`article`, `config`, `page`, `user_info`, and `resume_profile`. Search the report
for `failed`, confirm every legacy domain/path was found, and retain the report
with the deployment record.

## Apply

After verifying the database backup and R2 configuration:

```bash
go run ./tools/assetctl \
  -mode migrate \
  -config config.docker.yml \
  -apply \
  -backup-confirmed \
  -report /var/backups/gvb/asset-migration-apply.json
```

The destination key is content-addressed. Repeated source URLs are uploaded
once per run, and database references are changed only after upload and `Head`
verification. The tool never deletes a legacy source.

For managed assets, use Admin → Resource Management to upload and publish the
Before/After pair, the original GLB plus the script-generated delivery GLB, and
eight stickers. When the original GLB is already at most 6 MiB, the separate
delivery part is optional. The pair is one atomic version. Keep
`ASSET_MANIFEST_ENABLED=false` until both public manifest scopes return the
intended published versions.

## Reproducible delivery optimization

```bash
cd scripts
npm ci
npm run optimize:images
npm run optimize:model -- ../front/public/resume/resume-ready.glb ../front/public/resume/resume-ready.optimized.glb
```

The GLB command pins glTF Transform 4.5.0, keeps flatten/join/palette/simplify
disabled, and enables Meshopt plus WebP textures. The backend independently
checks the GLB container, required camera/animation/eyes/sticker nodes, and the
6 MiB hard limit before accepting it.

## Garbage collection

Versions beyond the current plus two historical releases first become
`purge_pending`. Objects are eligible only after seven days. Always inspect the
dry-run before applying:

```bash
go run ./tools/assetctl -mode gc -config config.docker.yml -report /tmp/asset-gc.json
go run ./tools/assetctl -mode gc -config config.docker.yml -apply -report /var/backups/gvb/asset-gc-apply.json
```

The GC command re-checks the active pointer and the three-version protection
set. A failed object deletion keeps the version row so the operation can be
retried.
