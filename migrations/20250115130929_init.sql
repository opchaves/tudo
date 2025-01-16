-- +goose Up
-- +goose StatementBegin
SET timezone TO 'GMT';

CREATE TABLE IF NOT EXISTS users(
  "id" INTEGER GENERATED ALWAYS AS IDENTITY (START WITH 10000),
  "uid" VARCHAR(10) NOT NULL,
  "first_name" VARCHAR NOT NULL,
  "last_name" VARCHAR,
  "email" VARCHAR NOT NULL,
  "password" VARCHAR,
  "avatar" VARCHAR,
  "verified" BOOLEAN NOT NULL DEFAULT false,
  "role" VARCHAR NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE,
  "last_active_at" TIMESTAMP WITHOUT TIME ZONE,
  CONSTRAINT "uq_users_uid" UNIQUE ("uid"),
  CONSTRAINT "uq_users_email" UNIQUE ("email"),
  CONSTRAINT "pk_users_id" PRIMARY KEY ("id")
);

-- user tokens (refresh, verification, invitation)
--   refresh: if deleted_at, no longer valid
--   verification: if deleted_at, already verified

CREATE TABLE IF NOT EXISTS user_tokens(
  "token" VARCHAR(21) NOT NULL,
  "tokey_type" VARCHAR NOT NULL,
  "user_id" INTEGER NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE,
  CONSTRAINT "pk_user_tokens_token" PRIMARY KEY ("token"),
  CONSTRAINT "fk_user_tokens_user_id" FOREIGN KEY ("user_id") REFERENCES users("id")
);

CREATE TABLE IF NOT EXISTS user_oauth_providers (
  "id" INTEGER GENERATED ALWAYS AS IDENTITY (START WITH 10000) PRIMARY KEY,
  "user_id" INTEGER NOT NULL,
  "provider" VARCHAR NOT NULL,
  "oauth_uid" VARCHAR NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "fk_user_oauth_providers_user_id" FOREIGN KEY ("user_id") REFERENCES users("id"),
  CONSTRAINT "uq_user_oauth_providers_user_provider" UNIQUE ("user_id", "provider")
);

CREATE TABLE IF NOT EXISTS notes (
  "id" INTEGER GENERATED ALWAYS AS IDENTITY (START WITH 10000),
  "uid" VARCHAR(21) NOT NULL,
  "title" VARCHAR NOT NULL,
  "content" TEXT NOT NULL,
  "user_id" INTEGER NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_notes_id" PRIMARY KEY ("id"),
  CONSTRAINT "uq_notes_uid" UNIQUE ("uid"),
  CONSTRAINT "fk_notes_user_id" FOREIGN KEY ("user_id") REFERENCES users("id")
);

CREATE TABLE IF NOT EXISTS workspaces (
  "id" INTEGER GENERATED ALWAYS AS IDENTITY (START WITH 10000),
  "uid" VARCHAR(10) NOT NULL,
  "name" VARCHAR NOT NULL,
  "owner_id" INTEGER NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_workspaces_id" PRIMARY KEY ("id"),
  CONSTRAINT "uq_workspaces_uid" UNIQUE ("uid"),
  CONSTRAINT "fk_workspaces_owner_id" FOREIGN KEY ("owner_id") REFERENCES users("id")
);

CREATE TABLE IF NOT EXISTS workspace_members (
  "workspace_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "permission" VARCHAR NOT NULL CHECK (permission IN ('viewer', 'editor', 'manager')),
  "active" BOOLEAN NOT NULL DEFAULT true,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_workspace_members" PRIMARY KEY ("workspace_id", "user_id"),
  CONSTRAINT "fk_workspace_members_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES workspaces("id"),
  CONSTRAINT "fk_workspace_members_user_id" FOREIGN KEY ("user_id") REFERENCES users("id")
);

CREATE TABLE IF NOT EXISTS invitations (
  "token" VARCHAR(21) NOT NULL,
  "permission" VARCHAR NOT NULL DEFAULT 'viewer' CHECK (permission IN ('viewer', 'editor', 'manager')),
  "workspace_id" INTEGER NOT NULL,
  "invited_user_id" INTEGER NOT NULL,
  "invited_by" INTEGER NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE,
  CONSTRAINT "pk_invitations_token" PRIMARY KEY ("token"),
  CONSTRAINT "fk_invitations_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES workspaces("id"),
  CONSTRAINT "fk_invitations_invited_user_id" FOREIGN KEY ("invited_user_id") REFERENCES users("id"),
  CONSTRAINT "fk_invitations_invited_by" FOREIGN KEY ("invited_by") REFERENCES users("id")
);

CREATE TABLE IF NOT EXISTS categories (
  "id" INTEGER GENERATED ALWAYS AS IDENTITY (START WITH 10000),
  "workspace_id" INTEGER NOT NULL,
  "name" VARCHAR NOT NULL,
  "description" VARCHAR,
  "type" VARCHAR NOT NULL CHECK (type IN ('income', 'expense')),
  "color" VARCHAR,
  "icon" VARCHAR,
  "created_by" INTEGER NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_categories_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES workspaces("id"),
  CONSTRAINT "fk_categories_created_by" FOREIGN KEY ("created_by") REFERENCES users("id")
);

CREATE TABLE IF NOT EXISTS accounts (
  "id" INTEGER GENERATED ALWAYS AS IDENTITY (START WITH 10000),
  "workspace_id" INTEGER NOT NULL,
  "name" VARCHAR NOT NULL,
  "balance" DECIMAL(15, 2) NOT NULL,
  "color" VARCHAR,
  "icon" VARCHAR,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_accounts_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_accounts_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES workspaces("id")
);

CREATE TABLE IF NOT EXISTS transactions (
  "id" INTEGER GENERATED ALWAYS AS IDENTITY,
  "uid" VARCHAR(21) NOT NULL,
  "workspace_id" INTEGER NOT NULL,
  "category_id" INTEGER NOT NULL,
  "account_id" INTEGER NOT NULL,
  "title" VARCHAR NOT NULL,
  "note" VARCHAR,
  "txn_date" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  "price" DECIMAL(15, 2) NOT NULL CHECK (price > 0),
  "type" VARCHAR NOT NULL CHECK (type IN ('income', 'expense', 'transfer')),
  "paid" BOOLEAN NOT NULL DEFAULT false,
  "created_by" INTEGER NOT NULL,
  "updated_by" INTEGER,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_transactions_id" PRIMARY KEY ("id"),
  CONSTRAINT "uq_transactions_uid" UNIQUE ("uid"),
  CONSTRAINT "fk_transactions_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES workspaces("id"),
  CONSTRAINT "fk_transactions_category_id" FOREIGN KEY ("category_id") REFERENCES categories("id"),
  CONSTRAINT "fk_transactions_account_id" FOREIGN KEY ("account_id") REFERENCES accounts("id"),
  CONSTRAINT "fk_transactions_created_by" FOREIGN KEY ("created_by") REFERENCES users("id"),
  CONSTRAINT "fk_transactions_updated_by" FOREIGN KEY ("updated_by") REFERENCES users("id"),
  CONSTRAINT "chk_transactions_paid_txn_date" CHECK (NOT paid OR txn_date <= now())
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS invitations;
DROP TABLE IF EXISTS workspace_members;
DROP TABLE IF EXISTS workspaces;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS user_oauth_providers;
DROP TABLE IF EXISTS user_tokens;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
