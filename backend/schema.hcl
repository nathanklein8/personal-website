schema "public" {}

table "test_items" {
  schema = schema.public
  column "id" {
    type = serial
  }
  column "count" {
    type = int
    default = 0
  }
  column "note" {
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
