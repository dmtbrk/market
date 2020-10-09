# Market

This is a product management service. It is using [AIexMoran/httpCRUD](https://github.com/AIexMoran/httpCRUD) service
for authorization.

## How to run

`Docker` and `docker-compose` must be installed on the system to be able to run this project.

- `make` builds and runs the project in Docker.

- `make run` creates and runs Docker containers of this project and `AIexMoran/httpCRUD`.

- `make build` builds or rebuilds services.

## Usage

### HTTP

`GET /products/?offset=0&limit=10` lists all products in the given range. Offset and limit parameters required.

Response example:
```
200 OK
```
```
[
    {
        "id": "1",
        "name": "Banana",
        "price": 1500,
        "seller": "1234"
    },
    {
        "id": "2",
        "name": "Carrot",
        "price": 1400,
        "seller": "bunny"
    }
]
```

`GET /products/{id}` shows product details by the specified id.

Response example:
```
200 OK
```
```
{
    "id": "1",
    "name": "Banana",
    "price": 1500,
    "seller": "1234"
}
```

`POST /products/` adds a product to the product list. Authorization required.

Request example:
```
{
    "name": "Banana",
    "price": 1500,
    "seller": "1234"
}
```
Response example:
```
200 OK
```
```
{
    "id": "1",
    "name": "Banana",
    "price": 1500,
    "seller": "1234"
}
```

`PATCH /products/{id}` updates product data by the specified id. Update only happens for provided fields. Authorization required.

Request example:
```
{
    "name": "Banana-nana-nana... Batman",
}
```

Response example:
```
200 OK
```
```
{
    "id": "1",
    "name": "Banana-nana-nana... Batman",
    "price": 1500,
    "seller": "1234"
}
```

`DELETE /products/{id}` removes the product by the specified id. Authorization required. 

Response example:
```
200 OK
```
```
{
    "id": "1",
    "name": "Banana",
    "price": 1500,
    "seller": "1234"
}
```

#### Authorization

Requests to protected resources are expected to have an `Authorization` header with a token issued by `AIexMoran/httpCRUD`.
This service starts alongside with the product service (`make` or `make run` does it), so additional actions
not required. 
The usage may be found [here](https://github.com/AIexMoran/httpCRUD).

Example:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```