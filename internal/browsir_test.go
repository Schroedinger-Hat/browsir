package browsir

import (
	"testing"

	"github.com/404answernotfound/browsir/utils"
)

// LinkLoader defines the interface for loading links and shortcuts
type LinkLoader interface {
	LoadLinks() (map[string]string, error)
	LoadLocalShortcuts() (map[string]string, error)
}

// MockLinkLoader implements LinkLoader interface for testing
type MockLinkLoader struct{}

func (m *MockLinkLoader) LoadLinks() (map[string]string, error) {
	return map[string]string{
		"https://example.com": "general",
		"https://test.com":    "test",
	}, nil
}

func (m *MockLinkLoader) LoadLocalShortcuts() (map[string]string, error) {
	return map[string]string{
		"shortcut1": "https://shortcut1.com",
		"shortcut2": "https://shortcut2.com",
	}, nil
}

// RealLinkLoader implements LinkLoader interface using the actual utils package
type RealLinkLoader struct{}

func (r *RealLinkLoader) LoadLinks() (map[string]string, error) {
	return utils.LoadLinks()
}

func (r *RealLinkLoader) LoadLocalShortcuts() (map[string]string, error) {
	return utils.LoadLocalShortcuts()
}

func TestListDefault(t *testing.T) {
	mockLoader := &MockLinkLoader{}
	command := Command{linkLoader: mockLoader}
	
	err := command.list([]string{})
	if err != nil {
		t.Errorf("list() returned an error: %v", err)
	}
}

func TestListLinks(t *testing.T) {
	mockLoader := &MockLinkLoader{}
	command := Command{linkLoader: mockLoader}
	
	err := command.list([]string{"links"})
	if err != nil {
		t.Errorf("list() returned an error: %v", err)
	}
} 