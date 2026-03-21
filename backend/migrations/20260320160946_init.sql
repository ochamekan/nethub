-- +goose Up
CREATE TABLE devices (
  id serial PRIMARY KEY,
  hostname text NOT NULL,
  location text NOT NULL,
  ip varchar(47) NOT NULL,
  is_active boolean NOT NULL DEFAULT false,
  is_deleted boolean NOT NULL DEFAULT false,
  created_at timestamp NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS devices;
