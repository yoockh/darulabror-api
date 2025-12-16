# darulabror-api

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](#)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)](#)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Darul Abror](https://img.shields.io/badge/Darul%20Abror-API-111827)](#)

Darul Abror API service (public + admin area) built with Echo, GORM, PostgreSQL, and optional Google Cloud Storage integration for public article media.

## Features

### Public
- List published articles
- Get published article by ID
- Student registration form
- Contact form

### Admin (`/admin`)
- Manage articles (CRUD)
- Manage registrations (list/detail/delete)
- Manage contacts (list/detail/update/delete)
- Profile endpoint (based on JWT claims)

### Superadmin (subset of `/admin`)
- Manage admins (create/list/update/delete)

## Requirements
- Go 1.22+
- PostgreSQL
- (Optional) Google Cloud Storage bucket + credentials

## Environment Variables

Required:
- `DATABASE_URL`  
  Example: `postgres://user:password@localhost:5432/darulabror?sslmode=disable`
- `JWT_SECRET`  
  Secret used to verify JWT tokens for `/admin` endpoints.

Optional:
- `PUBLIC_BUCKET`  
  If set, GCS client will be initialized and article media helpers can upload/resolve public URLs.

GCP auth (if using GCS):
- `GOOGLE_APPLICATION_CREDENTIALS` pointing to a service account JSON file.

## Run (local)

```bash
cp .env.example .env 2>/dev/null || true
go mod tidy
go run ./cmd/echo-server
```

Server defaults to `:8080` unless `PORT` is set.

## Routes (high level)

Public:
- `GET  /articles`
- `GET  /articles/:id`
- `POST /registrations`
- `POST /contacts`

Admin:
- `GET    /admin/profile`
- `GET    /admin/articles`
- `POST   /admin/articles`
- `PUT    /admin/articles/:id`
- `DELETE /admin/articles/:id`

- `GET    /admin/registrations`
- `GET    /admin/registrations/:id`
- `DELETE /admin/registrations/:id`

- `GET    /admin/contacts`
- `GET    /admin/contacts/:id`
- `PUT    /admin/contacts/:id`
- `DELETE /admin/contacts/:id`

Superadmin:
- `POST   /admin/admins`
- `GET    /admin/admins`
- `PUT    /admin/admins/:id`
- `DELETE /admin/admins/:id`

## Database

A PostgreSQL-compatible SQL schema is provided in `migrations/init.sql`.

## License
MIT. See [LICENSE](LICENSE).