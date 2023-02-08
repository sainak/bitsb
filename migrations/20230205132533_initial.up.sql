CREATE TABLE "users"
(
    "id"           SERIAL PRIMARY KEY,
    "email"        VARCHAR(255) UNIQUE NOT NULL,
    "first_name"   VARCHAR(255)        NOT NULL,
    "last_name"    VARCHAR(255)        NOT NULL,
    "access_level" INTEGER             NOT NULL,
    "password"     VARCHAR(255)        NOT NULL,
    "last_login"   TIMESTAMPTZ         NULL,
    "created_at"   TIMESTAMPTZ         NOT NULL,
    "updated_at"   TIMESTAMPTZ         NOT NULL
);