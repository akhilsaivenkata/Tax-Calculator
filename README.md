# Tax Calculator API using golang

This is a backend API built using Go and Gin for calculating income tax based on marginal tax brackets. The app fetches tax slab data from an external Flask-based API (provided in a Docker container), applies tax calculations, and returns the result via a RESTful endpoint.

---

## 🔧 Tech Stack

- **Go** (Gin framework)
- **Docker** & **Docker Compose**
- **Logrus** for structured logging
- **Testify** for unit test assertions
- **Flask** (Dockerized external tax API)

---

## 🚀 Run the App Locally (Without Docker)

### 1. Clone the repo

```bash
git clone https://github.com/akhilsaivenkata/Tax-Calculator.git
cd Tax-Calculator
```

### 2. Run the Flask Tax API (Dockerized)

```bash
docker pull ptsdocker16/interview-test-server
docker run -p 5001:5001 ptsdocker16/interview-test-server
```

The tax API will now be available at `http://localhost:5001`

```bash
export TAX_API_URL=http://localhost:5001
```

### 4. Run the Go server

```bash
go run ./cmd/server
```

The Go server will now be running at `http://localhost:8080`

---

## 🐳 Run the App with Docker Compose (Recommended)

This will spin up both services (Go + Flask) together in a Docker network.

### 1. Run the full stack

```bash
docker-compose up --build
```

- Go API: [http://localhost:8080](http://localhost:8080)
- Flask API: [http://localhost:5001](http://localhost:5001)

### 2. Stop the stack

```bash
Ctrl + C
```

---

## 📬 Example API Usage

### `POST /calculate-tax`

**Request:**

```json
{
  "income": 75000,
  "tax_year": 2022
}
```

**Response:**

```json
{
  "total_tax": 10000,
  "effective_tax": 0.1333,
  "breakdown": [
    {
      "min": 0,
      "max": 50000,
      "rate": 0.1,
      "tax_paid": 5000
    },
    {
      "min": 50000,
      "rate": 0.2,
      "tax_paid": 5000
    }
  ]
}
```

---

## 🧪 Running Tests

### Run all unit tests with coverage:

```bash
go test ./... -cover
```

### Generate coverage report files:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

> `coverage.html` will show you line-by-line coverage visually.

---


## 🧠 Notes

- The external tax API (Flask) randomly fails and only supports years 2019–2022. This is intentional to test error handling.
- The Go app validates all inputs and logs structured JSON logs (good for observability in production).
- The service layer is fully decoupled and unit-tested.

---


