CREATE TABLE "urls" (
  "id" integer PRIMARY KEY,
  "short_code" varchar,
  "url" varchar,
  "createdAt" timestamp,
  "updatedAt" timestamp,
  "accessCount" integer
);