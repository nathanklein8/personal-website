-- Create "test_items" table
CREATE TABLE "public"."test_items" (
  "id" serial NOT NULL,
  "count" integer NOT NULL DEFAULT 0,
  "note" text NOT NULL,
  PRIMARY KEY ("id")
);
