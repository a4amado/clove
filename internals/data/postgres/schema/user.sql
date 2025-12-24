CREATE TABLE "user" (
    "id" UUID PRIMARY KEY DEFAULT(uuid_generate_v4()),
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "hash" TEXT NOT NULL,
    "createdAt" TIMESTAMP DEFAULT NOW(),
    "updatedAt" TIMESTAMP DEFAULT NOW()
);
 