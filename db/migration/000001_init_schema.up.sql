CREATE TABLE "user_policy" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "sender" varchar NOT NULL,
  "receiver" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "routine_transaction_policy" (
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "sender" varchar NOT NULL,
  "receiver" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "schedule_start_date" date NOT NULL,
  "schedule_end_date" date NOT NULL,
  "frequency" varchar NOT NULL,
  "amount" int NOT NULL
)INHERITS("user_policy");

CREATE TABLE "transaction_trigger_policy" (
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "sender" varchar NOT NULL,
  "receiver" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "targeted_balance" jsonb NOT NULL,
  "repeated" int NOT NULL DEFAULT 1,
  "amount" int NOT NULL
)INHERITS("user_policy");

CREATE TABLE "transaction" (
  "id" bigserial PRIMARY KEY,
  "nym_id" varchar NOT NULL,
  "transfer_sequence" int NOT NULL,
  "policy_id" bigserial
);

ALTER TABLE "transaction" ADD FOREIGN KEY ("policy_id") REFERENCES "user_policy" ("id");
