-- Create "landing_card" table
CREATE TABLE "public"."landing_card" (
  "id" serial NOT NULL,
  "bio" text NOT NULL,
  "email" text NOT NULL,
  "linkedin" text NOT NULL,
  "skills" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- Drop "test_items" table
DROP TABLE "public"."test_items";
