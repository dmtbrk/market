# Market

This service lets users manage their products. User management is provided by [AIexMoran/httpCRUD](https://github.com/AIexMoran/httpCRUD) service.

## How to run

`Docker` and `docker-compose` must be installed on the system to be able to run this project.
- `make` builds and runs the project in Docker.

- `make up` creates and runs Docker containers of this project and `AIexMoran/httpCRUD` a path to which must be specified in the `docker-compose.yml` file.

- `make down` stops and removes containers, networks, images, and volumes.

- `make build` builds or rebuilds services.

- `make start` starts services.

- `make stop` stops services.

## Usage

`GET /products/?offset=0&limit=10` lists all products in the given range. Offset and limit parameters required.

Response example:
```
Status: 200 OK
```
```
[
    {
        "id": 1,
        "name": "Banana",
        "price": 1500,
        "seller": "1234"
    },
    {
        "id": 2,
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
    "id": 1,
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
    "id": 1,
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
    "id": 1,
    "name": "Banana-nana-nana... Batman",
    "price": 1500,
    "seller": "1234"
}
```

`DELETE /products/{id}` removes the product by the specified id. Authorization required. 

Response example:
```
203 No Content
```

### Authorization

The request is expected to have an `Authorization` header with the token issued by `AIexMoran/httpCRUD`. The usage may be found [here](https://github.com/AIexMoran/httpCRUD).

Example:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```