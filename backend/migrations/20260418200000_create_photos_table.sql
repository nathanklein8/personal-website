-- Add photo library paths and thumbnail columns to photos table
ALTER TABLE "public"."photos"
  ADD COLUMN "source_path" text NOT NULL DEFAULT '',
  ADD COLUMN "thumbnail_path" text,
  ADD COLUMN "medium_path" text;
