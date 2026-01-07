-- Create "projects" table
CREATE TABLE "public"."projects" (
  "id" serial NOT NULL,
  "icon" text NOT NULL,
  "title" text NOT NULL,
  "description" text NOT NULL,
  "technologies" text NOT NULL,
  "deployment_link" text NULL,
  "image" text NULL,
  "alt_text" text NULL,
  PRIMARY KEY ("id")
);
