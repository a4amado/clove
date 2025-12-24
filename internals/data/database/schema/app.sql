CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE region AS ENUM ('dk1');
CREATE TYPE app_type AS ENUM ('free', 'standard', 'pro');

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
create index app_slug_idx on "app"("app_slug");
