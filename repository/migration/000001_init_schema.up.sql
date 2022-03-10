CREATE TABLE "user_policies" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "description" VARCHAR NOT NULL,
  "nym_id" VARCHAR NOT NULL,
  "recipient" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE "routine_transaction_policies" (
  "schedule_start_date" DATE NOT NULL,
  "schedule_end_date" DATE NOT NULL,
  "frequency" VARCHAR NOT NULL,
  "amount" JSONB NOT NULL
)
INHERITS (
  "user_policies"
);

CREATE TABLE "transaction_trigger_policies" (
  "targeted_balance" JSONB NOT NULL,
  "amount" JSONB NOT NULL
)
INHERITS (
  "user_policies"
);

CREATE TABLE "transactions" (
  "id" BIGSERIAL PRIMARY KEY,
  "nym_id" VARCHAR NOT NULL,
  "transfer_sequence" INT UNIQUE NOT NULL,
  "transfer" JSONB UNIQUE NOT NULL,
  "policy_id" BIGINT NOT NULL
);

ALTER TABLE "transactions"
  ADD FOREIGN KEY ("policy_id") REFERENCES "user_policies" ("id");
