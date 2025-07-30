# PianPianino

**_PianPianino_** is a personal To-Do manager designed to _lazily_ help you manage your time. The name is a pun on a saying in the Florentine dialect of italian: "*vai pianino*" means "take it slow", a reminder to think before acting. <br>
The app is built using the following technologies:
- Backend:
    - The Go Programming language.
    - The Bun ORM.
- Frontend:
    - Vue.js.
    - Naive ui as a component library.
    - Axios for http requests.

![Demo](./assets/PianPianino.gif)

The project also has a test suite for the backend, to run those use:
```bash
cd backend/
go test ./tests/*
```

## Features
- Web based interface to manage your To-Dos.
- REST API with secure JWT authentication.
- User registration and login.
- Track your To-Dos based on their priority and organize them by date.

## How to run
First of all make sure you have installed the Go programming language.
- Clone the repository  ```git clone https://github.com/efive-dev/PianPianino.git```.
- Navigate to the backend directory and run ```go mod tidy```.
- Navigate to the frontend directory and run ```npm run install```.
- Creat a *.env* file in the top directory with at least the two following variables:
    - DATABASE_DSN=./../db.sqlite
    - JWT_SECRET=your_secret_key
- Navigate to the backend directory and run the local server ```go run ./cmd/```.
- Navigate to the frontend directory and run the local server ```npm run dev```.

The backend REST API will be available at the following address: `http://localhost/1323`. <br>
The frontend will be available at `http://localhost/5173`.

## Database Schema
### Tasks:
| Column        | Type      | Constraints                                                                 |
| ------------- | --------- | --------------------------------------------------------------------------- |
| `id`          | INTEGER   | Primary Key, Auto-increment                                                 |
| `user_id`     | INTEGER   | Not Null, Foreign Key → `users(id)`, On Delete: Cascade, On Update: Cascade |
| `description` | TEXT      | —                                                                           |
| `importance`  | INTEGER   | Not Null, Default: 0 (`CHECK importance IN (0,1,2,3)`)                      |
| `completed`   | BOOLEAN   | Not Null, Default: false                                                    |
| `created_at`  | TIMESTAMP | Not Null, Default: CURRENT\_TIMESTAMP                                       |
| `updated_at`  | TIMESTAMP | Not Null, Default: CURRENT\_TIMESTAMP                                       |
**Notes**:
- importance maps to the Importance enum with values:
  - 0 = NotSet
  - 1 = Low
  - 2 = Medium
  - 3 = High
- User is a related model, referenced via user_id using a belongs-to relationship.
- created_at and updated_at use nullzero, meaning they omit zero values in inserts/updates, but are always non-null by schema.

### Users:
| Column     | Type    | Constraints                 |
| ---------- | ------- | --------------------------- |
| `id`       | INTEGER | Primary Key, Auto-increment |
| `username` | TEXT    | Not Null, Unique            |
| `password` | TEXT    | Not Null                    |
**Notes**:
- Tasks represents a one-to-many relationship with the Task model (has-many), joined by users.id = tasks.user_id.
- username is unique and required.
- password is required (stored as text but hashed beforehand).

## Endpoints

The API is organized into:

- **Public routes:** Available without authentication (e.g., user registration and login).
- **Protected routes:** Require a valid **JWT** to access (e.g., task management).

Authentication is handled using JWT (JSON Web Tokens). To access protected routes, clients must include a valid token in the `Authorization` header using the `Bearer <token>` format.

| Method | Path                   | Description                   | Auth Required |
|--------|------------------------|-------------------------------|--------------|
| POST   | `/register`            | Register a new user            | No           |
| POST   | `/login`               | Log in a user                  | No           |
| GET    | `/api/tasks`           | List all tasks                 | Yes          |
| POST   | `/api/tasks`           | Create a new task              | Yes          |
| DELETE | `/api/tasks/:id`       | Delete a task by ID            | Yes          |
| PATCH  | `/api/tasks/:id/toggle`| Toggle task completion status | Yes          |

---

**Notes**:
- Protected routes are all prefixed with /api.
- JWT authentication middleware is applied on the /api group.
- CORS is configured to allow requests from `http://localhost:5173` and `http://localhost:1323/`.


