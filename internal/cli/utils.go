package cli

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

func RenderMarkdown(content string) error {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(150),
	)
	if err != nil {
		return err
	}

	out, err := r.Render(content)
	if err != nil {
		return err
	}

	fmt.Print(out)
	return nil
}
