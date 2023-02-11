CREATE TABLE locations
(
    id         SERIAL PRIMARY KEY  NOT NULL,
    name       VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ         NOT NULL,
    updated_at TIMESTAMPTZ         NOT NULL
);
CREATE INDEX idx_locations_created_at ON locations (created_at);


CREATE TABLE "users"
(
    "id"               SERIAL PRIMARY KEY                NOT NULL,
    "email"            VARCHAR(255) UNIQUE               NOT NULL,
    "first_name"       VARCHAR(255)                      NOT NULL,
    "last_name"        VARCHAR(255)                      NOT NULL,
    "access_level"     INTEGER                           NOT NULL,
    "password"         VARCHAR(255)                      NOT NULL,
    "last_login"       TIMESTAMPTZ                       NULL,
    "created_at"       TIMESTAMPTZ                       NOT NULL,
    "updated_at"       TIMESTAMPTZ                       NOT NULL,
    "home_location_id" INTEGER REFERENCES locations (id) NULL,
    "work_location_id" INTEGER REFERENCES locations (id) NULL
);
CREATE INDEX idx_users_email ON users (email);


CREATE TABLE bus_routes
(
    id           SERIAL PRIMARY KEY NOT NULL,
    name         VARCHAR(100)       NOT NULL,
    number       VARCHAR(5) UNIQUE  NOT NULL,
    start_time   TIME               NOT NULL,
    end_time     TIME               NOT NULL,
    interval     INTEGER            NOT NULL,
    location_ids INTEGER[]          NOT NULL,
    min_price    INTEGER DEFAULT 5  NOT NULL,
    max_price    INTEGER DEFAULT 25 NOT NULL,
    created_at   TIMESTAMP          NOT NULL,
    updated_at   TIMESTAMP          NOT NULL
);
CREATE INDEX idx_bus_routes_created_at ON bus_routes (created_at);
CREATE INDEX idx_bus_route_locations ON bus_routes USING GIN (location_ids);
