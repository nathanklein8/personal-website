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
