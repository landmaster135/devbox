package usecases

import (
	"testing"
)

func TestNewService(t *testing.T) {
	project := "test-project"
	location := "us-central1"
	serviceName := "test-service"
	serviceAccountID := "test-sa"

	service := NewService(project, location, serviceName, serviceAccountID)

	if service.project != project {
		t.Errorf("Expected project %s, got %s", project, service.project)
	}
	if service.location != location {
		t.Errorf("Expected location %s, got %s", location, service.location)
	}
	if service.serviceName != serviceName {
		t.Errorf("Expected serviceName %s, got %s", serviceName, service.serviceName)
	}
	if service.serviceAccountID != serviceAccountID {
		t.Errorf("Expected serviceAccountID %s, got %s", serviceAccountID, service.serviceAccountID)
	}
}

func TestBuildDashboardConfig(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "")

	config := service.buildDashboardConfig()

	if config.DisplayName != "Cloud Run Monitoring - test-service" {
		t.Errorf("Expected display name 'Cloud Run Monitoring - test-service', got %s", config.DisplayName)
	}

	if config.Layout == nil {
		t.Error("Expected layout to be set")
	}

	gridLayout := config.GetGridLayout()
	if gridLayout == nil {
		t.Error("Expected grid layout to be set")
	}

	if gridLayout.Columns != 12 {
		t.Errorf("Expected 12 columns, got %d", gridLayout.Columns)
	}

	// ウィジェット数の確認（16個のウィジェット）
	expectedWidgetCount := 16
	if len(gridLayout.Widgets) != expectedWidgetCount {
		t.Errorf("Expected %d widgets, got %d", expectedWidgetCount, len(gridLayout.Widgets))
	}
}

func TestCreateTextWidget(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "")

	title := "Test Widget"
	content := "Test content"

	widget := service.createTextWidget(title, content)

	if widget.Title != title {
		t.Errorf("Expected title %s, got %s", title, widget.Title)
	}

	textWidget := widget.GetText()
	if textWidget == nil {
		t.Error("Expected text widget to be set")
	}

	expectedContent := "**Test Widget**\n\nTest content\n\nサービス: test-service"
	if textWidget.Content != expectedContent {
		t.Errorf("Expected content %s, got %s", expectedContent, textWidget.Content)
	}
}
