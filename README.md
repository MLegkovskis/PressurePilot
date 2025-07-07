# Real-Time Time-Series Forecasting & Visualization Platform

[](https://go.dev/)
[](https://www.docker.com/)
[](https://opensource.org/licenses/MIT)

A multi-container Go application demonstrating a full-cycle, real-time data pipeline. This project simulates industrial sensor data, streams it into a PostgreSQL database, and serves it via a REST API to a live-updating web UI that includes a predictive forecast using Fourier analysis.

This isn't just a simple web app; it's a showcase of microservice architecture, containerization, data persistence, and CI/CD-ready design patterns.

\<br\>

> *Live UI showing historical data (blue), a smoothed moving average (black), and a real-time forecast (orange).*

-----

## 🏛️ System Architecture

This project follows a classic microservice pattern, with each component isolated in its own Docker container and managed by Docker Compose. This design ensures separation of concerns, scalability, and a production-ready deployment model.

```mermaid
graph TD
    subgraph "Docker Environment"
        subgraph "Service: server"
            A[▶️ Gin API Server <br> (cmd/server)] -- ":8080" --> U[👨‍💻 User's Browser]
        end

        subgraph "Service: generator"
            G[📈 Data Generator <br> (cmd/generator)]
        end

        subgraph "Service: db"
            DB[(📦 PostgreSQL <br> postgres:13)]
        end

        subgraph "Network: default"
            A -- "Reads/Writes" --> DB
            G -- "Writes" --> DB
        end
    end

    U -- "HTTP Requests" --> A
    A -- "GET /api/readings" --> D1[API Response: Historical Data]
    A -- "GET /api/predict" --> D2[API Response: Forecast Data]
    U --> D1
    U --> D2

    style U fill:#f9f,stroke:#333,stroke-width:2px
    style DB fill:#bbf,stroke:#333,stroke-width:2px
```

1.  **`db` (PostgreSQL)**: The stateful component. A standard `postgres:13` image that persists all time-series data. The schema is initialized on first run using the `init_db.sql` volume mount.
2.  **`generator` (Go)**: A headless microservice whose only job is to simulate sensor readings. It connects to the database and inserts a new data point every 500ms, creating a continuous data stream.
3.  **`server` (Go & Gin)**: The user-facing microservice.
      * It serves the static `index.html` frontend, which contains the Plotly.js visualization logic.
      * It exposes a REST API with two endpoints (`/api/readings` and `/api/predict`).
      * When a request for a prediction arrives, it fetches all historical data, computes Fourier features, fits a linear regression model, and returns the forecast.
4.  **Networking**: All services are connected on a shared Docker network. The `server` and `generator` use the `db` service's DNS name (`db`) to connect to the PostgreSQL instance. Only the `server`'s port `8080` and the `db`'s port `5432` are exposed to the host.

-----

## ✨ Key Features & Concepts Demonstrated

This project showcases skills across the development and operations lifecycle:

  * **Containerization & Orchestration**: Multi-stage `Dockerfile` creates lean, optimized runtime images from a single build definition. `docker-compose.yml` orchestrates the entire application stack for one-command local deployment.
  * **Microservice Architecture**: Clear separation of concerns between the data generator, the API server, and the database.
  * **CI/CD Friendly**: The entire application is built and run in a containerized environment, making it portable and easy to integrate into any automated CI/CD pipeline (e.g., GitHub Actions, GitLab CI).
  * **Data Modeling & Persistence**: A relational schema (`init_db.sql`) stores time-series data, and Go services interact with it safely using connection pooling (`pgxpool`).
  * **Backend Development (Go)**:
      * **API Serving**: A robust REST API built with the high-performance **Gin** web framework.
      * **Database Interaction**: Clean, efficient SQL queries using the `jackc/pgx` library.
      * **Concurrency**: The `generator` uses tickers and goroutines for scheduled, non-blocking task execution.
  * **Frontend Development (Vanilla JS)**:
      * **Dynamic Visualization**: Interactive charts rendered with **Plotly.js**.
      * **Live Data**: The frontend continuously polls the API every 500ms to create a real-time experience.
      * **User Interaction**: A slider control allows the user to adjust the forecast horizon on-the-fly.
  * **Mathematical Modeling**: A practical application of numerical computing (**Gonum**) to generate Fourier series features and perform a linear regression for time-series forecasting.

-----

## 🛠️ Technology Stack

| Category      | Technology                                                                                                                                                                    |
| :------------ | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Backend** |                                  |
| **Database** |                                                                                    |
| **Frontend** |   |
| **DevOps** |              |
| **Libraries** | `gonum/mat` (Linear Algebra), `jackc/pgx` (Postgres Driver), `plotly.js` (Charting)                                                                                             |

-----

## 🚀 Getting Started

### Prerequisites

  * [Docker](https://www.docker.com/products/docker-desktop/)
  * [Docker Compose](https://docs.docker.com/compose/install/) (Typically included with Docker Desktop)

### 1\. Clone & Launch

Clone the repository and use Docker Compose to build the images and start the services in detached mode.

```bash
# Clone the repository
git clone <your-repo-url>
cd <repo-name>

# Build and run all services in the background
docker-compose up --build -d
```

### 2\. Verify and Access

Check that all three containers are running:

```bash
docker-compose ps
```

You should see `db`, `generator`, and `server` with a status of `Up` or `Running`.

  * **View the UI**: Open your browser to **[http://localhost:8080](https://www.google.com/search?q=http://localhost:8080)**.
  * **Access the Database**: Connect to `localhost:5432` with user `user` and password `password`.

It may take a few seconds for the first data points to appear on the chart as the generator service begins populating the database.

### 3\. Stop the Application

To stop and remove the containers, run:

```bash
docker-compose down
```

-----

## ⚙️ Configuration

Application behavior can be modified through environment variables in `docker-compose.yml` or constants in the source code.

| Parameter           | Location                                | Description                                                               |
| :------------------ | :-------------------------------------- | :------------------------------------------------------------------------ |
| **DB Connection** | `docker-compose.yml` (`PG_CONN`)        | The DSN connection string used by both the `server` and `generator`.      |
| **Generator Rate** | `cmd/generator/main.go` (`time.Ticker`) | The interval at which new data points are generated (default: `500ms`).   |
| **Forecast Period** | `cmd/server/main.go` (`period`)         | The assumed period of the primary cycle in the data (default: `100`).     |
| **Harmonics** | `cmd/server/main.go` (`harmonics`)      | The number of Fourier terms (sin/cos pairs) to use in the model (default: `3`). |

To apply changes made in the Go source code, you must rebuild the Docker images: `docker-compose up --build`.

-----

## 🔬 Development & Extension Ideas

This project is a solid foundation. Here are some ways it could be extended to demonstrate further skills:

  * **Implement a CI/CD Pipeline**: Add a `.github/workflows/ci.yml` or `.gitlab-ci.yml` to automatically build, test, and push the Docker images to a registry (e.g., Docker Hub, GHCR).
  * **Add Unit & Integration Tests**: Write Go tests for the `data` package functions (`LinearFit`, `All`, etc.) and API endpoints to ensure correctness and prevent regressions.
  * **Introduce Health Checks**: Implement a `/healthz` endpoint in the `server` and add a `healthcheck` instruction in the `docker-compose.yml` to ensure services are truly operational.
  * **Improve the Prediction Model**: The current model is stateless. It could be enhanced by:
      * Adding a linear trend component to the feature set.
      * Replacing the linear regression with a more advanced model like ARIMA or a simple recurrent neural network (RNN).
  * **Switch to WebSockets**: Replace the frontend's polling mechanism with a WebSocket connection for more efficient, push-based data updates from the server.
