schema "public" {}

table "landing_card" {
  schema = schema.public
  column "id" {
    type = serial
  }
  column "bio" {
    type = text
  }
  column "email" {
    type = text
  }
  column "linkedin" {
    type = text
  }
  column "github" {
    type = text
  }
  column "skills" {
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
}

table "projects" {
  schema = schema.public
  column "id" {
    type = serial
  }
  column "icon" {
    type = text
    null = false
  }
  column "title" {
    type = text
    null = false
  }
  column "description" {
    type = text
    null = false
  }
  column "technologies" {
    type = text
    null = false
    // store technologies as json text ex: ["a", "b"]
  }
  column "deployment_link" {
    type = text
    null = true
  }
  column "image" {
    type = text
    null = true
  }
  column "alt_text" {
    type = text
    null = true
  }
  primary_key {
    columns = [column.id]
  }
}

table "photos" {
  schema = schema.public
  column "id" {
    type = serial
  }
  column "title" {
    type = text
    null = false
  }
  column "file_path" {
    type = text
    null = false
    // relative path within the volume mount, e.g. "landscapes/sunset.jpg"
  }
  column "alt_text" {
    type = text
    null = true
  }
  column "date_taken" {
    type = text
    null = true
  }
  column "location" {
    type = text
    null = true
  }
  column "camera" {
    type = text
    null = true
  }
  column "lens" {
    type = text
    null = true
  }
  column "aperture" {
    type = text
    null = true
  }
  column "shutter_speed" {
    type = text
    null = true
  }
  column "iso" {
    type = text
    null = true
  }
  column "visible" {
    type = boolean
    null = false
    default = true
  }
  column "sort_order" {
    type = integer
    null = false
    default = 0
  }
  primary_key {
    columns = [column.id]
  }
}
