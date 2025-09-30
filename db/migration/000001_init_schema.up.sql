CREATE TABLE "users" (
    "id" uuid PRIMARY KEY,
    "username" varchar UNIQUE NOT NULL,
    "name" varchar NOT NULL,
    "dni" varchar UNIQUE NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "wallets" (
    "id" uuid PRIMARY KEY,
    "user_id" uuid NOT NULL,
    "currency" varchar NOT NULL,
    "balance" decimal NOT NULL DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "wallets" ADD CONSTRAINT "fk_wallets_users" FOREIGN KEY ("user_id") REFERENCES "users"("id");