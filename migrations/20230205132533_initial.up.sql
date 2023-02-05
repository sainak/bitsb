CREATE TABLE "user"
(
    "id"         BIGINT       NOT NULL,
    "email"      VARCHAR(255) NOT NULL,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name"  VARCHAR(255) NOT NULL,
    "type"       INTEGER      NOT NULL,
    "password"   VARCHAR(255) NOT NULL,
    "last_login" TIMESTAMPTZ  NULL,
    "created_at" TIMESTAMPTZ  NOT NULL,
    "updated_at" TIMESTAMPTZ  NOT NULL
);