# Web frontend (demo)

## Run

From project root:

```bash
go run ./cmd/web
```

Open: [http://localhost:8080](http://localhost:8080)

## Implemented screens

- Main page with service presentation (`TurboLine AutoService`)
- User registration/login section
- Cars section (add/remove with VIN)
- Orders tab (parts search + cart + checkout)
- Order tracking section (poll every 10s)
- FAQ page section

## API contract expected by frontend

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/user/cars`
- `POST /api/user/cars`
- `DELETE /api/user/cars/{id}`
- `GET /api/parts/search?q=<query>`
  - response: `{ "items": [{ "part_id": "...", "name": "...", "brand": "...", "price": 1000, "delivery_day": 2 }] }`
- `POST /api/orders`
  - request: `{ "address": "...", "email": "...", "items": [...] }`
  - response: `{ "id": 123 }`
- `GET /api/orders/{id}/status`
  - response: `{ "status": "В обработке" }`

When endpoints are not ready, frontend uses demo fallback for search and status.
