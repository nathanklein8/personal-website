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
