package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zhangwlhaut/ypi-ai-cli/internal/model"
)

// Client is the YApi Open API HTTP client.
type Client struct {
	Server     string
	Token      string
	ProjectID  int
	httpClient *http.Client
}

// New creates a new YApi client.
func New(server, token string, projectID int) *Client {
	return &Client{
		Server:     server,
		Token:      token,
		ProjectID:  projectID,
		httpClient: &http.Client{},
	}
}

// GetProject fetches project basic info.
func (c *Client) GetProject() (*model.Project, error) {
	resp, err := c.get("/api/project/get", url.Values{"token": {c.Token}})
	if err != nil {
		return nil, err
	}
	var project model.Project
	if err := json.Unmarshal(resp, &project); err != nil {
		return nil, fmt.Errorf("unmarshal project: %w", err)
	}
	return &project, nil
}

// GetCatMenu fetches category menu list.
func (c *Client) GetCatMenu() ([]model.CatMenu, error) {
	resp, err := c.get("/api/interface/getCatMenu", url.Values{
		"project_id": {strconv.Itoa(c.ProjectID)},
		"token":      {c.Token},
	})
	if err != nil {
		return nil, err
	}
	var cats []model.CatMenu
	if err := json.Unmarshal(resp, &cats); err != nil {
		return nil, fmt.Errorf("unmarshal cat menu: %w", err)
	}
	return cats, nil
}

// AddCat creates a new interface category.
func (c *Client) AddCat(req *model.AddCatRequest) (json.RawMessage, error) {
	req.Token = c.Token
	req.ProjectID = c.ProjectID
	return c.postJSON("/api/interface/add_cat", req)
}

// ListInterfaces fetches all interfaces in the project (paginated).
func (c *Client) ListInterfaces(page, limit int) ([]model.Interface, error) {
	if limit <= 0 {
		limit = 1000
	}
	if page <= 0 {
		page = 1
	}
	resp, err := c.get("/api/interface/list", url.Values{
		"project_id": {strconv.Itoa(c.ProjectID)},
		"token":      {c.Token},
		"page":       {strconv.Itoa(page)},
		"limit":      {strconv.Itoa(limit)},
	})
	if err != nil {
		return nil, err
	}
	var result struct {
		List  []model.Interface `json:"list"`
		Total int               `json:"total"`
		Count int               `json:"count"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("unmarshal interface list: %w", err)
	}
	return result.List, nil
}

// ListInterfacesByCat fetches interfaces under a category (paginated).
func (c *Client) ListInterfacesByCat(catid, page, limit int) ([]model.Interface, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	resp, err := c.get("/api/interface/list_cat", url.Values{
		"token": {c.Token},
		"catid": {strconv.Itoa(catid)},
		"page":  {strconv.Itoa(page)},
		"limit": {strconv.Itoa(limit)},
	})
	if err != nil {
		return nil, err
	}
	var result struct {
		List  []model.Interface `json:"list"`
		Total int               `json:"total"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("unmarshal interface list by cat: %w", err)
	}
	return result.List, nil
}

// GetInterface fetches a single interface detail by ID.
func (c *Client) GetInterface(id int) (*model.Interface, error) {
	resp, err := c.get("/api/interface/get", url.Values{
		"id":    {strconv.Itoa(id)},
		"token": {c.Token},
	})
	if err != nil {
		return nil, err
	}
	var iface model.Interface
	if err := json.Unmarshal(resp, &iface); err != nil {
		return nil, fmt.Errorf("unmarshal interface: %w", err)
	}
	return &iface, nil
}

// AddInterface creates a new interface.
func (c *Client) AddInterface(req *model.AddInterfaceRequest) (json.RawMessage, error) {
	req.Token = c.Token
	if req.ProjectID == 0 {
		req.ProjectID = c.ProjectID
	}
	return c.postJSON("/api/interface/add", req)
}

// UpdateInterface updates an existing interface.
func (c *Client) UpdateInterface(req *model.AddInterfaceRequest) (json.RawMessage, error) {
	req.Token = c.Token
	if req.ProjectID == 0 {
		req.ProjectID = c.ProjectID
	}
	return c.postJSON("/api/interface/up", req)
}

// SaveInterface creates or updates an interface (upsert).
func (c *Client) SaveInterface(req *model.AddInterfaceRequest) (json.RawMessage, error) {
	req.Token = c.Token
	if req.ProjectID == 0 {
		req.ProjectID = c.ProjectID
	}
	return c.postJSON("/api/interface/save", req)
}

// ListMenu fetches the interface menu tree (categories with their interfaces).
func (c *Client) ListMenu() ([]model.CatMenu, error) {
	resp, err := c.get("/api/interface/list_menu", url.Values{
		"project_id": {strconv.Itoa(c.ProjectID)},
		"token":      {c.Token},
	})
	if err != nil {
		return nil, err
	}
	var cats []model.CatMenu
	if err := json.Unmarshal(resp, &cats); err != nil {
		return nil, fmt.Errorf("unmarshal interface menu: %w", err)
	}
	return cats, nil
}

// ImportData imports data from swagger/json.
func (c *Client) ImportData(req *model.ImportRequest) (json.RawMessage, error) {
	req.Token = c.Token
	return c.postForm("/api/open/import_data", req)
}

// FindInterfaceByPathMethod searches all interfaces for one matching path+method.
func (c *Client) FindInterfaceByPathMethod(path, method string) (*model.Interface, error) {
	menu, err := c.ListMenu()
	if err != nil {
		return nil, err
	}
	for _, cat := range menu {
		for _, iface := range cat.List {
			if iface.Path == path && stringsEqualIgnoreCase(iface.Method, method) {
				return &iface, nil
			}
		}
	}
	return nil, nil
}

func (c *Client) get(path string, params url.Values) (json.RawMessage, error) {
	u := strings.TrimRight(c.Server, "/") + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

func (c *Client) postJSON(path string, body any) (json.RawMessage, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	u := strings.TrimRight(c.Server, "/") + path
	req, err := http.NewRequest("POST", u, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

func (c *Client) postForm(path string, body any) (json.RawMessage, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	u := strings.TrimRight(c.Server, "/") + path
	req, err := http.NewRequest("POST", u, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) (json.RawMessage, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp model.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if apiResp.ErrCode != 0 {
		return nil, fmt.Errorf("API error %d: %s", apiResp.ErrCode, apiResp.ErrMsg)
	}

	return apiResp.Data, nil
}

func stringsEqualIgnoreCase(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		la := a[i]
		lb := b[i]
		if la >= 'A' && la <= 'Z' {
			la += 32
		}
		if lb >= 'A' && lb <= 'Z' {
			lb += 32
		}
		if la != lb {
			return false
		}
	}
	return true
}
