package browsir

import (
	"context"
	"fmt"
	"os"
	"strings"

	"io"
	"net/http"
	"time"

	"github.com/404answernotfound/browsir/utils"
	"github.com/PuerkitoBio/goquery"
)

type ICommand interface {
	add(args []string)
	rm(args []string)
	last(args []string)
	preview(args []string)
}

type Command struct{}

func (c Command) add(args []string) error {
	switch args[0] {
	case "link":
		err := utils.CheckInputArgs(len(args), 2)
		if err != nil {
			os.Exit(0)
		}

		flags := utils.GetFlags(args[2:])
		link := args[1]

		categories, ok := flags["-c"]
		if !ok {
			return fmt.Errorf("not a good flag")
		}

		err = utils.SaveLink(link, categories)
		if err != nil {
			return err
		}
		return nil
	case "shortcut":
		err := utils.CheckInputArgs(len(args), 3)
		if err != nil {
			os.Exit(0)
		}

		shortcut := args[1]
		url := args[2]
		err = utils.SaveLocalShortcut(shortcut, url)
		return err
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (c Command) remove(args []string) error {
	switch args[0] {
	case "link":
		fmt.Println("rm link is not yet implemented")
	case "shortcut":
		err := utils.CheckInputArgs(len(args), 2)
		if err != nil {
			os.Exit(0)
		}

		shortcut := args[1]
		err = utils.RemoveLocalShortcut(shortcut)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Command) list(args []string) error {
	if len(args) == 0 || (len(args) > 0 && (args[0] == "all")) {
		// List all links and all shortcuts
		links, err := utils.LoadLinks()
		if err == nil {
			fmt.Println("Links:")
			for link, categories := range links {
				fmt.Printf("  %s - Categories: %s\n", link, categories)
			}
		} else {
			fmt.Fprintf(os.Stderr, "There was an issue loading your links configuration, %v\n", err)
		}
		shortcuts, err := utils.LoadLocalShortcuts()
		if err == nil {
			fmt.Println("Shortcuts:")
			utils.PrintLocalShortcuts(shortcuts)
		} else {
			fmt.Fprintf(os.Stderr, "There was an issue loading your shortcuts configuration, %v\n", err)
		}
		return nil
	}

	if args[0] == "links" {
		links, err := utils.LoadLinks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "There was an issue loading your links configuration, %v\n", err)
			return nil
		}
		fmt.Println("Links:")
		for link, categories := range links {
			fmt.Printf("  %s - Categories: %s\n", link, categories)
		}
		return nil
	}

	if args[0] == "shortcuts" {
		shortcuts, err := utils.LoadLocalShortcuts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "There was an issue loading your shortcuts configuration, %v\n", err)
			return nil
		}
		fmt.Println("Shortcuts:")
		utils.PrintLocalShortcuts(shortcuts)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Unknown argument for list: %s\n", args[0])
	return nil
}

func (c Command) preview(args []string) error {
	ctx := context.Background()
	deadline := time.Now().Add(3000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", args[0], nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))

	title := doc.Find("title").Text()
	fmt.Printf("Title: %s\n", title)

	description, exists := doc.Find("meta[name='description']").Attr("content")
	if exists {
		fmt.Printf("Description: %s\n", description)
	} else {
		fmt.Println("No meta description found.")
	}

	allH1Tags := doc.Find("h1").Text()
	fmt.Printf("H1 Tags: %v\n", allH1Tags)

	return nil
}

func RunCommand(mainCmd string, otherArgs []string) error {
	command := &Command{}
	var err error

	switch mainCmd {
	case "add":
		err = command.add(otherArgs)
		return err
	case "rm":
		err = command.remove(otherArgs)
		return err
	case "list":
		err = command.list(otherArgs)
		return err
	case "preview":
		err = command.preview(otherArgs)
		return err
	default:
		return fmt.Errorf("not implemented")
	}
}
