# Infrastructure Health Dashboard - At A Glance

## What Is This Project?

This project is a **monitoring system** that tracks the health of computer infrastructure—things like virtual machines, servers, and containers running in the cloud. Think of it as a fitness tracker, but for computers instead of people.

Just like a fitness tracker monitors your heart rate and steps, this dashboard monitors:
- **CPU usage** (how hard the computer is working)
- **Memory usage** (how much information it's holding at once)
- **Health status** (is the system running smoothly or having problems?)

---

## Why Does This Exist?

When companies run applications in the cloud, they often have dozens or hundreds of servers working together. If one server starts struggling or fails, it can affect customers. This dashboard helps teams:

- **See problems early** before they affect users
- **Track performance over time** to spot trends
- **Get alerts** when something needs attention

---

## The Technology Stack

### Programming Language: Go

**Go** (also called Golang) is the programming language used to build this project. Created by Google, it's known for being fast and efficient—perfect for systems that need to handle many requests quickly.

### Database: MySQL

**MySQL** is where all the data gets stored. It's a database—essentially a very organized filing cabinet for digital information. When the system records that a server's CPU is at 45%, that information goes into MySQL so we can look at it later.

### Web Framework: Gin

**Gin** is a toolkit that helps build web APIs (Application Programming Interfaces). An API is like a waiter in a restaurant—it takes requests ("I'd like the CPU data for server #5"), goes to the kitchen (database), and brings back what you asked for.

### Container Platform: Docker

**Docker** packages the application and everything it needs into a "container"—a self-contained box that runs the same way on any computer. This makes it easy to deploy the application anywhere without worrying about setup differences.

### Cloud Platform: Azure (Planned)

**Microsoft Azure** is a cloud computing platform. In future phases, this project will connect to Azure to automatically pull real metrics from cloud resources instead of using test data.

### Kubernetes (Planned)

**Kubernetes** (often called K8s) is a system for managing containers at scale. When you have many containers running, Kubernetes keeps them organized, restarts them if they crash, and distributes work evenly.

---

## How It Works (Simple Version)

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│   1. COLLECT        2. STORE           3. SERVE             │
│                                                             │
│   Cloud servers     Database           Dashboard/API        │
│   send metrics  →   saves them    →    shows the data       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

1. **Collect**: The system gathers information from servers (CPU %, memory %, health status)
2. **Store**: That information is saved in a database with timestamps
3. **Serve**: Users and other systems can request the data through the API

---

## What's Inside the Project

### The API (Application Programming Interface)

The API is the main component. It's a web service that accepts requests and returns data. Here's what you can ask it:

| Request | What It Does |
|---------|--------------|
| "Show me all resources" | Lists every server/system being monitored |
| "Add a new resource" | Registers a new server to monitor |
| "Get metrics for server #5" | Returns CPU, memory, etc. for that server |
| "What's the health status?" | Shows if systems are healthy or having issues |

### The Database Tables

Data is organized into four tables (think of these as spreadsheets):

| Table | What It Stores |
|-------|----------------|
| **Resources** | The list of things being monitored (servers, VMs, containers) |
| **Metrics** | Performance measurements (CPU at 45% at 2:30 PM) |
| **Health Status** | Current state (healthy, degraded, unhealthy) |
| **Alerts** | Notifications when something goes wrong |

### Project Files

| Folder | Purpose |
|--------|---------|
| `cmd/server/` | The starting point—runs the application |
| `internal/api/` | Handles web requests and responses |
| `internal/models/` | Defines data structures (what a "resource" or "metric" looks like) |
| `internal/repository/` | Talks to the database |
| `migrations/` | Sets up the database tables |

---

## Key Features

### Automatic Setup
When the application starts, it automatically creates the database tables if they don't exist. No manual database setup required.

### Real-Time Data
The API responds instantly to requests, returning the latest data from the database.

### Historical Tracking
Every metric is stored with a timestamp, allowing you to see trends over time (e.g., "CPU usage has been climbing over the past week").

### Flexible Filtering
You can request specific data:
- Only metrics from a certain time range
- Only a specific type of metric (just CPU, or just memory)
- Only resources of a certain type (just VMs, or just containers)

---

## The Development Roadmap

This project is built in phases:

| Phase | Goal | Status |
|-------|------|--------|
| **Phase 1** | Build the basic API with database storage | ✅ Complete |
| **Phase 2** | Connect to Azure cloud for real metrics | 🔲 Planned |
| **Phase 3** | Package in Docker containers | 🔲 Planned |
| **Phase 4** | Run on local Kubernetes | 🔲 Planned |
| **Phase 5** | Deploy to Azure Kubernetes Service | 🔲 Planned |
| **Phase 6** | Monitor Kubernetes itself | 🔲 Planned |
| **Phase 7** | Production polish (logging, alerts) | 🔲 Planned |

---

## Who Is This For?

- **DevOps Engineers**: Monitor infrastructure health in one place
- **System Administrators**: Track server performance over time
- **Development Teams**: Ensure their applications have healthy underlying infrastructure
- **Anyone learning**: This project demonstrates real-world patterns for building monitoring systems

---

## Quick Facts

| Aspect | Detail |
|--------|--------|
| **Language** | Go (Golang) |
| **Database** | MySQL |
| **API Type** | REST (standard web API) |
| **Default Port** | 8080 |
| **Data Format** | JSON |
| **License** | Personal project |

---

## Glossary

| Term | Simple Explanation |
|------|-------------------|
| **API** | A way for programs to talk to each other over the internet |
| **Database** | Organized storage for information |
| **Container** | A packaged application that runs the same everywhere |
| **Cloud** | Computers owned by companies like Microsoft/Amazon that you rent |
| **Metrics** | Measurements (like CPU percentage) |
| **Endpoint** | A specific URL the API responds to |
| **JSON** | A text format for data that looks like `{"name": "value"}` |
| **VM (Virtual Machine)** | A simulated computer running inside a real computer |
| **Kubernetes** | Software that manages many containers |

---

## Summary

This Infrastructure Health Dashboard is a monitoring tool built with modern technologies. It collects performance data from computer systems, stores it in a database, and provides that data through a web API. The project demonstrates industry-standard practices for building reliable, scalable monitoring solutions.

While currently focused on local development with mock data, the roadmap includes connecting to real cloud resources in Microsoft Azure and deploying on Kubernetes—the same technologies used by major companies worldwide.
