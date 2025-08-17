package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yuya-takeyama/panama/internal/workspace"
)

type Format string

const (
	FormatPath Format = "path"
	FormatCD   Format = "cd"
	FormatJSON Format = "json"
)

func Print(path string, format Format) error {
	switch format {
	case FormatPath:
		fmt.Println(path)
	case FormatCD:
		fmt.Printf("cd \"%s\"\n", path)
	case FormatJSON:
		data := map[string]string{"path": path}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(data); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
	return nil
}

func PrintWorkspaces(workspaces []*workspace.Workspace, format Format) error {
	switch format {
	case FormatPath:
		for _, ws := range workspaces {
			fmt.Println(ws.Path)
		}
	case FormatJSON:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(workspaces); err != nil {
			return err
		}
	default:
		return fmt.Errorf("format %s is not supported for listing workspaces", format)
	}
	return nil
}

func ParseFormat(s string) (Format, error) {
	switch s {
	case "path":
		return FormatPath, nil
	case "cd":
		return FormatCD, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("invalid format: %s", s)
	}
}