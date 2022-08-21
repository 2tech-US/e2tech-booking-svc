CREATE TABLE "request" (
  "id" bigserial PRIMARY KEY,
  "type" varchar NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "pick_up_latitude" float8 NOT NULL,
  "pick_up_longitude" float8 NOT NULL,
  "drop_off_latitude" float8 NOT NULL,
  "drop_off_longitude" float8 NOT NULL,
  "status" varchar NOT NULL DEFAULT 'finding',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expire_at" timestamptz
);

CREATE TABLE "response" (
  "id" bigserial PRIMARY KEY,
  "request_id" bigint NOT NULL,
  "driver_phone" varchar NOT NULL DEFAULT 0,
  "driver_name" varchar NOT NULL,
  "driver_latitude" float8 NOT NULL,
  "driver_longitude" float8 NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "history" (
  "id" bigserial PRIMARY KEY,
  "type" varchar NOT NULL,
  "passenger_phone" varchar NOT NULL,
  "driver_phone" varchar NOT NULL,
  "pick_up_latitude" float8 NOT NULL,
  "pick_up_longitude" float8 NOT NULL,
  "drop_off_latitude" float8 NOT NULL,
  "drop_off_longitude" float8 NOT NULL,
  "price" float8 NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "done_at" timestamptz NOT NULL
);

CREATE INDEX ON "request" ("pick_up_longitude", "pick_up_latitude");

CREATE INDEX ON "request" ("drop_off_longitude", "drop_off_latitude");

CREATE INDEX ON "response" ("driver_longitude", "driver_latitude");

CREATE INDEX ON "history" ("type");

CREATE INDEX ON "history" ("passenger_phone");

CREATE INDEX ON "history" ("pick_up_longitude", "pick_up_latitude");

CREATE INDEX ON "history" ("drop_off_longitude", "drop_off_latitude");

COMMENT ON COLUMN "response"."driver_phone" IS '0 if found_driver is false';

COMMENT ON COLUMN "history"."passenger_phone" IS '0 if passenger is non-user';

ALTER TABLE "response" ADD FOREIGN KEY ("request_id") REFERENCES "request" ("id");
