# Vodka

REST Framework written in Go.

Experimental and highly not recommended  to use.

Decided to migrate here: https://github.com/syndicatedb/vodka



## QueryBuilder

### Save method

Creating new row or executing ON CONFLICT statement with provided params.
It's possible to build queries with ON CONFLICT on unique keys or constraint and action UPDATE or NOTHING.

*PARAMS*

Params may be provided in query.

| name | type | comment | reqired |
|------|------|---------|---------|
| __conflictKey | string | name of constraint (ex. unique key) | false |
| __conflictAction | string | action to do on conflict (update / nothing) | true |


If you haven't got __conflictKey, unique fields must be set in model structure with tag `unique`. In this example `name` and `amount` are unique.

```Go
type Item struct {
	ID        string      `db:"id" json:"id"`
	Name      interface{} `db:"name" json:"name" unique:"true"`
	CreatedAt time.Time   `db:"created_at" json:"createdAt"`
	Amount    float64     `db:"amount" json:"amount" unique:"true"`
	Count     int64       `db:"count" json:"count"`
	Status    string      `db:"status" json:"status"`
}
```

If you haven't provided unique fields or constraint name, builder will generate query without ON CONFLICT statement.

*RESULT*

On creating / updating queries you will get a slice with affected row. If DO NOTHING action was set, response will be empty.

For example, for this request:

`items?__conflictKey=items_name_amount_pk&__conflictAction=update`

```json
{
	"name": "name name name",
	"count": 1,
	"amount": 1.3,
	"status": "active"
}
```

you will get such response:

```json
{
    "data": [
        {
            "id": "42cf70b7-0ad1-4ff3-4572-0615675f07a8",
            "name": "name name name",
            "createdAt": "2018-06-19T17:41:16.586148Z",
            "amount": 1.3,
            "count": 1,
            "status": "active"
        }
    ],
    "error": null
}
```