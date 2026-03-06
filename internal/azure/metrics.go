package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
)

// MetricDefinition represents an available metric for a resource
type MetricDefinition struct {
	Name        string
	DisplayName string
	Unit        string
	Aggregation string
}

// MetricValue represents a single metric data point
type MetricValue struct {
	MetricName string
	Timestamp  time.Time
	Value      float64
	Unit       string
}

// VMMetricNames contains the metric names to collect for VMs
var VMMetricNames = []string{
	"Percentage CPU",
	"Available Memory Bytes",
	"Disk Read Bytes",
	"Disk Write Bytes",
	"Network In Total",
	"Network Out Total",
}

// MetricNameMapping maps Azure metric names to our internal names
var MetricNameMapping = map[string]string{
	"Percentage CPU":         "cpu_percent",
	"Available Memory Bytes": "memory_available_bytes",
	"Disk Read Bytes":        "disk_read_bytes",
	"Disk Write Bytes":       "disk_write_bytes",
	"Network In Total":       "network_in_bytes",
	"Network Out Total":      "network_out_bytes",
}

// MetricUnitMapping maps Azure metric names to units
var MetricUnitMapping = map[string]string{
	"Percentage CPU":         "percent",
	"Available Memory Bytes": "bytes",
	"Disk Read Bytes":        "bytes",
	"Disk Write Bytes":       "bytes",
	"Network In Total":       "bytes",
	"Network Out Total":      "bytes",
}

// ListAvailableMetrics lists all available metrics for a given resource
func (c *Client) ListAvailableMetrics(ctx context.Context, resourceURI string) ([]MetricDefinition, error) {
	// Create a metric definitions client
	definitionsClient, err := armmonitor.NewMetricDefinitionsClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric definitions client: %w", err)
	}

	var definitions []MetricDefinition
	pager := definitionsClient.NewListPager(resourceURI, nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list metric definitions: %w", err)
		}

		for _, def := range page.Value {
			if def.Name == nil || def.Name.Value == nil {
				continue
			}

			displayName := ""
			if def.Name.LocalizedValue != nil {
				displayName = *def.Name.LocalizedValue
			}

			unit := ""
			if def.Unit != nil {
				unit = string(*def.Unit)
			}

			// Get primary aggregation type
			aggregation := "Average"
			if len(def.SupportedAggregationTypes) > 0 && def.SupportedAggregationTypes[0] != nil {
				aggregation = string(*def.SupportedAggregationTypes[0])
			}

			definitions = append(definitions, MetricDefinition{
				Name:        *def.Name.Value,
				DisplayName: displayName,
				Unit:        unit,
				Aggregation: aggregation,
			})
		}
	}

	return definitions, nil
}

// FetchMetrics fetches metric values for a resource within a time range
// resourceURI is the full Azure resource ID
// metricNames is a slice of metric names to fetch (e.g., "Percentage CPU")
// timeRange specifies how far back to fetch (e.g., 5 minutes, 1 hour)
// interval specifies the granularity (e.g., PT1M for 1 minute, PT5M for 5 minutes)
func (c *Client) FetchMetrics(ctx context.Context, resourceURI string, metricNames []string, timeRange time.Duration, interval string) ([]MetricValue, error) {
	if len(metricNames) == 0 {
		return nil, fmt.Errorf("at least one metric name is required")
	}

	// Build metric names string (comma-separated)
	metricNamesStr := ""
	for i, name := range metricNames {
		if i > 0 {
			metricNamesStr += ","
		}
		metricNamesStr = metricNamesStr + name
	}

	// Calculate time range
	endTime := time.Now().UTC()
	startTime := endTime.Add(-timeRange)
	timespan := fmt.Sprintf("%s/%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	// Query metrics
	response, err := c.metricsClient.List(ctx, resourceURI, &armmonitor.MetricsClientListOptions{
		Timespan:    &timespan,
		Interval:    &interval,
		Metricnames: &metricNamesStr,
		Aggregation: ptrString("Average"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}

	var values []MetricValue

	for _, metric := range response.Value {
		if metric.Name == nil || metric.Name.Value == nil {
			continue
		}
		metricName := *metric.Name.Value

		// Determine unit
		unit := ""
		if metric.Unit != nil {
			unit = string(*metric.Unit)
		}

		// Extract time series data
		for _, ts := range metric.Timeseries {
			for _, dp := range ts.Data {
				if dp.TimeStamp == nil {
					continue
				}

				// Get the average value (most common aggregation)
				var value float64
				if dp.Average != nil {
					value = *dp.Average
				} else if dp.Total != nil {
					value = *dp.Total
				} else if dp.Maximum != nil {
					value = *dp.Maximum
				} else if dp.Minimum != nil {
					value = *dp.Minimum
				} else if dp.Count != nil {
					value = *dp.Count
				} else {
					continue // No value available
				}

				values = append(values, MetricValue{
					MetricName: metricName,
					Timestamp:  *dp.TimeStamp,
					Value:      value,
					Unit:       unit,
				})
			}
		}
	}

	return values, nil
}

// FetchLatestMetrics fetches the most recent metric values (last 5 minutes, 1 minute interval)
func (c *Client) FetchLatestMetrics(ctx context.Context, resourceURI string, metricNames []string) ([]MetricValue, error) {
	return c.FetchMetrics(ctx, resourceURI, metricNames, 5*time.Minute, "PT1M")
}

// FetchVMMetrics fetches standard VM metrics (CPU, memory, disk, network)
func (c *Client) FetchVMMetrics(ctx context.Context, resourceURI string) ([]MetricValue, error) {
	return c.FetchLatestMetrics(ctx, resourceURI, VMMetricNames)
}

// GetNormalizedMetricName converts Azure metric name to our internal name
func GetNormalizedMetricName(azureMetricName string) string {
	if name, ok := MetricNameMapping[azureMetricName]; ok {
		return name
	}
	return azureMetricName
}

// GetMetricUnit gets the unit for a metric
func GetMetricUnit(azureMetricName string) string {
	if unit, ok := MetricUnitMapping[azureMetricName]; ok {
		return unit
	}
	return ""
}

// helper function to create string pointer
func ptrString(s string) *string {
	return &s
}
