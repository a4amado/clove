CREATE TABLE "user" (
    "id" UUID PRIMARY KEY DEFAULT(uuid_generate_v4()),
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "hash" TEXT NOT NULL,
    "createdAt" TIMESTAMP DEFAULT NOW(),
    "updatedAt" TIMESTAMP DEFAULT NOW()
);

CREATE TYPE changeOperation AS ENUM ('INSERT', 'DELETE', 'UPDATE');

CREATE TABLE "userHistory" (
    "id" UUID PRIMARY KEY DEFAULT(uuid_generate_v4()),
    "userId" UUID NOT NULL,
    "email" VARCHAR(255) NOT NULL,  -- STILL WRONG - REMOVE UNIQUE
    "hash" TEXT NOT NULL,  -- YOU'RE STILL MISSING THIS
    "createdAt" TIMESTAMP NOT NULL,
    "updatedAt" TIMESTAMP NOT NULL,
    "changedAt" TIMESTAMP NOT NULL DEFAULT NOW(),
    "changedBy" UUID,  -- SHOULD BE NULLABLE
    "operation" changeOperation NOT NULL,
    
    CONSTRAINT fk_user_id
        FOREIGN KEY ("userId")
        REFERENCES "user"("id")
        ON DELETE CASCADE
);
