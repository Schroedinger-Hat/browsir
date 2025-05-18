package browsir

import (
	"testing"

	"github.com/404answernotfound/browsir/utils"
)

// MockLoadLinks returns a predefined map of links to test
func MockLoadLinks() (map[string]string, error) {
	return map[string]string{
		"https://example.com": "general",
		"https://test.com":    "test",
	}, nil
}

// MockLoadLocalShortcuts returns a predefined map of shortcuts to test
func MockLoadLocalShortcuts() (map[string]string, error) {
	return map[string]string{
		"shortcut1": "https://shortcut1.com",
		"shortcut2": "https://shortcut2.com",
	}, nil
}

func TestListDefault(t *testing.T) {
	// Save original functions and restore them after the test
	originalLoadLinks := utils.LoadLinks
	originalLoadLocalShortcuts := utils.LoadLocalShortcuts
	defer func() {
		utils.LoadLinks = originalLoadLinks
		utils.LoadLocalShortcuts = originalLoadLocalShortcuts
	}()

	// Replace with mocks
	utils.LoadLinks = MockLoadLinks
	utils.LoadLocalShortcuts = MockLoadLocalShortcuts

	command := Command{}
	err := command.list([]string{})
	if err != nil {
		t.Errorf("list() returned an error: %v", err)
	}
}

func TestListLinks(t *testing.T) {
	// Save original func and restore it after the test
	originalLoadLinks := utils.LoadLinks
	defer func() {
		utils.LoadLinks = originalLoadLinks
	}()

	// Replace with mokc
	utils.LoadLinks = MockLoadLinks

	command := Command{}
	err := command.list([]string{"links"})
	if err != nil {
		t.Errorf("list() returned an error: %v", err)
	}
} 