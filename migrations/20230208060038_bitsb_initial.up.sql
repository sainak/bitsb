CREATE TABLE locations
(
    id         SERIAL PRIMARY KEY  NOT NULL,
    name       VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ         NOT NULL,
    updated_at TIMESTAMPTZ         NOT NULL
);

CREATE INDEX idx_locations_created_at ON locations (created_at);

CREATE TABLE companies
(
    id          SERIAL PRIMARY KEY  NOT NULL,
    name        VARCHAR(255) UNIQUE NOT NULL,
    location_id INTEGER             REFERENCES locations (id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ         NOT NULL,
    updated_at  TIMESTAMPTZ         NOT NULL
);

CREATE INDEX idx_companies_created_at ON companies (created_at);
CREATE INDEX idx_companies_location_id ON companies (location_id);