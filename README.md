# 🧠 AI-Powered Task Management System (Go + OpenAI)

An AI-driven task management backend that allows users to **create, update, delete, and search tasks using natural language**.  
Built with **Golang**, **PostgreSQL**, and **OpenAI Function Calling**, following a **Clean / Layered Architecture**.

---

## 🚀 Features

- 🗣️ **Natural Language Task Management**
  - “Create a high priority task to buy milk”
  - “Show me tasks due tomorrow”
- 🤖 **Autonomous AI Agent**
  - Uses OpenAI Function Calling (ReAct / Tool Use pattern)
- 🔐 **JWT-based Authentication**
- 🧱 **Clean Architecture** (Handlers → Services → Repositories)
- 📦 **RESTful API**
- ⚙️ **Graceful Shutdown & Connection Pooling**
- 📈 Designed for **scalability, observability, and cost control**

---

## 🏗️ Architecture Overview

The system follows a **Layered Architecture inspired by Clean Architecture**.

### High-Level Flow

```
Client → Middleware → Handler → Service → Repository → Database
                             ↓
                         OpenAI API
```

### Core Layers

- **Entry Point** (`cmd/server/main.go`)
  - Initializes DB, services, handlers, and HTTP server
- **Transport Layer** (`internal/handlers`)
  - Handles HTTP requests/responses
- **Middleware** (`internal/middlewares`)
  - JWT Authentication, CORS
- **Business Logic** (`internal/services`)
  - AuthService, AgentService
- **Data Access** (`internal/repository`)
  - PostgreSQL CRUD operations

---

## 🤖 AI Agent Design (ReAct / Tool Use)

The **AgentService** implements a controlled autonomous loop:

1. User sends a natural language command
2. Prompt + tool definitions are sent to OpenAI
3. OpenAI decides whether to:
   - Call a tool (e.g., `create_task`)
   - Or respond directly
4. Tool is executed via Go repository
5. Result is fed back to OpenAI
6. Loop continues (max 10 iterations)
7. Final natural language response is returned

### Supported Tools

- `create_task`
- `list_tasks`
- `update_task`
- `complete_task`
- `delete_task`
- `search_tasks`
- `calculate_date`

> ⚠️ The AI cannot execute arbitrary code — it can only invoke predefined tools.

---

## 📂 Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── middlewares/
│   ├── services/
│   ├── repository/
│   └── models/
├── db/
│   └── migrations/
├── config/
├── go.mod
├── go.sum
└── README.md
```

---

## 🔐 Authentication

- JWT Bearer Token authentication
- Token contains `userID`
- Middleware injects `userID` into request context
- All task operations are user-scoped

---

## ⚙️ Environment Variables

Create a `.env` file:

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/tasks
JWT_SECRET=your_long_random_secret
OPENAI_API_KEY=your_openai_api_key
```

> 🔒 Never commit `.env` or secrets to version control

---

## ▶️ Running the Project

### 1. Clone the repository

```bash
git clone https://github.com/your-username/task-manager-ai-backend.git
cd task-manager-ai-backend
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Start PostgreSQL

Ensure PostgreSQL is running and accessible via `DATABASE_URL`.

### 4. Run the server

```bash
go run cmd/server/main.go
```

Server starts on:

```
http://localhost:8080
```

---

## 🚀 Deployment

This project is deployed using **Render** for both the **Go backend** and **PostgreSQL database**.  
The setup is optimized for **learning, demos, and portfolio presentation**.

### Deployment Architecture

- **Backend:** Go (REST API)
- **Database:** PostgreSQL (Render Managed Database)
- **Hosting:** Render (Free Tier)

```
Client → Go Backend (Render) → PostgreSQL (Render)
```

---

### Backend Deployment (Render)

- A **Web Service** is created on Render and connected to the GitHub repository
- The backend is built and started using:

```bash
go build -o app
./app
```

- The application listens on the `PORT` environment variable provided by Render

---

### Database (Render PostgreSQL)

- Uses a **managed PostgreSQL instance** provided by Render
- Database credentials are injected using the `DATABASE_URL` environment variable
- The free-tier database is intended for **temporary and portfolio usage**

---

### Environment Variables (Render)

```env
DATABASE_URL=postgresql://<user>:<password>@<host>/<db>
JWT_SECRET=<your-secret>
PORT=8080
```

---

### Notes & Limitations

- Render free-tier services may experience **cold-start latency** after periods of inactivity
- The free PostgreSQL database **expires after a limited period** and may require recreation
- This deployment setup is **not intended for production workloads**

---

### Why Render?

- Zero cost
- Simple deployment workflow
- Suitable for demos and portfolio projects

---

## ⭐ Future Improvements

- Migrate PostgreSQL to a long-lived serverless provider (e.g., Neon)
- Reduce cold-start latency by moving the backend to a VM-based platform (e.g., Fly.io)

````

---

## 📡 API Example

### AI Task Command

```http
POST /agent/query
Authorization: Bearer <JWT_TOKEN>
````

```json
{
  "message": "Create a high priority task to buy groceries tomorrow"
}
```

---

## 🔒 Security Considerations

- OpenAI API key is never logged or exposed
- Prompt injection mitigated via controlled tools
- JWT tokens should be short-lived and signed with a strong secret

---

## 📈 Scalability & Observability

- Stateless backend (JWT-based)
- Horizontally scalable
- Token usage tracking per user (recommended)
- Structured logging with trace IDs

---

## 👨‍💻 Author

**Vinay**  
Backend Developer | Go | Distributed Systems | AI Integration
