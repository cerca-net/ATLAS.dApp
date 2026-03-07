package monitoring

import (
	"testing"
)

func TestAlertDebug(t *testing.T) {
	// Create a fresh monitor
	monitor := NewMonitor(nil)
	
	// Add first alert
	monitor.AddAlert("warning", "Test warning alert", nil)
	
	// Check immediately
	alerts1 := monitor.GetAlerts()
	t.Logf("After first alert: %d alerts", len(alerts1))
	for i, alert := range alerts1 {
		t.Logf("Alert %d: ID=%s, Level=%s, Message=%s", i, alert.ID, alert.Level, alert.Message)
	}
	
	// Add second alert
	monitor.AddAlert("error", "Test error alert", nil)
	
	// Check immediately
	alerts2 := monitor.GetAlerts()
	t.Logf("After second alert: %d alerts", len(alerts2))
	for i, alert := range alerts2 {
		t.Logf("Alert %d: ID=%s, Level=%s, Message=%s", i, alert.ID, alert.Level, alert.Message)
	}
	
	if len(alerts2) != 2 {
		t.Errorf("Expected 2 alerts, got %d", len(alerts2))
	}
} 