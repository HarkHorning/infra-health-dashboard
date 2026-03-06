-- Infrastructure Health Dashboard - Initial Schema
-- Run this migration to set up the database tables

-- Resources being monitored (VMs, pods, nodes, etc.)
CREATE TABLE IF NOT EXISTS resources (
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
CREATE TABLE IF NOT EXISTS metrics (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    resource_id INT NOT NULL,
    metric_name VARCHAR(100) NOT NULL,         -- 'cpu_percent', 'memory_percent', etc.
    metric_value DOUBLE NOT NULL,
    unit VARCHAR(50),                          -- 'percent', 'bytes', 'count'
    recorded_at TIMESTAMP NOT NULL,            -- When the metric was recorded
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
    INDEX idx_resource_metric_time (resource_id, metric_name, recorded_at)
);

-- Resource health status snapshots
CREATE TABLE IF NOT EXISTS health_status (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource_id INT NOT NULL,
    status VARCHAR(20) NOT NULL,               -- 'healthy', 'unhealthy', 'degraded', 'unknown'
    reason TEXT,                               -- Why it's unhealthy
    checked_at TIMESTAMP NOT NULL,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
    INDEX idx_resource_checked (resource_id, checked_at)
);

-- Alert history
CREATE TABLE IF NOT EXISTS alerts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource_id INT NOT NULL,
    alert_type VARCHAR(100) NOT NULL,          -- 'high_cpu', 'pod_crash_loop', etc.
    severity VARCHAR(20) NOT NULL,             -- 'info', 'warning', 'critical'
    message TEXT,
    triggered_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
    INDEX idx_resource_alert (resource_id, triggered_at)
);
