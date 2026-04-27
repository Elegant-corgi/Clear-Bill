package bll

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type PermissionItem struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Tag    string `json:"tag"`
}

type PermissionCatalog struct {
	items   []PermissionItem
	byID    map[string]PermissionItem
	byRoute map[string]PermissionItem
}

func NewPermissionCatalog() *PermissionCatalog {
	items := loadPermissionsFromSwaggerAnnotations()
	byID := make(map[string]PermissionItem, len(items))
	byRoute := make(map[string]PermissionItem, len(items))
	for _, item := range items {
		byID[item.ID] = item
		byRoute[routePermissionKey(item.Method, item.Path)] = item
	}
	return &PermissionCatalog{
		items:   items,
		byID:    byID,
		byRoute: byRoute,
	}
}

func (c *PermissionCatalog) List() []PermissionItem {
	result := make([]PermissionItem, len(c.items))
	copy(result, c.items)
	return result
}

func (c *PermissionCatalog) Exists(permissionID string) bool {
	_, ok := c.byID[permissionID]
	return ok
}

func (c *PermissionCatalog) GetByID(permissionID string) (PermissionItem, bool) {
	item, ok := c.byID[permissionID]
	return item, ok
}

func (c *PermissionCatalog) GetByRoute(method, path string) (PermissionItem, bool) {
	item, ok := c.byRoute[routePermissionKey(method, normalizeSwaggerPath(path))]
	return item, ok
}

func loadPermissionsFromSwaggerAnnotations() []PermissionItem {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to resolve permission catalog source path")
	}

	actionDir := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "../action"))
	entries, err := os.ReadDir(actionDir)
	if err != nil {
		panic(err)
	}

	type partial struct {
		id     string
		method string
		path   string
		tag    string
	}

	items := make([]PermissionItem, 0, 32)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		file, err := os.Open(filepath.Join(actionDir, entry.Name()))
		if err != nil {
			panic(err)
		}

		var current partial
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			switch {
			case strings.HasPrefix(line, "// @Router"):
				raw := strings.TrimSpace(strings.TrimPrefix(line, "// @Router"))
				parts := strings.Fields(raw)
				if len(parts) >= 2 {
					current.path = normalizeSwaggerPath(parts[0])
					current.method = strings.ToUpper(strings.Trim(parts[1], "[]"))
				}
			case strings.HasPrefix(line, "// @ID"):
				current.id = strings.TrimSpace(strings.TrimPrefix(line, "// @ID"))
			case strings.HasPrefix(line, "// @Tags"):
				current.tag = strings.TrimSpace(strings.TrimPrefix(line, "// @Tags"))
			case strings.HasPrefix(line, "func ") && current.id != "" && current.path != "" && current.method != "":
				items = append(items, PermissionItem{
					ID:     current.id,
					Method: current.method,
					Path:   current.path,
					Tag:    current.tag,
				})
				current = partial{}
			}
		}
		_ = file.Close()
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Tag != items[j].Tag {
			return items[i].Tag < items[j].Tag
		}
		if items[i].Path != items[j].Path {
			return items[i].Path < items[j].Path
		}
		if items[i].Method != items[j].Method {
			return items[i].Method < items[j].Method
		}
		return items[i].ID < items[j].ID
	})
	return items
}

func routePermissionKey(method, path string) string {
	return strings.ToUpper(strings.TrimSpace(method)) + " " + normalizeSwaggerPath(path)
}

func normalizeSwaggerPath(path string) string {
	path = strings.TrimSpace(path)
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = "{" + strings.TrimPrefix(part, ":") + "}"
		}
	}
	return strings.Join(parts, "/")
}
