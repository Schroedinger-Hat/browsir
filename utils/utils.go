package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/404answernotfound/browsir/config"
	"github.com/PuerkitoBio/goquery"
)

func LoadLocalShortcuts() map[string]string {
	shortcuts := make(map[string]string)

	base, err := config.GetConfigBase()
	if err != nil {
		return shortcuts
	}

	// First check system config directory
	shortcutsPath := base + "/shortcuts"

	// Finally check current directory
	if _, err := os.Stat(shortcutsPath); os.IsNotExist(err) {
		shortcutsPath = "shortcuts"
	}

	data, err := os.ReadFile(shortcutsPath)
	if err != nil {
		return shortcuts
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			shortcuts[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return shortcuts
}

// TODO: Add other paths
func LoadLinks() map[string]string {
	var links = make(map[string]string)
	var linksPath string

	base, err := config.GetConfigBase()
	if err != nil {
		return links
	}

	// First check system config directory
	linksPath = base + "/links"

	if _, err := os.Stat(linksPath); os.IsNotExist(err) {
		linksPath = "links"
	}

	data, err := os.ReadFile(linksPath)
	if err != nil {
		return links
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			links[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return links
}

func SaveLink(link string, categories string) error {

	base, err := config.GetConfigBase()
	if err != nil {
		return err
	}

	linksPath := base + "/links"

	categoriesSlice := strings.Split(categories, ",")
	var categoriesToString string

	if len(categoriesSlice) == 0 {
		categoriesToString = "general"
	}

	categoriesToString = strings.Join(categoriesSlice, ",")

	f, err := os.OpenFile(linksPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	CheckDuplicates(f, link, 0, "link", "|")

	if err == nil {
		defer f.Close()
		// Writing to file https://somelink.com|some category
		_, err = fmt.Fprintf(f, "%s|%s\n", link, categoriesToString)
		return err
	}

	_, err = fmt.Fprintf(f, "%s|%s\n", link, categoriesToString)
	return err
}

func SaveLocalShortcut(shortcut, url string) error {
	base, err := config.GetConfigBase()
	if err != nil {
		return err
	}

	shortcutsPath := base + "/shortcuts"

	f, err := os.OpenFile(shortcutsPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	CheckDuplicates(f, url, 1, "shortcut", "=")

	if err == nil {
		defer f.Close()
		_, err = fmt.Fprintf(f, "%s=%s\n", shortcut, url)
		if err == nil {
			fmt.Printf("Shortcut %s correctly saved\n", shortcut)
		}
		return err
	}

	return err
}

func RemoveLocalShortcut(shortcut string) error {
	base, err := config.GetConfigBase()
	if err != nil {
		return err
	}

	shortcutsPath := base + "/shortcuts"

	// First open shortcut file
	f, err := os.OpenFile(shortcutsPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a temp empty file
	tempFilePath := shortcutsPath + ".tmp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(f)
	writer := bufio.NewWriter(tempFile)

	// Loop the shortcut file and copy in the temp only the ones that don't match the requested one
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, shortcut+"=") {
			found = true
			continue
		}
		fmt.Fprintln(writer, line)
	}

	// if not found returns related error
	if !found {
		return fmt.Errorf("shortcut '%v' not found", shortcut)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	writer.Flush()

	// Replace the temp file as the new shortcut file
	if err := os.Rename(tempFilePath, shortcutsPath); err != nil {
		return err
	}

	err = os.Remove(tempFilePath)

	if err != nil {
		return err
	}

	fmt.Printf("Shortcut %s correctly removed!\n", shortcut)
	return nil
}

func GetBrowserPath(browserName string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		switch browserName {
		case "chrome":
			return "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", nil
		case "brave":
			return "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser", nil
		case "arc":
			return "/Applications/Arc.app/Contents/MacOS/Arc", nil
		case "zen":
			return "/Applications/Zen.app/Contents/MacOS/zen", nil
		case "firefox":
			return "/Applications/Firefox.app/Contents/MacOS/firefox", nil
		case "firefox-developer-edition":
			return "/Applications/Firefox Developer Edition.app/Contents/MacOS/firefox", nil
		}
	case "linux":
		switch browserName {
		case "firefox":
			return "firefox", nil
		case "vivaldi":
			return "vivaldi", nil
		case "chrome":
			return "google-chrome", nil
		case "brave":
			return "brave-browser", nil
		}
	case "windows":
		switch browserName {
		case "chrome":
			return filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"), nil
		case "brave":
			return filepath.Join(
				os.Getenv("ProgramFiles"),
				"BraveSoftware",
				"Brave-Browser",
				"Application",
				"brave.exe",
			), nil
		}
	}
	return "", fmt.Errorf("unsupported browser %s on %s", browserName, runtime.GOOS)
}

func OpenBrowser(browserName string, profile config.Profile, url string) error {
	browserPath, err := GetBrowserPath(browserName)
	if err != nil {
		return err
	}

	var args []string
	if browserName == "firefox" || browserName == "firefox-developer-edition" || browserName == "zen" {
		args = []string{"-profile", profile.ProfileDir}
	} else {
		args = []string{"--profile-directory=" + profile.ProfileDir}
	}

	if url != "" {
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}
		args = append(args, url)
	}

	cmd := exec.Command(browserPath, args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start browser: %v", err)
	}

	return nil
}

func PrintUsage(profiles []config.Profile, shortcuts map[string]string, localShortcuts map[string]string) {
	fmt.Printf("browsir v1.0.0\n\n")
	fmt.Println("Usage: browsir [profile] [url|shortcut]")
	fmt.Println("\nProfiles:")
	for _, p := range profiles {
		fmt.Printf("  %-12s - %s\n", p.Name, p.Description)
	}
	fmt.Println("\nShortcuts:")
	for shortcut, url := range shortcuts {
		fmt.Printf("  %-12s -> %s\n", shortcut, url)
	}
	if len(localShortcuts) > 0 {
		fmt.Println("\nLocal Shortcuts:")
		for shortcut, url := range localShortcuts {
			fmt.Printf("  %-12s -> %s\n", shortcut, url)
		}
	}
	fmt.Println("\nExamples:")
	fmt.Println("  browsir work                    # Open browser with work profile")
	fmt.Println("  browsir personal gmail.com      # Open Gmail with personal profile")
	fmt.Println("  browsir default mail           # Open mail shortcut with default profile")

	fmt.Println("\nOther commands:")
	fmt.Println("  -h, --help            # Print this help message")
	fmt.Println("  -v, --version         # Print browsir version")
	fmt.Println("  -ls, --list-shortcuts # List all shortcuts")
	fmt.Println("  -p, --profiles        # List all profiles")
	fmt.Println("  -q        		     # Search the web with a query")
	fmt.Println("  -se, --search-engine  # Specify search engine (google, duckduckgo, brave)")

	fmt.Println("   browsir add link <link> -c <categories>	# Add a link with categories")
	fmt.Println("   browsir add shortcut <shortcut> <url>	# Add a local shortcut, do not include http:// or https://")
	fmt.Println("   browsir rm link <link>					# Remove a link")
	fmt.Println("   browsir rm shortcut <shortcut>			# Remove a local shortcut")
	fmt.Println("   browsir list links						# List all links")
	fmt.Println("   browsir list all						# List all links and categories")
	fmt.Println("   browsir preview <link>					# Preview a link")
}

func PrintProfiles(profiles []config.Profile) {
	fmt.Println("\nProfiles:")
	for _, p := range profiles {
		fmt.Printf("  %-12s - %s\n", p.Name, p.Description)
	}
}

func FindSimilarShortcuts(input string, shortcuts, localShortcuts map[string]string) []string {
	var similar []string
	for shortcut := range shortcuts {
		if strings.Contains(shortcut, input) || strings.Contains(input, shortcut) {
			similar = append(similar, shortcut)
		}
	}
	for shortcut := range localShortcuts {
		if strings.Contains(shortcut, input) || strings.Contains(input, shortcut) {
			similar = append(similar, shortcut)
		}
	}
	return similar
}

func PromptYesNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s (y/n): ", prompt)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		}
		if response == "n" || response == "no" {
			return false
		}
	}
}

func PrintLocalShortcuts(shortcuts map[string]string) {
	for shortcut, url := range shortcuts {
		fmt.Printf("  %-12s -> %s\n", shortcut, url)
	}
}

func Search(browserName string, profile config.Profile, searchEngine string, searchTerm string) error {

	var url string
	if searchTerm != "" {
		switch searchEngine {
		case "google":
			url = fmt.Sprintf("google.com/search?q=%s", searchTerm)
		case "duckduckgo":
			url = fmt.Sprintf("duckduckgo.com//?q=%s", searchTerm)
		case "brave":
			url = fmt.Sprintf("search.brave.com/search?q=%s", searchTerm)
		default:
			url = fmt.Sprintf("google.com/search?q=%s", searchTerm)
		}
		url = "https://" + url
	}

	OpenBrowser(browserName, profile, url)

	return nil
}

func GetFlags(args []string) map[string]string {
	flags := make(map[string]string)
	if args == nil {
		return flags
	}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				flags[parts[0]] = parts[1]
			} else {
				flags[arg] = ""
			}
		}
	}
	return flags
}

func Contains(args []string, value string) bool {
	for _, arg := range args {
		if arg == value {
			return true
		}
	}
	return false
}

func Log(args ...any) {
	for _, arg := range args {
		fmt.Println(arg)
	}
}

func ExitLog(exitMessage string) {
	fmt.Println(exitMessage)
	os.Exit(0)
}

func CheckDuplicates(f io.Reader, url string, pos int, t string, separator string) {
	buf := bufio.NewScanner(f)
	for buf.Scan() {
		line := buf.Text()
		splittedLine := strings.Split(line, separator)
		if len(splittedLine) > 1 && splittedLine[pos] == url {
			fmt.Printf("%s already exists with url %s", t, splittedLine[0])
			os.Exit(0)
		}
	}
}

func FindHtmlNode(doc *goquery.Document, search []string) []string {
	values := make([]string, 0)
	for _, term := range search {
		values = append(values, doc.Find(term).Text())
	}
	return values
}

func CheckInputArgs(currentArgs, expectedArgs int) error {
	if currentArgs < expectedArgs {
		fmt.Printf("You provided %d arguments, while %d are needed.\nPlease, see --help flag to check usage for add command\n", currentArgs, expectedArgs)
		return errors.New("not enough arguments")
	}
	return nil
}
