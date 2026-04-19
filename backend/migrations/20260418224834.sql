-- Modify "photos" table
ALTER TABLE "public"."photos" ADD COLUMN "featured" boolean NOT NULL DEFAULT false;
