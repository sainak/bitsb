CREATE TABLE bus_routes
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    number       VARCHAR(5)   NOT NULL UNIQUE,
    start_time   TIME         NOT NULL,
    end_time     TIME         NOT NULL,
    interval     INTEGER      NOT NULL,
    location_ids INTEGER[]    NOT NULL,
    created_at   TIMESTAMP    NOT NULL,
    updated_at   TIMESTAMP    NOT NULL
);

CREATE INDEX idx_bus_route_locations ON bus_routes USING GIN (location_ids);
