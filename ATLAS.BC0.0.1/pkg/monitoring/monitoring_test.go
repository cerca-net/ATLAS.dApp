package monitoring

import (
	"atlas-blockchain/pkg/sharding"
	"testing"
	"time"
)

func TestRealMonitoringSystem(t *testing.T) {
	// Test 1: Monitor creation
	t.Run("Monitor Creation", func(t *testing.T) {
		monitor := NewMonitor(nil)
		if monitor == nil {
			t.Fatal("Monitor should not be nil")
		}

		if monitor.startTime.IsZero() {
			t.Fatal("Monitor start time should be set")
		}

		if monitor.metrics == nil {
			t.Fatal("Monitor metrics map should be initialized")
		}

		if monitor.healthChecks == nil {
			t.Fatal("Monitor health checks map should be initialized")
		}

		if monitor.alerts == nil {
			t.Fatal("Monitor alerts map should be initialized")
		}
	})

	// Test 2: Real system resource monitoring
	t.Run("System Resource Monitoring", func(t *testing.T) {
		monitor := NewMonitor(nil)

		// Perform health checks
		monitor.performHealthChecks()

		// Check that health checks were created
		healthChecks := monitor.GetHealthChecks()
		if len(healthChecks) == 0 {
			t.Fatal("Health checks should be performed")
		}

		// Find system resources check
		var systemCheck *HealthCheck
		for _, check := range healthChecks {
			if check.Name == "system_resources" {
				systemCheck = check
				break
			}
		}

		if systemCheck == nil {
			t.Fatal("System resources health check should exist")
		}

		// Verify the check has real data
		if systemCheck.Details == nil {
			t.Fatal("System resources check should have details")
		}

		// Check for required fields
		requiredFields := []string{"memory_usage", "cpu_usage", "disk_usage"}
		for _, field := range requiredFields {
			if _, exists := systemCheck.Details[field]; !exists {
				t.Errorf("System resources check should have %s field", field)
			}
		}
	})

	// Test 3: Real performance metrics
	t.Run("Performance Metrics", func(t *testing.T) {
		monitor := NewMonitor(nil)

		// Update performance metrics
		monitor.updatePerformanceMetrics()

		// Get performance metrics
		metrics := monitor.GetPerformanceMetrics()
		if metrics == nil {
			t.Fatal("Performance metrics should not be nil")
		}

		// Verify metrics are real (not hardcoded)
		if metrics.MemoryUsage <= 0 {
			t.Error("Memory usage should be greater than 0")
		}

		if metrics.CPUUsage < 0 || metrics.CPUUsage > 100 {
			t.Error("CPU usage should be between 0 and 100")
		}

		if metrics.DiskUsage < 0 || metrics.DiskUsage > 100 {
			t.Error("Disk usage should be between 0 and 100")
		}
	})

	// Test 4: Transaction and block time recording
	t.Run("Transaction and Block Recording", func(t *testing.T) {
		monitor := NewMonitor(nil)

		// Record some transactions
		monitor.RecordTransaction()
		monitor.RecordTransaction()
		monitor.RecordTransaction()

		// Record some block times
		monitor.RecordBlockTime(12 * time.Second)
		monitor.RecordBlockTime(15 * time.Second)
		monitor.RecordBlockTime(10 * time.Second)

		// Update performance metrics to recalculate TPS and block time
		monitor.updatePerformanceMetrics()

		metrics := monitor.GetPerformanceMetrics()
		if metrics == nil {
			t.Fatal("Performance metrics should not be nil")
		}

		// TPS should be calculated (even if small)
		if metrics.TPS < 0 {
			t.Error("TPS should not be negative")
		}

		// Block time should be calculated from recorded times
		if metrics.BlockTime <= 0 {
			t.Error("Block time should be greater than 0")
		}

		// Block time should be close to the average of recorded times (12.33 seconds)
		expectedBlockTime := 12.33 // Average of 12, 15, 10
		if metrics.BlockTime < expectedBlockTime-2 || metrics.BlockTime > expectedBlockTime+2 {
			t.Errorf("Block time should be close to %f, got %f", expectedBlockTime, metrics.BlockTime)
		}
	})

	// Test 5: Alert system
	t.Run("Alert System", func(t *testing.T) {
		monitor := NewMonitor(nil)

		// Verify initial state
		initialAlerts := monitor.GetAlerts()
		if len(initialAlerts) != 0 {
			t.Errorf("Expected 0 initial alerts, got %d", len(initialAlerts))
		}

		// Add first alert
		monitor.AddAlert("warning", "Test warning alert", map[string]interface{}{
			"test_field": "test_value",
		})

		// Check after first alert
		alerts1 := monitor.GetAlerts()
		if len(alerts1) != 1 {
			t.Errorf("Expected 1 alert after first addition, got %d", len(alerts1))
		}

		// Add second alert
		monitor.AddAlert("error", "Test error alert", nil)

		// Get alerts
		alerts := monitor.GetAlerts()
		if len(alerts) != 2 {
			t.Errorf("Expected 2 alerts, got %d", len(alerts))
		}

		// Check alert properties
		for _, alert := range alerts {
			if alert.ID == "" {
				t.Error("Alert should have an ID")
			}

			if alert.Message == "" {
				t.Error("Alert should have a message")
			}

			if alert.Timestamp.IsZero() {
				t.Error("Alert should have a timestamp")
			}

			if alert.Acknowledged {
				t.Error("New alerts should not be acknowledged")
			}
		}

		// Test alert acknowledgment
		if len(alerts) > 0 {
			alertID := alerts[0].ID
			err := monitor.AcknowledgeAlert(alertID)
			if err != nil {
				t.Errorf("Failed to acknowledge alert: %v", err)
			}

			// Check that alert is now acknowledged
			updatedAlerts := monitor.GetAlerts()
			for _, alert := range updatedAlerts {
				if alert.ID == alertID && !alert.Acknowledged {
					t.Error("Alert should be acknowledged")
				}
			}
		}
	})

	// Test 6: Network monitoring
	t.Run("Network Monitoring", func(t *testing.T) {
		monitor := NewMonitor(nil)

		// Perform health checks to trigger network monitoring
		monitor.performHealthChecks()

		// Check that network connectivity check was performed
		healthChecks := monitor.GetHealthChecks()
		var networkCheck *HealthCheck
		for _, check := range healthChecks {
			if check.Name == "network_connectivity" {
				networkCheck = check
				break
			}
		}

		if networkCheck == nil {
			t.Fatal("Network connectivity health check should exist")
		}

		// Verify network check has real data
		if networkCheck.Details == nil {
			t.Fatal("Network connectivity check should have details")
		}

		// Check for network-related fields
		requiredFields := []string{"bytes_sent", "bytes_received", "interfaces"}
		for _, field := range requiredFields {
			if _, exists := networkCheck.Details[field]; !exists {
				t.Errorf("Network connectivity check should have %s field", field)
			}
		}
	})

	// Test 7: Shard health monitoring
	t.Run("Shard Health Monitoring", func(t *testing.T) {
		// Create a proper shard manager with configuration
		shardConfig := &sharding.ShardConfig{
			TotalShards:     4,
			ShardSize:       10,
			CrossShardDelay: time.Second * 5,
			ConsensusType:   "pbft",
		}
		shardManager := sharding.NewShardManager(shardConfig)
		monitor := NewMonitor(shardManager)

		// Perform health checks
		monitor.performHealthChecks()

		// Check that shard health check was performed
		healthChecks := monitor.GetHealthChecks()
		var shardCheck *HealthCheck
		for _, check := range healthChecks {
			if check.Name == "shard_health" {
				shardCheck = check
				break
			}
		}

		if shardCheck == nil {
			t.Fatal("Shard health check should exist")
		}

		// Verify shard check has data
		if shardCheck.Details == nil {
			t.Fatal("Shard health check should have details")
		}

		// Check for shard-related fields
		requiredFields := []string{"active_shards", "total_shards"}
		for _, field := range requiredFields {
			if _, exists := shardCheck.Details[field]; !exists {
				t.Errorf("Shard health check should have %s field", field)
			}
		}
	})

	// Test 8: Metric recording and retrieval
	t.Run("Metric Recording", func(t *testing.T) {
		monitor := NewMonitor(nil)

		// Record some metrics
		monitor.RecordMetric("test_metric", MetricTypeGauge, 42.5, map[string]string{
			"label1": "value1",
			"label2": "value2",
		})

		monitor.RecordMetric("test_counter", MetricTypeCounter, 100, nil)

		// Get all metrics
		allMetrics := monitor.GetAllMetrics()
		if len(allMetrics) != 2 {
			t.Errorf("Expected 2 metrics, got %d", len(allMetrics))
		}

		// Check specific metric
		metric, exists := monitor.GetMetric("test_metric")
		if !exists {
			t.Fatal("test_metric should exist")
		}

		if metric.Value != 42.5 {
			t.Errorf("Expected metric value 42.5, got %f", metric.Value)
		}

		if metric.Type != MetricTypeGauge {
			t.Errorf("Expected metric type gauge, got %s", metric.Type)
		}

		if len(metric.Labels) != 2 {
			t.Errorf("Expected 2 labels, got %d", len(metric.Labels))
		}

		if metric.Labels["label1"] != "value1" {
			t.Error("Label1 should be 'value1'")
		}
	})

	// Test 9: System status
	t.Run("System Status", func(t *testing.T) {
		monitor := NewMonitor(nil)

		// Add a small delay to ensure uptime is measurable
		time.Sleep(10 * time.Millisecond)

		// Add some health checks and alerts
		monitor.UpdateHealthCheck("test_check", "healthy", "Test check", nil)
		monitor.AddAlert("info", "Test alert", nil)

		// Get system status
		status := monitor.GetSystemStatus()
		if status == nil {
			t.Fatal("System status should not be nil")
		}

		// Verify status fields
		if status.Status == "" {
			t.Error("System status should have a status")
		}

		if status.Uptime <= 0 {
			t.Error("System status should have uptime")
		}

		if status.Version == "" {
			t.Error("System status should have a version")
		}

		if status.Performance == nil {
			t.Error("System status should have performance metrics")
		}

		if len(status.HealthChecks) == 0 {
			t.Error("System status should have health checks")
		}

		if len(status.Alerts) == 0 {
			t.Error("System status should have alerts")
		}
	})
}
