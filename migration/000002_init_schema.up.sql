CREATE TABLE "notification_sent" (
  "id" bigserial PRIMARY KEY,
  "request_id" bigint NOT NULL,
  "driver_phone_rejected" varchar[] NOT NULL
);