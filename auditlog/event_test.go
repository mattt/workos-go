package auditlog

import "testing"

func TestNewEvent(t *testing.T) {
	event := NewEvent("user.login", Create)
	if event.Action != "user.login" {
		t.Errorf("got %q, wanted user.login", event.Action)
	}

	if event.ActionType != Create {
		t.Errorf("got %q, wanted %q", event.ActionType, Create)
	}

	if event.Location == "" {
		t.Error("event.Location should not be empty")
	}
}

func TestNewEventWithMetadata(t *testing.T) {
	event := NewEventWithMetadata("user.login", Create, map[string]interface{}{
		"user_agent": "Mozilla",
	})

	if len(event.Metadata) != 1 {
		t.Error("event.Metadata should contain one entry")
	}

	if event.Metadata["user_agent"] != "Mozilla" {
		t.Error("event.Metadata should contain one entry")
	}
}

func TestEventAddingMetadata(t *testing.T) {
	event := NewEvent("user.login", Create)
	if len(event.Metadata) != 0 {
		t.Error("event.Metadata should be empty")
	}

	event.AddMetadata("key", "value")

	if len(event.Metadata) != 1 {
		t.Error("event.Metadata should contain key/value")
	}

	if event.Metadata["key"] != "value" {
		t.Error("event.Metadata[key] should equal value")
	}
}

func TestEventAddingMetadataLimit(t *testing.T) {
	event := NewEvent("user.login", Create)

	buffer := make([]int, 500)
	for i := range buffer {
		event.AddMetadata(string(i), string(i))
	}

	event.AddMetadata("key", "value")

	if event.Metadata["key"] != nil {
		t.Error("event.Metadata should not have added entries for key/value")
	}

	if len(event.Metadata) != 500 {
		t.Error("event.Metadata should not contain over 500 entries")
	}
}

func TestEventPublishMergesGlobalMetadata(t *testing.T) {
	SetMetadata(map[string]interface{}{
		"environment": "testing",
	})

	event := NewEvent("user.login", Create)

	if event.Metadata["environment"] == "testing" {
		t.Errorf("expected event to not have any metadata for environment, got %q", event.Metadata["environment"])
	}

	event.Publish()

	if event.Metadata["environment"] != "testing" {
		t.Errorf("expected event to have metadata for environment, got %q", event.Metadata["environment"])
	}
}
func TestEventGlobalMetadataOverridesLocalMetadata(t *testing.T) {
	SetMetadata(map[string]interface{}{
		"environment": "testing",
	})

	event := NewEventWithMetadata("user.login", Create, map[string]interface{}{
		"environment": "production",
	})

	if event.Metadata["environment"] != "production" {
		t.Errorf("expected event to have metadata for environment set to production, got %q", event.Metadata["environment"])
	}

	event.Publish()

	if event.Metadata["environment"] != "testing" {
		t.Errorf("expected event to have had metadata overrwritten for environment, got %q", event.Metadata["environment"])
	}
}
