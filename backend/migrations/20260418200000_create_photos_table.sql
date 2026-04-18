-- Create "photos" table
CREATE TABLE "public"."photos" (
  "id" serial NOT NULL,
  "title" text NOT NULL,
  "file_path" text NOT NULL,
  "alt_text" text NULL,
  "date_taken" text NULL,
  "location" text NULL,
  "camera" text NULL,
  "lens" text NULL,
  "aperture" text NULL,
  "shutter_speed" text NULL,
  "iso" text NULL,
  "visible" boolean NOT NULL DEFAULT true,
  "sort_order" integer NOT NULL DEFAULT 0,
  "source_path" text NOT NULL,
  "thumbnail_path" text NULL,
  "medium_path" text NULL,
  PRIMARY KEY ("id")
);
