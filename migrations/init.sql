-- Table: admins
CREATE TABLE admins (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL, -- superadmin / admin 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: articles
CREATE TABLE articles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title TEXT NOT NULL,
    content JSONB NOT NULL,
    author TEXT NOT NULL,
    status TEXT DEFAULT 'draft', -- draft / published
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: registrations
CREATE TABLE registrations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    student_type ENUM('NEW','TRANSFER') NOT NULL,
    gender ENUM('MALE','FEMALE') NOT NULL,

    email VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,

    place_of_birth VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,

    address TEXT NOT NULL,
    origin_school VARCHAR(150) NOT NULL,

    nisn VARCHAR(20) NOT NULL UNIQUE,

    father_name VARCHAR(100) NOT NULL,
    father_occupation VARCHAR(100) NOT NULL,
    phone_father VARCHAR(20) NOT NULL,
    date_of_birth_father DATE NOT NULL,

    mother_name VARCHAR(100) NOT NULL,
    mother_occupation VARCHAR(100) NOT NULL,
    phone_mother VARCHAR(20) NOT NULL,
    date_of_birth_mother DATE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
