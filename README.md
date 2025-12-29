# darulabror-api

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](#)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)](#)
[![Swagger](https://img.shields.io/badge/Swagger-API%20Docs-85EA2D?logo=swagger&logoColor=000)](https://darulabror-717070183986.asia-southeast2.run.app/swagger/index.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Darul Abror backend API (public + admin) built with **Go (Echo)**, **GORM**, **PostgreSQL**, and **Google Cloud Storage (GCS)** for article media uploads.

---

## Documentation (Swagger)

- Swagger UI: https://darulabror-717070183986.asia-southeast2.run.app/swagger/index.html
- OpenAPI JSON: https://darulabror-717070183986.asia-southeast2.run.app/swagger/doc.json
- OpenAPI YAML: https://darulabror-717070183986.asia-southeast2.run.app/swagger/doc.yaml

---

## Features

### Public
- List published articles (pagination)
- Get published article detail
- Create registration
- Create contact message

### Admin (JWT)
- Login + profile
- Manage articles (CRUD)
  - Create/Update uses **multipart/form-data**
  - `photo_header` is **required**
  - Inline image/video for `content` supported via **single request** (placeholders + multipart files)
- Manage registrations (list/detail/delete)
- Manage contacts (list/detail/update/delete)

### Superadmin (JWT + role)
- Manage admins (create/list/update/delete)

---

## Project Structure

```text
.
├── api/
│   ├── middleware/          # JWT auth + role guards
│   └── routes/              # HTTP route registration
├── cmd/
│   └── echo-server/         # Server entrypoint (main.go)
├── config/                  # DB configuration
├── docs/                    # Generated Swagger docs (swag)
├── internal/
│   ├── dto/                 # DTOs for requests/responses
│   ├── handler/             # HTTP handlers + Swagger annotations
│   ├── models/              # GORM models
│   ├── repository/          # Data access layer (Postgres + GCS)
│   ├── service/             # Business logic
│   └── utils/               # Response helpers, pagination, auth context, etc.
└── migrations/
    └── init.sql             # SQL schema bootstrap
```

---

## Environment Variables

Required:
- `DATABASE_URL` — PostgreSQL DSN/URL
- `JWT_SECRET` — JWT HMAC secret
- `CORS_ORIGINS` — comma-separated allowlist (required)

Optional:
- `PUBLIC_BUCKET` — enables GCS uploads for article media
- `PORT` — default `8080`
- `ALLOW_LOCALHOST_CORS` — set to `true` to allow `http://localhost:3000` and `http://127.0.0.1:3000` for local development (default: `false`)


---

## Response Format (Convention)

Most JSON responses use this envelope (see `internal/utils/response.go`):

### Success
```json
{
  "status": "success",
  "message": "OK",
  "data": {}
}
```

### Error
```json
{
  "status": "error",
  "message": "something went wrong"
}
```

Notes:
- Some endpoints intentionally return **No Content** (`201/204` with empty body) because handlers use `c.NoContent(...)`.

---

## Authentication (Admin)

### POST /admin/login
Request (JSON):
```json
{
  "email": "admin@darulabror.com",
  "password": "StrongPassword123"
}
```

Response `200`:
```json
{
  "status": "success",
  "message": "login success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "admin": {
      "id": 1,
      "username": "admin",
      "email": "admin@darulabror.com",
      "role": "admin",
      "is_active": true,
      "created_at": 1734567890,
      "updated_at": 1734567890
    }
  }
}
```

Use token for admin endpoints:
- Header: `Authorization: Bearer <token>`

Role rules:
- `/admin/*` requires role: `admin` or `superadmin`
- `/admin/admins*` requires role: `superadmin`

---

## Pagination

List endpoints support:
- `page` (default 1)
- `limit` (default 10, max 100)

Response:
```json
{
  "status": "success",
  "message": "OK",
  "data": {
    "items": [],
    "meta": { "page": 1, "limit": 10, "total": 123 }
  }
}
```

---

## Articles: Core Rules

### 1) `photo_header` is REQUIRED
`photo_header` is the cover image used by article cards in list view.  
**Every article must have it** on create/update.

Provide it by:
- uploading a file via `photo_header_file` (recommended), OR
- sending a URL via `photo_header`

### 2) `content` is FLEXIBLE JSON
`content` is stored as JSONB (no fixed schema enforced by backend). Frontend can use canvas/block-editor style.

### 3) Inline image/video inside `content` (single request)
Binary files (image/video) cannot be embedded directly into JSON. This API supports single-request upload:

- Put placeholders inside `content` using `"upload_key": "<key>"`
- Attach the actual files in the same multipart request using:
  - `content_files[<key>]`

Server behavior:
- uploads `content_files[...]` to GCS
- replaces `upload_key` with a `url` field inside `content` before saving

Example content sent by frontend:
```json
{
  "blocks": [
    { "type": "paragraph", "text": "Hello" },
    { "type": "image", "upload_key": "img1", "caption": "Image in body" },
    { "type": "video", "upload_key": "vid1" }
  ]
}
```

Stored content will contain:
- `{ "type":"image", "url":"https://storage.googleapis.com/<bucket>/articles/content/img1_....jpg", "caption":"..." }`
- `{ "type":"video", "url":"https://storage.googleapis.com/<bucket>/articles/content/vid1_....mp4" }`

Requirement:
- `PUBLIC_BUCKET` must be configured, otherwise uploads will fail.

---

## Endpoints (Detailed)

### Health & Docs

#### GET /healthz
Response:
- `200 OK` (plain text): `ok`

#### GET /swagger/index.html
Response:
- `200 OK` (Swagger UI)

---

## Public Endpoints

### GET /articles
Query:
- `page` (optional)
- `limit` (optional)

Response `200` example:
```json
{
  "status": "success",
  "message": "articles fetched",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "Example",
        "photo_header": "https://storage.googleapis.com/<bucket>/articles/header.jpg",
        "content": {},
        "author": "Admin",
        "status": "published",
        "created_at": 1734567890,
        "updated_at": 1734567890
      }
    ],
    "meta": { "page": 1, "limit": 10, "total": 1 }
  }
}
```

### GET /articles/:id
Response:
- `200 OK` → full article object
- `400` → invalid `id`
- `404` → not found / not published

---

### POST /registrations
Request body: `dto.RegistrationDTO` (see `internal/dto/registration_dto.go`)

Example request:
```json
{
  "student_type": "new",
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone": "081234567890",
  "gender": "male",
  "place_of_birth": "Bandung",
  "date_of_birth": "2007-01-02",
  "address": "Jl. Contoh No. 1",
  "origin_school": "SMP Contoh",
  "nisn": "1234567890",
  "father_name": "Father",
  "father_occupation": "Employee",
  "phone_father": "081234567890",
  "date_of_birth_father": "1980-01-02",
  "mother_name": "Mother",
  "mother_occupation": "Homemaker",
  "phone_mother": "081234567891",
  "date_of_birth_mother": "1982-01-02"
}
```

Response:
- `201 Created` (no body)

---

### POST /contacts
Request:
```json
{
  "email": "user@example.com",
  "subject": "Question",
  "message": "Hello..."
}
```

Response:
- `201 Created` (no body)

---

## Admin Endpoints (JWT required)

### GET /admin/profile
Response `200` example:
```json
{
  "status": "success",
  "message": "profile fetched",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@darulabror.com",
    "role": "admin",
    "is_active": true,
    "created_at": 1734567890,
    "updated_at": 1734567890
  }
}
```

---

## Articles (Admin) — Multipart (Create/Update)

### POST /admin/articles
Request: **multipart/form-data**

Required fields:
- `title` (string)
- `author` (string)
- `content` (string) → must be **valid JSON string**
- `photo_header_file` (file) **OR** `photo_header` (string URL) → **required one of them**

Optional fields:
- `status` (`draft|published`)

Optional inline media fields (repeatable):
- `content_files[img1]` (file)
- `content_files[vid1]` (file)
- `content_files[any_key]` (file)

Response:
- `201 Created` (no body)

Example (curl):
```bash
curl -i -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -F 'title=My Article' \
  -F 'author=Admin' \
  -F 'status=published' \
  -F 'content={"blocks":[{"type":"paragraph","text":"Hello"},{"type":"image","upload_key":"img1"},{"type":"video","upload_key":"vid1"}]}' \
  -F 'photo_header_file=@/path/to/header.jpg' \
  -F 'content_files[img1]=@/path/to/body-image.jpg' \
  -F 'content_files[vid1]=@/path/to/body-video.mp4' \
  https://darulabror-717070183986.asia-southeast2.run.app/admin/articles
```

### PUT /admin/articles/:id
Same fields/behavior as create.

Response:
- `200 OK` (no body)

### DELETE /admin/articles/:id
Response:
- `204 No Content`

---

## Registrations (Admin)
- `GET /admin/registrations` (list)
- `GET /admin/registrations/:id` (detail)
- `DELETE /admin/registrations/:id` (delete)

---

## Contacts (Admin)
- `GET /admin/contacts` (list)
- `GET /admin/contacts/:id` (detail)
- `PUT /admin/contacts/:id` (update)
- `DELETE /admin/contacts/:id` (delete)

---

## Superadmin (role=superadmin)
- `POST /admin/admins`
- `GET /admin/admins`
- `PUT /admin/admins/:id`
- `DELETE /admin/admins/:id`

---

## Local Development

```bash
go mod download
go run ./cmd/echo-server
```

### Enabling CORS for Local Frontend

When developing locally with a frontend running on `http://localhost:3000`, you need to enable CORS for localhost:

1. Set the required `CORS_ORIGINS` environment variable with your production origins
2. Set `ALLOW_LOCALHOST_CORS=true` to allow localhost access

Example:
```bash
export CORS_ORIGINS="https://www.darulabror.com,https://admin.darulabror.com"
export ALLOW_LOCALHOST_CORS=true
export DATABASE_URL="postgres://user:pass@localhost/darulabror"
export JWT_SECRET="your-secret"
go run ./cmd/echo-server
```

This will allow requests from:
- `http://localhost:3000`
- `http://127.0.0.1:3000`
- All origins in `CORS_ORIGINS`

**Note**: In production, never set `ALLOW_LOCALHOST_CORS=true`. The server will only accept origins explicitly listed in `CORS_ORIGINS`.

Regenerate Swagger docs:
```bash
swag init -g cmd/echo-server/main.go -o docs --parseDependency --parseInternal
```

---

## Database

A PostgreSQL-compatible SQL schema is provided in `migrations/init.sql`.

---

## License

MIT. See [LICENSE](LICENSE).