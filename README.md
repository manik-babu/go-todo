# Todo Server

A simple REST API for managing todos, built with Go and PostgreSQL.

## Requirements

- Go 1.26+
- PostgreSQL

## Setup

1. Clone the repository:
```bash
git clone https://github.com/manik-babu/go-todo.git
cd go-todo
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file with your database connection string:
```
DB_CONNECT=postgres://user:password@localhost:5432/tododb
```

4. Create the `todos` table in PostgreSQL:
```sql
CREATE TABLE todos (
  id SERIAL PRIMARY KEY,
  message TEXT NOT NULL,
  is_completed BOOLEAN DEFAULT FALSE
);
```

## Running

```bash
go run .
```

The server will start on `http://localhost:3000`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Health check |
| POST | `/todos` | Create a new todo |
| GET | `/todos` | Retrieve all todos |
| PATCH | `/todos/{id}` | Update a todo |
| DELETE | `/todos/{id}` | Delete a todo |

## Request/Response

**Create Todo:**
```json
POST /todos
{
  "message": "Buy groceries",
  "is_completed": false
}
```

**Response:**
```json
{
  "ok": true,
  "message": "Todo created successfully",
  "data": {
    "id": 1,
    "message": "Buy groceries",
    "is_completed": false
  }
}
```
