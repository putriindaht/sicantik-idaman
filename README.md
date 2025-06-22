# Sicantik Idaman (Sistem Pencatatan Cuti Karawan Terintegrasi, Mudah, dan Aman)

**Sicantik Idaman** is a digital leave management system designed to simplify the process of leave request, approval, and tracking for employees and organizations.

---

## 🚀 Features

- 🔐 **JWT Authentication**
- 🔄 **Role-Based Access** (Employee, Manager, HR, Director)
- 📌 **Leave Request Lifecycle**:
  - Submit, update, delete requests
  - Manager/Director approvals
- 📆 **Leave Balance Tracking**
- 📊 **View Approved Leaves**
- 💬 **Leave Reactions** (like/dislike, one per user per leave)

---

## 🛠 Tech Stack

| Component     | Description                                  |
|---------------|----------------------------------------------|
| Go            | Backend programming language (Golang)        |
| Gin           | Web framework for routing                    |
| GORM          | ORM for PostgreSQL                           |
| PostgreSQL    | Primary database                             |
| JWT           | JSON Web Tokens for authentication           |
| Bcrypt        | For secure password hashing                  |
| Docker (opt.) | Containerization for local setup             |

---

## 📦 Installation & Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/sicantik-idaman.git
   cd sicantik-idaman
   ```

2. **Configure Environment**

   * Create a `.env` file based on `.env.example`
   * Fill in your database config:

     ```env
        APP_PORT=1010
        JWT_SECRET=yourjwtsecret
        DSN_DB=your dsndb
        POSTGRES_USER=postgres
        POSTGRES_PASSWORD=postgresoryourdbpassword
        POSTGRES_DB=sicantik-db
        DB_HOST=localhost
        DB_PORT=5432
        DB_USER=postgres
        DB_PASSWORD=postgresoryourdbpassword
        DB_NAME=sicantik-db
     ```

3. **Run the Application**

   ```bash
   go run cmd/sicantik-idaman/main.go
   ```

   This will:

   * Load config
   * Connect to PostgreSQL
   * Run auto migration & seeding
   * Start the server on the defined port

---

## 🔐 Roles and Permissions

| Role     | Capabilities                                             |
| -------- | -------------------------------------------------------- |
| Employee | Submit and manage their own leave requests               |
| Manager  | Approve leave requests from their team, view team leaves |
| HR       | Get notified of approvals across teams   |
| Director | Final approver, auto-approved requests, can approve all  |

---

## 📡 API Endpoints

### ➤ Base

* `GET /` — Test API

### ➤ Auth

* `POST /api/v1/login` — Login with email and password

### ➤ Leave Types

* `GET /api/v1/leaves/types` — Get all available leave types

### ➤ Leave Requests

* `POST /api/v1/leaves/requests` — Submit a leave request
* `GET /api/v1/leaves/requests/me` — View own leave requests
* `GET /api/v1/leaves/requests/approved` — View approved leaves (can filter with `start` and `end` query)
* `PATCH /api/v1/leaves/requests/:id` — Approve or reject leave (Manager/Director only)
* `PUT /api/v1/leaves/requests/:id` — Edit own pending request
* `DELETE /api/v1/leaves/requests/:id` — Soft delete own request
* `GET /api/v1/leaves/requests/:id/reactions` — View reactions for a leave request

### ➤ Leave Balance

* `GET /api/v1/leaves/balances/me` — View own leave balance

### ➤ Leave Reactions

* `POST /api/v1/leaves/reactions` — Add reaction to approved leave
* `PATCH /api/v1/leaves/reactions/:id` — Update your reaction
* `DELETE /api/v1/leaves/reactions/:id` — Soft delete your reaction

---

## 📌 Notes

* All routes (except `/login`) require JWT authentication.
* Users can only manage their own data unless authorized by role.
* Reactions are allowed only for approved leave requests.
* One reaction per user per leave.

---

## 👩‍💻 Author

**Putri Indah**