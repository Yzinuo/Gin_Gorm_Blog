# Cloudflare R2 setup

This runbook keeps all R2 credentials on the server. Never put the API token,
access key, or secret key in Git, frontend environment files, screenshots, issue
text, or database rows.

## 1. Rotate exposed credentials

If a credential has ever been pasted into chat or another shared system, create
a replacement bucket-scoped token and revoke the old token before production
deployment. Grant only object read/write/list access to `my-blog`.

The Cloudflare API token shown when the key is created is not used by the S3
client. The backend needs the S3 Access Key ID and Secret Access Key.

## 2. Bucket and public domain

1. Keep `my-blog` on R2 Standard storage.
2. In R2 → `my-blog` → Settings → Custom Domains, bind
   `assets.heliar.top`.
3. Wait for the custom domain to show Active and verify DNS/TLS.
4. Disable the public `r2.dev` URL after the custom domain works.
5. Do not create a public browser PUT policy; every upload goes through Gin.

Recommended bucket CORS policy (replace the origins with the exact production
frontend and Admin origins; add localhost only in development):

```json
[
  {
    "AllowedOrigins": [
      "https://heliar.top",
      "https://www.heliar.top"
    ],
    "AllowedMethods": ["GET", "HEAD"],
    "AllowedHeaders": ["Range", "If-None-Match"],
    "ExposeHeaders": [
      "Accept-Ranges",
      "Content-Length",
      "Content-Range",
      "ETag"
    ],
    "MaxAgeSeconds": 86400
  }
]
```

Create a Cloudflare Cache Rule for `assets.heliar.top/*` that permits caching
of versioned objects, including `.glb`. Objects are uploaded with
`Cache-Control: public, max-age=31536000, immutable`; never purge or overwrite
an existing key.

## 3. Server secrets

Copy `.env.example` to a server-only environment file outside the repository,
set mode `0600`, and populate these values:

```dotenv
STORAGE_PROVIDER=r2
R2_ACCOUNT_ID=<account-id>
R2_ACCESS_KEY_ID=<bucket-scoped-access-key>
R2_SECRET_ACCESS_KEY=<bucket-scoped-secret>
R2_BUCKET=my-blog
R2_PUBLIC_BASE_URL=https://assets.heliar.top
ASSET_MANIFEST_ENABLED=false
CORS_ALLOWED_ORIGINS=https://heliar.top,https://www.heliar.top
```

`ASSET_MANIFEST_ENABLED` stays false until the migration, initial version
publication, and smoke tests have passed. The backend fails at startup when R2
is selected and any required R2 variable is missing. Startup logs must never
print either key.

## 4. Verification

After deployment and after publishing one test asset:

```bash
curl -fsSI https://assets.heliar.top/<versioned-object-key>
curl -fsS 'https://heliar.top/api/front/assets?scope=home'
curl -fsS 'https://heliar.top/api/front/assets?scope=resume'
```

Check for a one-year immutable cache header, a stable ETag, a successful CORS
request from each allowed origin, and a Cloudflare cache hit on a repeated GET.
Configure R2 usage notifications and a budget alert in Cloudflare. The Admin
resource page also raises a visible warning at 8 GiB; neither mechanism is a
hard spending cap.
