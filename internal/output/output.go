package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/zhangwlhaut/yapi-ai-cli/internal/model"
)

// Exit codes
const (
	ExitSuccess   = 0
	ExitParamError = 1
	ExitAuthError  = 2
	ExitAPIError   = 3
)

// PrintJSON writes data as JSON. If compact is true, output is minified.
func PrintJSON(w io.Writer, data interface{}, compact bool) error {
	var bytes []byte
	var err error
	if compact {
		bytes, err = json.Marshal(data)
	} else {
		bytes, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(bytes))
	return err
}

// PrintError writes a structured error to stderr.
func PrintError(w io.Writer, code, message, hint string) {
	err := struct {
		Error  string `json:"error"`
		Code   string `json:"code"`
		Hint   string `json:"hint,omitempty"`
	}{
		Error: message,
		Code:  code,
		Hint:  hint,
	}
	data, _ := json.Marshal(err)
	fmt.Fprintln(w, string(data))
}

// PrintProjectMarkdown writes a project in Markdown format.
func PrintProjectMarkdown(w io.Writer, p *model.Project) error {
	fmt.Fprintf(w, "# %s\n\n", p.Name)
	fmt.Fprintf(w, "- **ID**: %d\n", p.ID)
	fmt.Fprintf(w, "- **Base Path**: %s\n", p.BasePath)
	if p.Desc != "" {
		fmt.Fprintf(w, "- **Description**: %s\n", p.Desc)
	}
	if p.GroupName != "" {
		fmt.Fprintf(w, "- **Group**: %s (%d)\n", p.GroupName, p.GroupID)
	} else {
		fmt.Fprintf(w, "- **Group ID**: %d\n", p.GroupID)
	}
	return nil
}

// PrintInterfaceMarkdown writes an interface detail in Markdown format.
func PrintInterfaceMarkdown(w io.Writer, iface *model.Interface) error {
	fmt.Fprintf(w, "# %s\n\n", iface.Title)
	fmt.Fprintf(w, "- **Method**: %s\n", iface.Method)
	fmt.Fprintf(w, "- **Path**: %s\n", iface.Path)
	fmt.Fprintf(w, "- **Status**: %s\n", iface.Status)
	fmt.Fprintf(w, "- **Category ID**: %d\n", iface.CatID)

	if iface.Desc != "" {
		desc := iface.Desc
		// Strip HTML tags for markdown output
		desc = stripHTMLTags(desc)
		if desc != "" {
			fmt.Fprintf(w, "\n## Description\n\n%s\n", desc)
		}
	}

	// Request Headers
	if len(iface.ReqHeaders) > 0 {
		fmt.Fprintf(w, "\n## Request Headers\n\n")
		fmt.Fprintf(w, "| Name | Value | Required |\n")
		fmt.Fprintf(w, "|------|-------|----------|\n")
		for _, h := range iface.ReqHeaders {
			fmt.Fprintf(w, "| %s | %s | %s |\n", h.Name, h.Example, requiredLabel(h.Required))
		}
	}

	// Query Parameters
	if len(iface.ReqQuery) > 0 {
		fmt.Fprintf(w, "\n## Query Parameters\n\n")
		fmt.Fprintf(w, "| Name | Type | Required | Description | Example |\n")
		fmt.Fprintf(w, "|------|------|----------|-------------|--------|\n")
		for _, q := range iface.ReqQuery {
			fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n", q.Name, q.Type, requiredLabel(q.Required), q.Desc, q.Example)
		}
	}

	// Path Parameters
	if len(iface.ReqParams) > 0 {
		fmt.Fprintf(w, "\n## Path Parameters\n\n")
		fmt.Fprintf(w, "| Name | Description | Example |\n")
		fmt.Fprintf(w, "|------|-------------|--------|\n")
		for _, p := range iface.ReqParams {
			fmt.Fprintf(w, "| %s | %s | %s |\n", p.Name, p.Desc, p.Example)
		}
	}

	// Request Body (form)
	if len(iface.ReqBodyForm) > 0 {
		fmt.Fprintf(w, "\n## Request Body (form)\n\n")
		fmt.Fprintf(w, "| Name | Type | Required | Description | Example |\n")
		fmt.Fprintf(w, "|------|------|----------|-------------|--------|\n")
		for _, f := range iface.ReqBodyForm {
			fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n", f.Name, f.Type, requiredLabel(f.Required), f.Desc, f.Example)
		}
	}

	// Request Body (json/raw)
	if iface.ReqBodyOther != "" {
		fmt.Fprintf(w, "\n## Request Body (%s)\n\n", iface.ReqBodyType)
		fmt.Fprintf(w, "```%s\n%s\n```\n", bodyLang(iface.ReqBodyType), iface.ReqBodyOther)
	}

	// Response
	if iface.ResBody != "" {
		fmt.Fprintf(w, "\n## Response (%s)\n\n", iface.ResBodyType)
		fmt.Fprintf(w, "```%s\n%s\n```\n", bodyLang(iface.ResBodyType), iface.ResBody)
	}

	return nil
}

// PrintCategoryListMarkdown writes category list in Markdown format.
func PrintCategoryListMarkdown(w io.Writer, cats []model.CatMenu) error {
	for _, cat := range cats {
		fmt.Fprintf(w, "## %s (id: %d)\n\n", cat.Name, cat.ID)
		if cat.Desc != "" {
			fmt.Fprintf(w, "%s\n\n", cat.Desc)
		}
		if len(cat.List) > 0 {
			fmt.Fprintf(w, "| ID | Title | Path | Method | Status |\n")
			fmt.Fprintf(w, "|----|-------|------|--------|--------|\n")
			for _, iface := range cat.List {
				fmt.Fprintf(w, "| %d | %s | %s | %s | %s |\n", iface.ID, iface.Title, iface.Path, iface.Method, iface.Status)
			}
			fmt.Fprintln(w)
		}
	}
	return nil
}

// ConfirmDuplication prints a warning about a matching interface and prompts for confirmation.
// Returns true if the user confirms the update.
func ConfirmDuplication(iface *model.Interface) bool {
	fmt.Fprintf(os.Stderr, "⚠ Found existing interface with same path+method:\n")
	fmt.Fprintf(os.Stderr, "  ID: %d | %s (%s) | CatID: %d | Status: %s\n", iface.ID, iface.Path, iface.Method, iface.CatID, iface.Status)
	fmt.Fprintf(os.Stderr, "\nUpdate this interface? [y/N]: ")

	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

func requiredLabel(r string) string {
	if r == "1" {
		return "Yes"
	}
	return "No"
}

func bodyLang(t string) string {
	switch t {
	case "json":
		return "json"
	case "raw":
		return "text"
	default:
		return ""
	}
}

func stripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}
