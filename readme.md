# Industrial Time-Series Prediction Demo

A small Go-based microservice demo that simulates reservoir pressure readings, stores them in Postgres, and serves both raw history and short-term Fourier-based forecasts over an HTTP+Plotly UI.

> **Not for production** — this is a toy/demo app to show how you can stitch together data generation, server/API, and in-browser visualization with Go, Postgres, Docker, Gin and Plotly.

---

## 🚀 Features

- **Synthetic data generator** (`cmd/generator`):  
  - Emits a new “pressure” reading every 0.5 s (simulates 30 min of process time per tick)  
  - Daily sinusoidal pattern + random noise
- **API server** (`cmd/server`):  
  - `GET /api/readings` → all historical readings  
  - `GET /api/predict?horizon=N` → next N-step forecast via least-squares + Fourier features  
  - Serves a single‐page HTML/JS UI (`internal/web/index.html`)
- **Plotly.js UI**:  
  - Live‐updating chart of raw, smoothed, and forecasted data  
  - Slider control for forecast horizon  
  - Auto‐refresh every 0.5 s

---

## 📁 Directory Layout

```

go-industrial-ts/
├── cmd/
│   ├── generator/      # data generator binary
│   │   └── main.go
│   └── server/         # API + UI server binary
│       └── main.go
├── internal/
│   ├── data/           # Postgres I/O + Fourier features + OLS
│   │   ├── db.go
│   │   ├── fourier.go
│   │   └── model.go
│   └── web/            # static HTML/JS Plotly client
│       └── index.html
├── init\_db.sql         # creates `pressure` table
├── Dockerfile          # multi-stage build for generator & server
└── docker-compose.yml  # brings up Postgres, generator & server

````

---

## 🛠️ Prerequisites

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- Go 1.22+ (for local builds inside Docker)

---

## ⚡ Quick Start with Docker

1. **Clone & enter**  
   ```bash
   git clone <repo-url>
   cd go-industrial-ts
````

2. **Build & launch**

   ```bash
   docker-compose up --build
   ```

   * **db** (`postgres:13`) runs on `localhost:5432` and applies `init_db.sql`
   * **generator** begins inserting data every 0.5 s
   * **server** listens on `localhost:8080`

3. **Browse**
   Open [http://localhost:8080](http://localhost:8080) in your browser.

4. **Stop**

   ```bash
   docker-compose down
   ```

---

## 🔌 Configuration

* **Postgres DSN**: via `PG_CONN` environment variable (both generator & server).
* **Generator pace**: in `cmd/generator/main.go`, `time.NewTicker(500 * time.Millisecond)`.
* **Fourier model**:

  * Period (default 100) and harmonics (default 3) in `cmd/server/main.go`.
  * Intercept + harmonics defined in `internal/data/fourier.go`.

---

## 📊 API Endpoints

| Method | Path                     | Description                                                          |
| ------ | ------------------------ | -------------------------------------------------------------------- |
| GET    | `/api/readings`          | Returns JSON array of all `{id, pressure}` readings                  |
| GET    | `/api/predict?horizon=N` | Returns `{ ids: [...], pred: [...] }` for the next N forecast points |

---

## 📊 UI Controls

* **Slider** at the top to choose forecast horizon (1–2000 samples).
* **Auto-refresh** every 0.5 s.
* **Traces**:

  * **Blue**: raw data (`id`, `pressure`)
  * **Black dashed**: 3-point moving average
  * **Green**: forecast for the next *N* readings

---

## 🧪 Extending & Testing

* Tweak generator’s noise or cycle in `cmd/generator/main.go` → rebuild.
* Increase harmonics or add a trend term in `internal/data/fourier.go`.
* Swap Plotly for another charting library in `internal/web/index.html`.