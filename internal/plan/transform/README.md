
A transform directive instructs doppelganger to alter the value of this column as per the expression provided.
An example coud be:
```hcl
    postgres "production" {
        url = "postgres://user:${password}@localhost:5432/db"
        table "users" {
            query = "SELECT id, name, surname FROM users"
            transform "name" "scramble()"
            transform "surname" "this is a constant value"
        }
    }
```
See the `transform` package to see a list of available functions and how to create yours using a plugin.