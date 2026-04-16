CREATE TABLE "urls" (
  "id" integer PRIMARY KEY,
  "short_code" varchar UNIQUE,
  "url" varchar,
  "createdAt" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updatedAt" timestamp DEFAULT CURRENT_TIMESTAMP,
  "accessCount" integer DEFAULT 0
);