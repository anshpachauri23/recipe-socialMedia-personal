# 🚀 Deployment Guide

Production architecture and the exact steps to stand it up from scratch.

## 📐 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  Users → Vercel (Next.js frontend, automatic from GitHub)   │
│            ↓ HTTPS                                          │
│         Northflank service (Go backend, Docker)             │
│            ↓ private network        ↓ HTTPS (IAM)           │
│         Northflank Postgres       AWS S3 (us-east-1)        │
│            (managed addon)        recipe-images/* prefix    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

| Layer | Provider | Notes |
|---|---|---|
| Frontend | Vercel | Auto-deploys on push to `main` |
| Backend | Northflank service | Built from [backend/Dockerfile](backend/Dockerfile), config in [backend/northflank.yaml](backend/northflank.yaml) |
| Database | Northflank Postgres addon | Same Northflank project as the backend → reachable over the internal private network |
| Image storage | AWS S3 | `us-east-1`, public-read on the `recipe-images/*` prefix only |

## 📋 Prerequisites

- Vercel account
- Northflank account
- AWS account (only used for S3)
- GitHub repository connected to both Vercel and Northflank

---

## 1. Frontend on Vercel

1. **Import the repo** at [vercel.com](https://vercel.com) → New Project → select this GitHub repo.
2. **Settings**: Framework = Next.js, Root Directory = `./`, Build Command = `npm run build`, Output Directory = `.next` (Vercel auto-detects all of these).
3. **Environment Variables** (Settings → Environment Variables):
   ```env
   NEXT_PUBLIC_API_URL=https://<your-backend>.northflank.app/api
   ```
4. Push to `main` → Vercel auto-deploys.

---

## 2. Backend on Northflank

### 2a. The Postgres addon

1. In your Northflank project → **Addons** → **Add addon** → **PostgreSQL**.
2. Pick a region close to where the backend service runs (same region = lowest latency).
3. After it provisions, open the addon → **Connection details** tab. You'll see two hosts: an **internal** one (e.g. `xxx.internal`) and a **public** one. The backend uses the internal host; the public one is for one-off `psql` from your laptop.
4. **Initialize the schema once**: Addon → Console / Query tab → paste the contents of [backend/database/schema.sql](backend/database/schema.sql) → Run. Or from your laptop:
   ```bash
   psql "postgresql://USER:PASS@PUBLIC_HOST:PORT/DBNAME?sslmode=require" \
        -f backend/database/schema.sql
   ```
   Verify with `\dt` — you should see `users`, `posts`, `post_images`, `follows`, `likes`, `comments`.

### 2b. The backend service

1. **Services** → **Create service** → **Combined service** (build + deploy from a Dockerfile).
2. **Source**: this GitHub repo, branch `main`, build context `backend/`, Dockerfile `Dockerfile`.
3. **Port**: `8080`, protocol HTTP (publicly exposed).
4. **Health check** (optional): GET `/api/posts/search?q=health` — already wired up in [backend/northflank.yaml](backend/northflank.yaml).
5. **Environment** → **Linked secrets** → link the Postgres addon and **rename** each variable to what the Go code reads ([backend/database/connection.go](backend/database/connection.go)):

   | Addon variable | Renamed to |
   |---|---|
   | `POSTGRES_HOST` (internal) | `DB_HOST` |
   | `POSTGRES_PORT` | `DB_PORT` |
   | `POSTGRES_USERNAME` | `DB_USER` |
   | `POSTGRES_PASSWORD` | `DB_PASSWORD` |
   | `POSTGRES_DATABASE` | `DB_NAME` |

6. **Environment** → **Plain-text variables / Secrets** → add the rest:
   ```env
   PORT=8080
   DB_SSLMODE=require
   JWT_SECRET=<long-random-string>
   S3_REGION=us-east-1
   S3_BUCKET=<your-bucket-name>
   AWS_ACCESS_KEY_ID=<from-step-3>
   AWS_SECRET_ACCESS_KEY=<from-step-3>   # mark as secret
   ```
7. **Deploy** → Northflank builds the Docker image and rolls it out.
8. **CORS for Vercel**: confirm the Go backend's CORS middleware allows the Vercel frontend origin.

---

## 3. AWS S3 (separate AWS account)

### 3a. Create the bucket

1. **AWS Console → S3 → Create bucket**.
2. Name: globally unique (e.g. `recipe-social-media-images-<suffix>`). Region: **us-east-1**.
3. **Object Ownership**: ACLs disabled (bucket-owner-enforced).
4. **Block Public Access**: **uncheck "Block all public access"** (objects need to be readable by `<img>` tags). Confirm the warning.
5. Versioning / encryption: defaults are fine.

### 3b. Bucket policy (public read on the recipe-images prefix only)

Permissions tab → Bucket policy → Edit:
```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Sid": "PublicReadRecipeImages",
    "Effect": "Allow",
    "Principal": "*",
    "Action": "s3:GetObject",
    "Resource": "arn:aws:s3:::YOUR_BUCKET_NAME/recipe-images/*"
  }]
}
```

### 3c. CORS

Permissions tab → CORS → Edit:
```json
[{
  "AllowedHeaders": ["*"],
  "AllowedMethods": ["GET"],
  "AllowedOrigins": ["*"],
  "ExposeHeaders": []
}]
```

### 3d. IAM user for the backend (least-privilege)

1. **IAM → Users → Create user** → name `recipe-social-backend-s3`. No console access.
2. **Attach policy** → inline policy:
   ```json
   {
     "Version": "2012-10-17",
     "Statement": [{
       "Effect": "Allow",
       "Action": ["s3:PutObject", "s3:GetObject"],
       "Resource": "arn:aws:s3:::YOUR_BUCKET_NAME/recipe-images/*"
     }]
   }
   ```
3. **Security credentials → Create access key** → "Application running outside AWS". Copy both keys (the secret is shown once). Paste them into the Northflank service env vars from step 2b.6.

The S3 object URL pattern produced by [backend/services/s3.go](backend/services/s3.go) is:
```
https://<bucket>.s3.us-east-1.amazonaws.com/recipe-images/<timestamp>-<filename>
```

---

## 4. Verify end-to-end

1. **Backend logs** in Northflank should show a clean startup (no `failed to ping database`).
2. `GET https://<backend>.northflank.app/api/posts/search?q=health` → 200 with JSON.
3. In the Postgres console: `SELECT COUNT(*) FROM users;` → `0`.
4. Sign up on the Vercel frontend → row appears in `users`.
5. Upload a profile photo → file appears under the bucket's `recipe-images/` prefix, the returned URL loads anonymously in a fresh browser tab, and `users.profile_photo_url` matches.
6. Create a post with images → repeat the same checks against `posts` / `post_images`.

---

## 🛠️ Common gotchas

| Symptom | Likely cause |
|---|---|
| Backend logs: `pq: SSL is not enabled on the server` or similar | `DB_SSLMODE` missing or set to `disable`. Northflank Postgres requires `require`. |
| Backend logs: `dial tcp: lookup ...: no such host` | Used the internal Postgres host but the service is in a different Northflank project. Either move it or use the public host (worse). |
| S3 bucket policy save fails with "Public access is blocked" | You forgot to uncheck "Block all public access" before applying the policy. |
| Image URLs return AccessDenied | IAM policy resource ARN is missing the `/recipe-images/*` suffix, OR the bucket policy is on a different prefix. |
| Frontend gets CORS errors against the backend | Backend CORS middleware doesn't include the Vercel origin. |
| Frontend gets 404s from the backend | `NEXT_PUBLIC_API_URL` missing the `/api` suffix or pointing at the wrong Northflank URL. |

---

## 🔁 Local development

For local dev, the backend can either point at the Northflank addon's **public** connection string (with `DB_SSLMODE=require`) or run a local Postgres in Docker. S3 vars can stay the same as production — uploads will just go straight to the prod bucket. Set everything in `backend/local.env` (gitignored) — see [backend/local.env.example](backend/local.env.example).

```bash
# Terminal 1
cd backend && go run main.go

# Terminal 2
npm run dev
```

Frontend on `http://localhost:3000`, backend on `http://localhost:8080`.
