CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE region AS ENUM ('dk1');
CREATE TYPE app_type AS ENUM ('free', 'standard', 'pro');
CREATE TYPE UserRole AS ENUM ('super', 'admin', 'user');

CREATE TABLE "user" (
    "id" UUID PRIMARY KEY DEFAULT(uuid_generate_v4()),
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "hash" TEXT NOT NULL,
    "createdAt" TIMESTAMP DEFAULT NOW(),
    "updatedAt" TIMESTAMP DEFAULT NOW(),
    "role" UserRole DEFAULT('user')
);
create index user_msg_idx on "user"("id");
create index user_email_idx on "user"("email");


CREATE TABLE "app" (
    "id" UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "appSlug" VARCHAR(255) NOT NULL UNIQUE,
    "region" region[] NOT NULL,
    "appType" app_type NOT NULL,
    "userId" UUID NOT NULL,
    "allowedOrigins" VARCHAR(255)[],
    
    constraint dk_userId
        FOREIGN KEY ("userId")
        References "user"("id")
        ON DELETE CASCADE     
);

create index app_id_idx on "app"("id");
create index app_slug_idx on "app"("appSlug");

CREATE TABLE "app_api_key" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "appId" UUID,
    "createdAt" TIMESTAMP DEFAULT NOW(),
    "updatedAt" TIMESTAMP DEFAULT NOW(),
    "key" TEXT,
    "name" VARCHAR(50),
    CONSTRAINT "api_key_app_fk" FOREIGN KEY ("appId") REFERENCES "app"("id")
);
create INDEX "app_api_key_appId_idx"  on "app_api_key"("appId")