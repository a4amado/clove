CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE credential_type AS ENUM ('session', 'sdk', 'ott');
CREATE TYPE region AS ENUM ('dk1');
CREATE TYPE app_type AS ENUM ('free', 'standard', 'pro');
CREATE TYPE user_role AS ENUM ('super', 'admin', 'user');
 
CREATE TABLE "user" (
    "id" UUID PRIMARY KEY DEFAULT(uuid_generate_v4()),
    "hash" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW(),
    "role" user_role DEFAULT('user'),
    "number_of_emails" INT NOT NULL DEFAULT 0 CHECK ("number_of_emails" >= 0 AND "number_of_emails" <= 10)
);
create index user_msg_idx on "user"("id");
CREATE TABLE "app" (
    "id" UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "app_slug" VARCHAR(255) NOT NULL UNIQUE,
    "region" region[] NOT NULL,
    "app_type" app_type NOT NULL,
    "user_id" UUID NOT NULL,
    "allowed_origins" VARCHAR(255)[],
    
    constraint dk_userId
        FOREIGN KEY ("user_id")
        References "user"("id")
        ON DELETE CASCADE     
);

CREATE TABLE "user_email" (
    "id" UUID PRIMARY KEY DEFAULT(uuid_generate_v4()),
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW(),
    "verified_at"  TIMESTAMP,
    "code" VARCHAR(90) NOT NULL,
    CONSTRAINT user_email_user_id_fk
        FOREIGN KEY ("user_id")
        REFERENCES "user"("id")
        ON DELETE CASCADE


);
create index user_email_email_idx on "user_email"("email");
create index user_email_user_id_idx on "user_email"("user_id");

 

create index app_id_idx on "app"("id");
create index app_slug_idx on "app"("app_slug");

CREATE TABLE "app_api_key" (
    "id" VARCHAR(64) PRIMARY KEY, -- hased version of the actual key
    "app_id" UUID,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW(),
    "name" VARCHAR(50),
    "prefix" VARCHAR(5),
    "suffix" VARCHAR(5),
    "expires_at" TIMESTAMP NOT NULL,
    CONSTRAINT "api_key_app_fk" FOREIGN KEY ("app_id") REFERENCES "app"("id")
);
create INDEX "app_api_key_appId_idx"  on "app_api_key"("app_id");

CREATE TABLE "credential" (
    "id"          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "token"       TEXT NOT NULL UNIQUE,
    "user_id"     UUID,
    "app_id"      UUID,
    "channel_id"  TEXT,
    "type"        credential_type NOT NULL,
    "permissions" TEXT NOT NULL,
    "expires_at"  TIMESTAMP NOT NULL,
    "created_at"  TIMESTAMP DEFAULT NOW(),
    CONSTRAINT "credential_user_fk"
        FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE,
    CONSTRAINT "credential_app_fk"
        FOREIGN KEY ("app_id") REFERENCES "app"("id") ON DELETE CASCADE
);
CREATE INDEX credential_token_idx ON "credential"("token");

