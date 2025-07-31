# 📈 Stock Tracker - AI-Powered Investment Platform

A comprehensive **stock analysis** and **portfolio management** platform powered by **AI** and real-time market data. Make informed investment decisions with **advanced analytics** and **intelligent recommendations**.

---

## 🚀 Main Features

### 🤖 AI-Powered Stock Analysis

Get intelligent recommendations based on **real-time market data** and **AI algorithms**.

* Advanced recommendation engine
* Personalized suggestions by subscription level
* Machine learning-driven market analysis

### 📊 Real-Time Analytics

* **Portfolio Analytics**: Track and optimize investments with broker-based insights
* **Live Market Data**: Real-time prices, market trends, financial indicators (via Yahoo Finance)
* **Performance Tracking**: Monitor portfolio returns with detailed metrics

### 💬 AI Featured Chat (FinancIA)

* Recommend stocks aligned with market trends
* Provide company news and analysis
* Generate automated reports with scores and news
* Answer investment questions with AI insights

---

## 🏗️ Project Structure

```
Stock-Tracker/
├── api/          # Go backend API with CockroachDB
├── webui/        # Vue.js frontend application
├── infra/        # Infrastructure & deployment (Terraform, Docker)
├── docs/         # Documentation
├── scripts/      # Utility scripts
```

### Key Components

* **Backend**: Go + Chi Router, CockroachDB Cloud, JWT Auth, Yahoo Finance ingestion
* **Frontend**: Vue 3 + TypeScript + TailwindCSS + ECharts
* **Infrastructure**: Docker, Terraform (IaC), Redis & LocalStack for dev services

---

## 🛠️ Tech Stack

**Backend:** Go 1.24.4 • Chi Router • CockroachDB Cloud • JWT • Yahoo Finance API
**Frontend:** Vue 3 • Vite • Tailwind CSS • ECharts • Pinia • Axios
**Infra:** Docker & Compose • Terraform • Redis • LocalStack

---

## ⚙️ Quick Start

### Prerequisites

* Go **1.24.4+**
* Node.js **20.19.0+ or 22.12.0+**
* Docker & Docker Compose
* CockroachDB Cloud account (and SSL cert)
* PostgreSQL CLI (`psql`)

---

### 1️⃣ Clone & Setup

```bash
git clone https://github.com/JeySeb/Stock-Tracker
cd Stock-Tracker
```

### 2️⃣ Configure Environment

Create `api/.env` with (check `api/.env.example` for more details):

```bash
DATABASE_URL=postgresql://username:password@host:26257/db?sslmode=verify-full&sslrootcert=certs/cc-ca.crt
JWT_SECRET=your-jwt-secret
YAHOO_FINANCE_API_KEY=your-api-key
```

Frontend `webui/.env`:

```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=Stock Tracker
```

---

### 3️⃣ Setup Dev Environment

```bash
make api-setup      # Setup CockroachDB certs, dev services, and migrations
```

Start Redis & LocalStack:

```bash
make dev-up         # Start development services (Redis, LocalStack)
make dev-down       # Stop development services
make dev-logs       # View logs
```

---

## ▶️ Running the App

### Backend API

```bash
make backend-deps    # Install backend deps
make backend-run     # Run API locally (http://localhost:8080)
```

### Frontend WebUI

```bash
make frontend-deps   # Install frontend deps
make frontend-dev    # Start dev server (http://localhost:5173)
```

---

## 🗄️ Database Commands

```bash
make migrate-up         # Apply migrations
make migrate-down       # Rollback migrations
make migrate-reset      # Full reset (down+up)
make migrate-status     # Show current migration version
make migrate-specific MIGRATION=004 DIRECTION=up   # Run specific migration

make db-reset           # ⚠️ Full reset (destructive)
make db-shell           # Open CockroachDB shell
make db-test-connection # Test DB connection
```

---

## 🧪 Testing & Linting

### Backend

```bash
make backend-test           # Full suite
make backend-test-unit      # Unit tests only
make backend-test-api       # API & auth tests
make backend-test-coverage  # Coverage report (HTML in coverage/)
make backend-lint           # Lint Go code
make backend-test-security  # Security analysis (gosec)
```

### Frontend

```bash
make frontend-test          # Run frontend tests
make frontend-lint          # Lint Vue/TS code
```

---

## 🐳 Docker Commands

```bash
make docker-build-all       # Build all Docker images
make docker-up-full         # Run full stack (backend + frontend)
make docker-clean           # Clean Docker resources
```

---

## ☁️ Infrastructure (Terraform)

```bash
make infra-plan-local       # Plan local infra
make infra-apply-local      # Apply infra
make infra-destroy-local    # Destroy infra
```

---

## 📚 Documentation

* **[API Docs](api/README.md)** – Endpoints, migrations, DB config
* **[WebUI Docs](webui/README.md)** – Frontend architecture & components
* **[Infra Docs](infra/README.md)** – Terraform & deployment setup

---

## 🔧 Utilities

```bash
make help            # Show all available commands
make check-deps      # Verify all required tools
make status          # Show status of services & DB connection
make clean           # Remove build artifacts
```

---

## 🚀 Deployment

### Development:

```bash
make dev-up
```

---

**Built with ❤️ for intelligent investing**
