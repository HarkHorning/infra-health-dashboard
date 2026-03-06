# Infrastructure Health Dashboard

A Go-based service that collects metrics from Azure resources and Kubernetes clusters, stores historical data in SQL, and exposes REST APIs for monitoring infrastructure health.

## Tech Stack

- **Language:** Go
- **Cloud:** Azure (AKS, Azure SQL, Azure Monitor)
- **Container Orchestration:** Kubernetes (AKS)
- **Database:** MySQL or Azure SQL
- **APIs:** REST (with potential gRPC extension)

---

## Cost Considerations & Free Alternatives

### Why the Original Tech Choices?

The technologies in this project were chosen because they reflect **real-world production environments**:

| Original Choice | Why It's Used in Production |
|-----------------|----------------------------|
| **Azure** | One of the "Big 3" cloud providers (AWS, Azure, GCP). Many enterprises use Azure, especially those with Microsoft ecosystems. |
| **AKS** | Managed Kubernetes is the industry standard for container orchestration. Companies pay for managed services to reduce operational burden. |
| **Azure Container Registry** | Enterprise teams need private registries with security scanning, access controls, and integration with their cloud provider. |
| **Azure Database for MySQL** | Managed databases handle backups, patching, scaling, and high availability automatically—critical for production workloads. |
| **Azure Monitor** | Native integration with Azure resources, built-in dashboards, and alerting capabilities that enterprises rely on. |

**The tradeoff:** These production-grade services cost money. Azure's free tier is limited, and after the initial $200 credit expires, you'd pay ~$25-35/month.

---

### Free Tier Alternatives

For learning and development, you can substitute paid services with free alternatives:

#### Database Options

| Service | Free Tier | Limitations | Best For |
|---------|-----------|-------------|----------|
| **Docker MySQL** (local) | ✅ Completely free | Only works locally | Phase 1-4 development |
| **[PlanetScale](https://planetscale.com)** | ✅ 5GB storage, 1 billion row reads/mo | Serverless MySQL, no traditional connections | Cloud deployment without Azure |
| **[Railway](https://railway.app)** | ✅ $5 free credit/month | Limited hours | Quick cloud testing |
| **[Neon](https://neon.tech)** | ✅ 512MB storage | PostgreSQL, not MySQL (requires code changes) | If you're okay switching databases |

**Why not originally included:** PlanetScale and Railway are great for small projects but aren't commonly used in enterprise environments. Employers expect familiarity with major cloud providers' managed database offerings.

#### Container Registry Options

| Service | Free Tier | Limitations | Best For |
|---------|-----------|-------------|----------|
| **[GitHub Container Registry](https://ghcr.io)** | ✅ Free for public repos | Private repos count against storage quota | Open source projects |
| **[Docker Hub](https://hub.docker.com)** | ✅ 1 free private repo | Rate limits on pulls | Simple projects |
| **Local Registry** | ✅ Completely free | Only works locally | Local Kubernetes testing |

**Why not originally included:** Azure Container Registry integrates seamlessly with AKS (no authentication setup needed). In production, teams use their cloud provider's registry for security, compliance, and operational simplicity.

#### Kubernetes Options

| Service | Free Tier | Limitations | Best For |
|---------|-----------|-------------|----------|
| **[kind](https://kind.sigs.k8s.io)** | ✅ Completely free | Runs in Docker, local only | Learning Kubernetes basics |
| **[minikube](https://minikube.sigs.k8s.io)** | ✅ Completely free | Local only, resource-heavy | Local development |
| **[Docker Desktop K8s](https://www.docker.com/products/docker-desktop)** | ✅ Free for personal use | Local only | Easiest setup for Windows |
| **[Civo](https://www.civo.com)** | ✅ $250 free credit | Credit expires | Cheap managed K8s |
| **[Google GKE Autopilot](https://cloud.google.com/kubernetes-engine)** | ⚠️ $200 free credit | Credit expires after 90 days | Alternative to AKS |

**Why not originally included:** Local Kubernetes (kind/minikube) is already in Phase 4. AKS was chosen for Phase 5 because it demonstrates real cloud deployment skills. Employers often use managed Kubernetes, and experience with any major provider (AKS/EKS/GKE) is valuable.

#### Cloud Metrics Alternatives

| Service | Free Tier | Limitations | Best For |
|---------|-----------|-------------|----------|
| **Mock Data** | ✅ Completely free | Not real metrics | Learning the codebase |
| **[Prometheus](https://prometheus.io)** | ✅ Open source | Self-hosted, more setup | Kubernetes-native monitoring |
| **Local System Metrics** | ✅ Completely free | Only monitors your machine | Testing collectors |

**Why not originally included:** Azure Monitor demonstrates integration with a real cloud provider's monitoring stack. However, Prometheus is equally valuable to learn—it's the industry standard for Kubernetes monitoring.

---

### Recommended Path: 100% Free

To complete this project without spending money:

| Phase | What to Use | Cost |
|-------|-------------|------|
| **Phase 1** | Local Go + Docker MySQL | Free |
| **Phase 2** | Skip OR use mock Azure data | Free |
| **Phase 3** | Docker + Docker Hub or GHCR | Free |
| **Phase 4** | kind or minikube | Free |
| **Phase 5** | Stay on local K8s OR use Civo free credit | Free |
| **Phase 6** | Local Kubernetes metrics | Free |
| **Phase 7** | Prometheus + Grafana (self-hosted) | Free |

### Recommended Path: Use Azure Free Credits

If you want the full Azure experience:

1. **Sign up** for Azure free account ($200 credit for 30 days)
2. **Complete Phases 1-4** locally first (no Azure costs)
3. **Speed-run Phases 5-7** within the 30-day window
4. **Delete all resources** when done to avoid charges
5. **Document everything** with screenshots for your portfolio

---

### Phase-by-Phase Alternatives

Below are specific commands for free alternatives in each phase:

#### Phase 3: Free Container Registry

Instead of Azure Container Registry, use GitHub Container Registry:

```powershell
# Login to GitHub Container Registry
echo $env:GITHUB_TOKEN | docker login ghcr.io -u YOUR_USERNAME --password-stdin

# Tag and push
docker tag infra-dashboard ghcr.io/YOUR_USERNAME/infra-dashboard:v1
docker push ghcr.io/YOUR_USERNAME/infra-dashboard:v1
```

#### Phase 5: Free Database (PlanetScale)

Instead of Azure Database for MySQL:

1. Create account at [planetscale.com](https://planetscale.com)
2. Create a database named `infra_dashboard`
3. Get connection string from dashboard
4. Update environment variables:
```powershell
$env:DB_HOST = "aws.connect.psdb.cloud"
$env:DB_USER = "your-username"
$env:DB_PASSWORD = "your-password"
$env:DB_NAME = "infra_dashboard"
```

Note: PlanetScale uses a different connection method. You may need to add `?tls=true` to the DSN.

#### Phase 5: Skip AKS, Stay Local

Instead of deploying to AKS, enhance your local Kubernetes setup:

```powershell
# Use kind with a custom config for more realistic setup
kind create cluster --config kind-config.yaml
```

Create `kind-config.yaml`:
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
```

This gives you a multi-node cluster locally—similar to a real cloud environment.

---

### Summary: Original vs Free

| Component | Original (Production) | Free Alternative |
|-----------|----------------------|------------------|
| Database | Azure Database for MySQL | Docker MySQL / PlanetScale |
| Container Registry | Azure Container Registry | GitHub Container Registry |
| Kubernetes | Azure Kubernetes Service | kind / minikube |
| Cloud Metrics | Azure Monitor | Mock data / Prometheus |
| VM to Monitor | Azure VM | Local machine / Skip |

**Bottom line:** You can complete this entire project for $0. The Azure services teach cloud-specific skills, but the core concepts (Go, Kubernetes, databases, APIs) are the same regardless of which platform you use.

## Architecture Overview

```
┌─────────────────┐     ┌─────────────────┐
│  Azure Monitor  │     │   Kubernetes    │
│     Metrics     │     │   Metrics API   │
└────────┬────────┘     └────────┬────────┘
         │                       │
         └───────────┬───────────┘
                     │
              ┌──────▼──────┐
              │   Go API    │
              │   Service   │
              └──────┬──────┘
                     │
              ┌──────▼──────┐
              │    MySQL    │
              │  (History)  │
              └─────────────┘
```

---

## Data Model

### What Metrics to Collect

#### Azure Resource Metrics (Phase 2)
| Metric | Description | Source |
|--------|-------------|--------|
| CPU Percentage | VM/Container CPU usage | Azure Monitor |
| Memory Percentage | RAM usage | Azure Monitor |
| Disk Read/Write Bytes | Storage I/O | Azure Monitor |
| Network In/Out | Network traffic bytes | Azure Monitor |
| Resource Health Status | Healthy/Unhealthy/Degraded | Azure Resource Health API |

#### Kubernetes Metrics (Phase 6)
| Metric | Description | Source |
|--------|-------------|--------|
| Pod CPU Usage | CPU millicores used by pod | Metrics API |
| Pod Memory Usage | Memory bytes used by pod | Metrics API |
| Node CPU Capacity/Usage | Node-level CPU | Metrics API |
| Node Memory Capacity/Usage | Node-level memory | Metrics API |
| Pod Status | Running/Pending/Failed | Kubernetes API |
| Pod Restart Count | Number of container restarts | Kubernetes API |

### Database Schema

```sql
-- Resources being monitored (VMs, pods, nodes, etc.)
CREATE TABLE resources (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource_id VARCHAR(255) NOT NULL UNIQUE,  -- Azure resource ID or K8s uid
    resource_type VARCHAR(50) NOT NULL,        -- 'azure_vm', 'k8s_pod', 'k8s_node'
    name VARCHAR(255) NOT NULL,
    region VARCHAR(100),
    tags JSON,                                 -- Flexible metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Time-series metrics data
CREATE TABLE metrics (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    resource_id INT NOT NULL,
    metric_name VARCHAR(100) NOT NULL,         -- 'cpu_percent', 'memory_percent', etc.
    metric_value DOUBLE NOT NULL,
    unit VARCHAR(50),                          -- 'percent', 'bytes', 'count'
    recorded_at TIMESTAMP NOT NULL,            -- When the metric was recorded
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES resources(id),
    INDEX idx_resource_metric_time (resource_id, metric_name, recorded_at)
);

-- Resource health status snapshots
CREATE TABLE health_status (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource_id INT NOT NULL,
    status VARCHAR(20) NOT NULL,               -- 'healthy', 'unhealthy', 'degraded', 'unknown'
    reason TEXT,                               -- Why it's unhealthy
    checked_at TIMESTAMP NOT NULL,
    FOREIGN KEY (resource_id) REFERENCES resources(id),
    INDEX idx_resource_checked (resource_id, checked_at)
);

-- Optional: Alert history
CREATE TABLE alerts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource_id INT NOT NULL,
    alert_type VARCHAR(100) NOT NULL,          -- 'high_cpu', 'pod_crash_loop', etc.
    severity VARCHAR(20) NOT NULL,             -- 'info', 'warning', 'critical'
    message TEXT,
    triggered_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES resources(id)
);
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Service health check |
| GET | `/api/v1/resources` | List all monitored resources |
| GET | `/api/v1/resources/{id}` | Get single resource details |
| POST | `/api/v1/resources` | Register a resource to monitor |
| DELETE | `/api/v1/resources/{id}` | Stop monitoring a resource |
| GET | `/api/v1/resources/{id}/metrics` | Get metrics for a resource (query params: metric_name, from, to) |
| GET | `/api/v1/resources/{id}/health` | Get latest health status |
| GET | `/api/v1/metrics/summary` | Aggregated metrics across all resources |
| GET | `/api/v1/alerts` | List recent alerts |

---

## Roadmap

### Phase 1: Local Development Foundation
**Goal:** Get a working Go API with MySQL storage running locally.

#### 1.1 Project Setup
```
infra-health-dashboard/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers.go       # HTTP handlers
│   │   ├── router.go         # Route definitions
│   │   └── middleware.go     # Logging, auth, etc.
│   ├── models/
│   │   ├── resource.go       # Resource struct
│   │   ├── metric.go         # Metric struct
│   │   └── health.go         # Health status struct
│   ├── repository/
│   │   ├── resource_repo.go  # Resource DB operations
│   │   ├── metric_repo.go    # Metric DB operations
│   │   └── mysql.go          # DB connection setup
│   └── service/
│       └── collector.go      # Metric collection logic
├── configs/
│   └── config.yaml           # Configuration file
├── migrations/
│   └── 001_initial.sql       # Database migrations
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

#### 1.2 Tasks
- [ ] Initialize Go module: `go mod init github.com/yourusername/infra-health-dashboard`
- [ ] Install dependencies:
  ```powershell
  go get github.com/gin-gonic/gin          # HTTP router
  go get github.com/go-sql-driver/mysql    # MySQL driver
  go get github.com/jmoiron/sqlx           # SQL helper library
  ```
- [ ] Create `cmd/server/main.go` with basic HTTP server
- [ ] Implement `/health` endpoint that returns `{"status": "ok"}`
- [ ] Run MySQL via Docker:
  ```powershell
  docker run -d --name mysql-dev `
    -e MYSQL_ROOT_PASSWORD=devpassword `
    -e MYSQL_DATABASE=infra_dashboard `
    -p 3306:3306 `
    mysql:8
  ```
- [ ] Create database tables using the schema above
- [ ] Implement `internal/repository/mysql.go` - connection pool setup
- [ ] Implement `internal/repository/resource_repo.go`:
  - `Create(resource)` - insert new resource
  - `GetByID(id)` - fetch single resource
  - `List()` - fetch all resources
  - `Delete(id)` - remove resource
- [ ] Implement `internal/repository/metric_repo.go`:
  - `Insert(metric)` - store a metric reading
  - `GetByResourceID(resourceID, from, to)` - fetch metrics in time range
- [ ] Implement handlers for `/api/v1/resources` endpoints
- [ ] Implement handlers for `/api/v1/resources/{id}/metrics`
- [ ] Add mock data seeder for testing
- [ ] Write unit tests for repository layer

**Test your Phase 1:**
```powershell
# Start the server
go run cmd/server/main.go

# Test health endpoint
Invoke-RestMethod -Uri http://localhost:8080/health

# Create a mock resource
$body = @{
    resource_id = "test-vm-001"
    resource_type = "azure_vm"
    name = "Test VM"
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/api/v1/resources `
    -Method Post `
    -ContentType "application/json" `
    -Body $body

# Insert a mock metric
$metricBody = @{
    metric_name = "cpu_percent"
    metric_value = 45.2
    unit = "percent"
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/api/v1/resources/1/metrics `
    -Method Post `
    -ContentType "application/json" `
    -Body $metricBody

# Query metrics
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/resources/1/metrics?metric_name=cpu_percent"
```

**Deliverable:** A Go API that can store and retrieve mock metrics from a local database.

---

### Phase 2: Azure Integration
**Goal:** Connect to real Azure resources and pull metrics.

#### 2.1 Azure Setup
- [ ] Create free Azure account (if needed)
- [ ] Install Azure CLI: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-windows
- [ ] Login: `az login`
- [ ] Create a resource group:
  ```powershell
  az group create --name infra-dashboard-rg --location eastus
  ```
- [ ] Create a test VM (to have something to monitor):
  ```powershell
  az vm create `
    --resource-group infra-dashboard-rg `
    --name test-vm `
    --image Ubuntu2204 `
    --size Standard_B1s `
    --admin-username azureuser `
    --generate-ssh-keys
  ```

#### 2.2 Authentication Setup
- [ ] Create a Service Principal for your app:
  ```powershell
  az ad sp create-for-rbac --name "infra-dashboard-sp" `
    --role "Monitoring Reader" `
    --scopes /subscriptions/{your-subscription-id}
  ```
  Save the output - you'll need `appId`, `password`, and `tenant`
- [ ] Store credentials in environment variables:
  ```powershell
  $env:AZURE_CLIENT_ID = "your-app-id"
  $env:AZURE_CLIENT_SECRET = "your-password"
  $env:AZURE_TENANT_ID = "your-tenant-id"
  $env:AZURE_SUBSCRIPTION_ID = "your-subscription-id"
  ```
- [ ] To persist environment variables across sessions, use:
  ```powershell
  [System.Environment]::SetEnvironmentVariable("AZURE_CLIENT_ID", "your-app-id", "User")
  [System.Environment]::SetEnvironmentVariable("AZURE_CLIENT_SECRET", "your-password", "User")
  [System.Environment]::SetEnvironmentVariable("AZURE_TENANT_ID", "your-tenant-id", "User")
  [System.Environment]::SetEnvironmentVariable("AZURE_SUBSCRIPTION_ID", "your-subscription-id", "User")
  ```

#### 2.3 Code Changes
- [ ] Install Azure SDK:
  ```powershell
  go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
  go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor
  go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources
  ```
- [ ] Create `internal/azure/client.go`:
  - Initialize Azure credentials using `azidentity.NewDefaultAzureCredential()`
  - Create monitor client
- [ ] Create `internal/azure/metrics.go`:
  - Function to list available metrics for a resource
  - Function to fetch metric values for a time range
- [ ] Create `internal/service/collector.go`:
  - Background goroutine that polls Azure every N minutes
  - Fetches CPU, memory, network metrics
  - Stores in database
- [ ] Add configuration for:
  - Azure credentials (from env vars)
  - Polling interval
  - Which resources to monitor

#### 2.4 Metrics to Fetch from Azure Monitor
```go
// Example metrics to request for a VM
metrics := []string{
    "Percentage CPU",
    "Available Memory Bytes",
    "Disk Read Bytes",
    "Disk Write Bytes",
    "Network In Total",
    "Network Out Total",
}
```

**Test your Phase 2:**
```powershell
# Run the collector
go run cmd/server/main.go

# After a few collection cycles, query the metrics
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/resources/1/metrics?from=2024-01-01T00:00:00Z"
```

**Deliverable:** Service pulls real metrics from Azure and stores them.

---

### Phase 3: Containerization
**Goal:** Package the application for Kubernetes deployment.

#### 3.1 Dockerfile
- [ ] Create `Dockerfile`:
  ```dockerfile
  # Build stage
  FROM golang:1.22-alpine AS builder
  WORKDIR /app
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux go build -o /infra-dashboard ./cmd/server

  # Run stage
  FROM alpine:latest
  RUN apk --no-cache add ca-certificates
  WORKDIR /root/
  COPY --from=builder /infra-dashboard .
  EXPOSE 8080
  CMD ["./infra-dashboard"]
  ```

#### 3.2 Docker Compose
- [ ] Create `docker-compose.yml`:
  ```yaml
  version: '3.8'
  services:
    app:
      build: .
      ports:
        - "8080:8080"
      environment:
        - DB_HOST=mysql
        - DB_PORT=3306
        - DB_USER=root
        - DB_PASSWORD=devpassword
        - DB_NAME=infra_dashboard
        - AZURE_CLIENT_ID=${AZURE_CLIENT_ID}
        - AZURE_CLIENT_SECRET=${AZURE_CLIENT_SECRET}
        - AZURE_TENANT_ID=${AZURE_TENANT_ID}
        - AZURE_SUBSCRIPTION_ID=${AZURE_SUBSCRIPTION_ID}
      depends_on:
        - mysql

    mysql:
      image: mysql:8
      environment:
        - MYSQL_ROOT_PASSWORD=devpassword
        - MYSQL_DATABASE=infra_dashboard
      ports:
        - "3306:3306"
      volumes:
        - mysql_data:/var/lib/mysql
        - ./migrations:/docker-entrypoint-initdb.d

  volumes:
    mysql_data:
  ```

#### 3.3 Tasks
- [ ] Build and test locally:
  ```powershell
  docker-compose up --build
  Invoke-RestMethod -Uri http://localhost:8080/health
  ```
- [ ] Create Azure Container Registry:
  ```powershell
  az acr create --resource-group infra-dashboard-rg --name youracrname --sku Basic
  az acr login --name youracrname
  ```
- [ ] Push image:
  ```powershell
  docker tag infra-dashboard youracrname.azurecr.io/infra-dashboard:v1
  docker push youracrname.azurecr.io/infra-dashboard:v1
  ```

**Deliverable:** Docker image ready for deployment.

---

### Phase 4: Kubernetes Basics (Local)
**Goal:** Learn Kubernetes fundamentals with local tooling.

#### 4.1 Local Cluster Setup
- [ ] Install Docker Desktop (includes Kubernetes) OR install kind:
  ```powershell
  # Using kind (Kubernetes in Docker)
  go install sigs.k8s.io/kind@latest
  kind create cluster --name infra-dashboard
  ```
- [ ] Verify: `kubectl cluster-info`

#### 4.2 Kubernetes Concepts to Learn
| Concept | What it does | Your use case |
|---------|--------------|---------------|
| Pod | Smallest deployable unit, runs containers | Runs your Go app |
| Deployment | Manages pod replicas, rolling updates | Ensures app stays running |
| Service | Stable network endpoint for pods | Exposes app internally/externally |
| ConfigMap | Non-sensitive configuration | DB host, polling interval |
| Secret | Sensitive data (base64 encoded) | DB password, Azure credentials |
| Namespace | Logical isolation | Organize your resources |

#### 4.3 Create Kubernetes Manifests
- [ ] Create `k8s/` directory:
  ```powershell
  New-Item -ItemType Directory -Path k8s -Force
  ```
- [ ] Create `k8s/namespace.yaml`:
  ```yaml
  apiVersion: v1
  kind: Namespace
  metadata:
    name: infra-dashboard
  ```
- [ ] Create `k8s/configmap.yaml`:
  ```yaml
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: app-config
    namespace: infra-dashboard
  data:
    DB_HOST: "mysql"
    DB_PORT: "3306"
    DB_NAME: "infra_dashboard"
    POLLING_INTERVAL: "300"  # seconds
  ```
- [ ] Create `k8s/secret.yaml` (don't commit real values!):
  ```yaml
  apiVersion: v1
  kind: Secret
  metadata:
    name: app-secrets
    namespace: infra-dashboard
  type: Opaque
  stringData:
    DB_USER: "root"
    DB_PASSWORD: "your-password"
    AZURE_CLIENT_ID: "your-client-id"
    AZURE_CLIENT_SECRET: "your-secret"
    AZURE_TENANT_ID: "your-tenant"
    AZURE_SUBSCRIPTION_ID: "your-sub"
  ```
- [ ] Create `k8s/deployment.yaml`:
  ```yaml
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: infra-dashboard
    namespace: infra-dashboard
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: infra-dashboard
    template:
      metadata:
        labels:
          app: infra-dashboard
      spec:
        containers:
        - name: app
          image: infra-dashboard:local
          ports:
          - containerPort: 8080
          envFrom:
          - configMapRef:
              name: app-config
          - secretRef:
              name: app-secrets
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 15
            periodSeconds: 20
  ```
- [ ] Create `k8s/service.yaml`:
  ```yaml
  apiVersion: v1
  kind: Service
  metadata:
    name: infra-dashboard
    namespace: infra-dashboard
  spec:
    selector:
      app: infra-dashboard
    ports:
    - port: 80
      targetPort: 8080
    type: ClusterIP
  ```

#### 4.4 Deploy Locally
- [ ] Load image into kind:
  ```powershell
  kind load docker-image infra-dashboard:local --name infra-dashboard
  ```
- [ ] Apply manifests:
  ```powershell
  kubectl apply -f k8s/namespace.yaml
  kubectl apply -f k8s/configmap.yaml
  kubectl apply -f k8s/secret.yaml
  kubectl apply -f k8s/deployment.yaml
  kubectl apply -f k8s/service.yaml
  ```
- [ ] Check status:
  ```powershell
  kubectl get pods -n infra-dashboard
  ```
- [ ] Port forward:
  ```powershell
  kubectl port-forward svc/infra-dashboard 8080:80 -n infra-dashboard
  ```
- [ ] Test:
  ```powershell
  Invoke-RestMethod -Uri http://localhost:8080/health
  ```

**Deliverable:** App running on local Kubernetes.

---

### Phase 5: AKS Deployment
**Goal:** Deploy to Azure Kubernetes Service.

#### 5.1 Create AKS Cluster
- [ ] Create cluster:
  ```powershell
  az aks create `
    --resource-group infra-dashboard-rg `
    --name infra-dashboard-aks `
    --node-count 1 `
    --node-vm-size Standard_B2s `
    --enable-managed-identity `
    --attach-acr youracrname `
    --generate-ssh-keys
  ```
- [ ] Get credentials:
  ```powershell
  az aks get-credentials --resource-group infra-dashboard-rg --name infra-dashboard-aks
  ```
- [ ] Verify:
  ```powershell
  kubectl get nodes
  ```

#### 5.2 Set Up MySQL
Option A - Azure Database for MySQL:
```powershell
az mysql flexible-server create `
  --resource-group infra-dashboard-rg `
  --name infra-dashboard-mysql `
  --admin-user adminuser `
  --admin-password YourPassword123! `
  --sku-name Standard_B1ms `
  --public-access 0.0.0.0
```

Option B - MySQL in Kubernetes (simpler for learning):
- [ ] Create `k8s/mysql.yaml` with PersistentVolumeClaim, Deployment, Service

#### 5.3 Deploy Application
- [ ] Update `k8s/deployment.yaml` to use ACR image:
  ```yaml
  image: youracrname.azurecr.io/infra-dashboard:v1
  ```
- [ ] Update secrets with real values
- [ ] Apply all manifests to AKS
- [ ] Expose externally:
  ```yaml
  # Change service type to LoadBalancer
  type: LoadBalancer
  ```
- [ ] Get external IP:
  ```powershell
  kubectl get svc -n infra-dashboard
  ```

**Deliverable:** Live application running on AKS.

---

### Phase 6: Kubernetes Metrics Collection
**Goal:** Add Kubernetes cluster monitoring to the dashboard.

#### 6.1 Tasks
- [ ] Install client-go:
  ```powershell
  go get k8s.io/client-go@latest
  go get k8s.io/metrics@latest
  ```
- [ ] Create `internal/k8s/client.go`:
  - Initialize Kubernetes client (in-cluster config when running in K8s)
  - Function to list pods
  - Function to list nodes
- [ ] Create `internal/k8s/metrics.go`:
  - Fetch pod metrics from Metrics API
  - Fetch node metrics
- [ ] Update collector to also gather K8s metrics
- [ ] Add new resource types: `k8s_pod`, `k8s_node`
- [ ] Create unified dashboard endpoint that shows all resources

**Deliverable:** Dashboard monitors both Azure resources and K8s cluster.

---

### Phase 7: Polish and Interview-Ready
**Goal:** Add features that demonstrate production-readiness.

- [ ] Add structured logging with `zerolog`:
  ```powershell
  go get github.com/rs/zerolog
  ```
- [ ] Implement graceful shutdown (handle SIGTERM)
- [ ] Add `/metrics` endpoint for Prometheus:
  ```powershell
  go get github.com/prometheus/client_golang/prometheus/promhttp
  ```
- [ ] Write integration tests
- [ ] Create GitHub Actions workflow (`.github/workflows/ci.yml`):
  - Run tests
  - Build Docker image
  - Push to ACR
- [ ] Document all API endpoints with examples
- [ ] Add basic alerting: if CPU > 80% for 5 minutes, create alert record

**Deliverable:** Production-quality codebase ready to discuss in interviews.

---

## Learning Resources

### Go
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go database/sql tutorial](https://go.dev/doc/tutorial/database-access)

### Kubernetes
- [Kubernetes Basics Tutorial](https://kubernetes.io/docs/tutorials/kubernetes-basics/)
- [kind - Local K8s](https://kind.sigs.k8s.io/)
- [client-go examples](https://github.com/kubernetes/client-go/tree/master/examples)

### Azure
- [Azure CLI Getting Started](https://learn.microsoft.com/en-us/cli/azure/get-started-with-azure-cli)
- [Azure SDK for Go](https://learn.microsoft.com/en-us/azure/developer/go/)
- [AKS Documentation](https://learn.microsoft.com/en-us/azure/aks/)
- [Azure Monitor REST API](https://learn.microsoft.com/en-us/rest/api/monitor/)

## Interview Talking Points

When discussing this project, be prepared to explain:

1. **Why Go?** - Performance, simplicity, strong concurrency model, popular in cloud-native tooling
2. **Database design decisions** - Schema choices, indexing for time-series queries, why you chose certain data types
3. **Kubernetes concepts** - How deployments, services, and secrets work; difference between ConfigMap and Secret
4. **Azure authentication** - Service principals vs managed identities, principle of least privilege
5. **Observability** - Logging, metrics, health checks, why structured logging matters
6. **CI/CD approach** - How you would automate deployments, what triggers builds
7. **Scaling considerations** - Horizontal scaling, database connection pooling, metric aggregation strategies
8. **Error handling** - What happens when Azure API is unavailable? How do you handle partial failures?
9. **Security** - How secrets are managed, network policies, RBAC
