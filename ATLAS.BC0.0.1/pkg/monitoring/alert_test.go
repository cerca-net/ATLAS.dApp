package monitoring

import (
	"testing"
)

func TestAlertSystemStandalone(t *testing.T) {
	// Create a fresh monitor
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
	
	// Check after second alert
	alerts2 := monitor.GetAlerts()
	if len(alerts2) != 2 {
		t.Errorf("Expected 2 alerts after second addition, got %d", len(alerts2))
	}
	
	// Verify alert properties
	for i, alert := range alerts2 {
		if alert.ID == "" {
			t.Errorf("Alert %d should have an ID", i)
		}
		
		if alert.Message == "" {
			t.Errorf("Alert %d should have a message", i)
		}
		
		if alert.Timestamp.IsZero() {
			t.Errorf("Alert %d should have a timestamp", i)
		}
		
		if alert.Acknowledged {
			t.Errorf("Alert %d should not be acknowledged initially", i)
		}
	}
	
	// Test acknowledgment
	if len(alerts2) > 0 {
		alertID := alerts2[0].ID
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
} 