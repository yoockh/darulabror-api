# darulabror-api

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](#)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)](#)
[![Swagger](https://img.shields.io/badge/Swagger-API%20Docs-85EA2D?logo=swagger&logoColor=000)](https://darulabror-717070183986.asia-southeast2.run.app/swagger/index.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Darul Abror](https://img.shields.io/badge/Darul%20Abror-API-111827)](#)

Darul Abror API service (public + admin area) built with Echo, GORM, PostgreSQL, and optional Google Cloud Storage integration for public article media.

## API Documentation (Swagger)

- Swagger UI: https://darulabror-717070183986.asia-southeast2.run.app/swagger/index.html
- OpenAPI JSON: https://darulabror-717070183986.asia-southeast2.run.app/swagger/doc.json
- OpenAPI YAML: https://darulabror-717070183986.asia-southeast2.run.app/swagger/doc.yaml

## Response format (convention)

Most JSON responses follow this envelope (see `internal/utils/response.go`):

Success:
```json
{
  "status": "success",
  "message": "OK",
  "data": {}
}
```

Error:
```json
{
  "status": "error",
  "message": "something went wrong"
}
```

Notes:
- Some endpoints intentionally return **No Content** (`201/204` with empty body) because handlers use `c.NoContent(...)`.

## Authentication (Admin)

1) Login to get JWT:
- `POST /admin/login`

2) Use token for admin endpoints:
- Header: `Authorization: Bearer <token>`

Role rules (see `api/routes/routes.go`):
- `/admin/*` requires role: `admin` or `superadmin`
- `/admin/admins*` requires role: `superadmin`

## Pagination

List endpoints support:
- `page` (default 1)
- `limit` (default 10, max 100)

Response uses:
```json
{
  "items": [],
  "meta": { "page": 1, "limit": 10, "total": 123 }
}
```

---

## Articles: `content` is flexible JSON (frontend-defined)

In the database/model, `Article.content` is stored as **JSONB** (`gorm.io/datatypes.JSON`).
That means the backend does **not** enforce a fixed schema for article body content.

**Frontend decides the shape**, e.g. a block editor style:
```json
{
  "blocks": [
    { "type": "heading", "level": 2, "text": "Judul" },
    { "type": "paragraph", "text": "Teks panjang..." },
    { "type": "image", "url": "https://storage.googleapis.com/<bucket>/articles/xxx.jpg", "caption": "..." },
    { "type": "video", "url": "https://storage.googleapis.com/<bucket>/articles/yyy.mp4" }
  ]
}
```

Important:
- File binary (image/video) **is not stored inside `content`**.
- Store only **URLs** (or object names) returned by the media upload endpoint.

### `photo_header`
`photo_header` is intended for the article card/cover image (thumbnail/banner). It should be a URL string.

---

## Media Upload (images/videos)

### POST /admin/articles/media (Admin)
Uploads a file to storage and returns a URL (public bucket) that you can put into:
- `photo_header`
- `content` JSON (e.g., image/video blocks)

Swagger: see `Articles (Admin) -> POST /admin/articles/media`

Request:
- `multipart/form-data`
- field name: `file`

Response `201` (example):
```json
{
  "status": "success",
  "message": "media uploaded",
  "data": {
    "url": "https://storage.googleapis.com/<bucket>/articles/12345_file.jpg"
  }
}
```

Requirements:
- Set `PUBLIC_BUCKET` env var in the server runtime to enable GCS uploads.
- If `PUBLIC_BUCKET` is empty / GCS not configured, upload will fail.

---

## Health & Swagger

### GET /healthz
Response:
- `200 OK` body: `ok`

### GET /swagger/index.html
Swagger UI

---

## Public Endpoints

### GET /articles
Query:
- `page` (optional)
- `limit` (optional)

Response `200`:
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
- `200` success envelope with `data` = article
- `400` if id invalid
- `404` if not found / not published

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

### POST /contacts
Request body:
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

## Admin Endpoints (requires BearerAuth)

### GET /admin/profile
Response `200`:
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

### Articles (Admin)
- `GET    /admin/articles` (list draft + published)
- `POST   /admin/articles` (create)
- `PUT    /admin/articles/:id` (update)
- `DELETE /admin/articles/:id` (delete)
- `POST   /admin/articles/media` (upload image/video for `photo_header` or `content`)

Create/Update request body: `dto.ArticleDTO` (see `internal/dto/article_dto.go`)

Create response:
- `201 Created` (no body)

Update response:
- `200 OK` (no body)

Delete response:
- `204 No Content` (no body)

### Registrations (Admin)
- `GET    /admin/registrations` (list)
- `GET    /admin/registrations/:id` (detail)
- `DELETE /admin/registrations/:id` (delete)

### Contacts (Admin)
- `GET    /admin/contacts` (list)
- `GET    /admin/contacts/:id` (detail)
- `PUT    /admin/contacts/:id` (update)
- `DELETE /admin/contacts/:id` (delete)

---

## Superadmin Endpoints (requires role=superadmin)

### Admin management
- `POST   /admin/admins`
- `GET    /admin/admins`
- `PUT    /admin/admins/:id`
- `DELETE /admin/admins/:id`

Create/Update body: `dto.AdminDTO` (see `internal/dto/admin_dto.go`)
Notes:
- Create requires `password` non-empty (see `internal/handler/admin_handler.go`).

## Database

A PostgreSQL-compatible SQL schema is provided in `migrations/init.sql`.

## License
MIT. See [LICENSE](LICENSE).