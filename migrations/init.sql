-- Table: admins
CREATE TABLE IF NOT EXISTS admins (
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin','superadmin')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Table: articles
CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    photo_header TEXT,
    content JSONB NOT NULL,
    author TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','published')),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Table: registrations
CREATE TABLE IF NOT EXISTS registrations (
    id BIGSERIAL PRIMARY KEY,

    student_type TEXT NOT NULL CHECK (student_type IN ('new','transfer')),
    gender TEXT NOT NULL CHECK (gender IN ('male','female')),

    email TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    phone TEXT NOT NULL,

    place_of_birth TEXT NOT NULL,
    date_of_birth DATE NOT NULL,

    address TEXT NOT NULL,
    origin_school TEXT NOT NULL,

    nisn TEXT NOT NULL UNIQUE,

    father_name TEXT NOT NULL,
    father_occupation TEXT NOT NULL,
    phone_father TEXT NOT NULL,
    date_of_birth_father DATE NOT NULL,

    mother_name TEXT NOT NULL,
    mother_occupation TEXT NOT NULL,
    phone_mother TEXT NOT NULL,
    date_of_birth_mother DATE NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table: contacts
CREATE TABLE IF NOT EXISTS contacts (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    subject TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at BIGINT NOT NULL
);
