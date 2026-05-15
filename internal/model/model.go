package model

import "encoding/json"

// APIResponse is the common YApi API response wrapper.
type APIResponse struct {
	ErrCode int             `json:"errcode"`
	ErrMsg  string           `json:"errmsg"`
	Data    json.RawMessage `json:"data"`
}

// Project represents a YApi project.
type Project struct {
	ID        int    `json:"_id"`
	Name      string `json:"name"`
	BasePath  string `json:"basepath"`
	Desc      string `json:"desc"`
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	UID       int    `json:"uid"`
	UserName  string `json:"username"`
	AddTime   int64  `json:"add_time"`
	UpTime    int64  `json:"up_time"`
}

// Category represents an interface category (分类).
type Category struct {
	ID        int    `json:"_id"`
	Name      string `json:"name"`
	ProjectID int    `json:"project_id"`
	Desc      string `json:"desc"`
	UID       int    `json:"uid"`
	AddTime   int64  `json:"add_time"`
	UpTime    int64  `json:"up_time"`
	List      []Interface `json:"list,omitempty"`
}

// CatMenu represents a category menu item.
type CatMenu struct {
	ID        int    `json:"_id"`
	Name      string `json:"name"`
	ProjectID int    `json:"project_id"`
	Desc      string `json:"desc"`
	UID       int    `json:"uid"`
	AddTime   int64  `json:"add_time"`
	UpTime    int64  `json:"up_time"`
	List      []Interface `json:"list,omitempty"`
}

// Interface represents a YApi interface definition.
type Interface struct {
	ID                int              `json:"_id"`
	ProjectID         int              `json:"project_id"`
	CatID             int              `json:"catid"`
	Title             string           `json:"title"`
	Path              string           `json:"path"`
	Method            string           `json:"method"`
	ReqBodyType       string           `json:"req_body_type"`
	ReqBodyOther      string           `json:"req_body_other,omitempty"`
	ResBody           string           `json:"res_body"`
	ResBodyType       string           `json:"res_body_type"`
	ResBodyIsJSONSchema bool           `json:"res_body_is_json_schema"`
	UID               int              `json:"uid"`
	UserName          string           `json:"username,omitempty"`
	AddTime           int64            `json:"add_time"`
	UpTime            int64            `json:"up_time"`
	Status            string           `json:"status"`
	EditUID           int              `json:"edit_uid"`
	Desc              string           `json:"desc,omitempty"`
	ReqBodyForm       []ReqBodyItem    `json:"req_body_form,omitempty"`
	ReqParams         []ReqParamItem   `json:"req_params,omitempty"`
	ReqHeaders        []ReqHeaderItem  `json:"req_headers,omitempty"`
	ReqQuery          []ReqQueryItem   `json:"req_query,omitempty"`
	SwitchNotice      bool             `json:"switch_notice,omitempty"`
	Message           string           `json:"message,omitempty"`
}

// ReqBodyItem represents a form body parameter.
type ReqBodyItem struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Example  string `json:"example"`
	Desc     string `json:"desc"`
	Required string `json:"required"`
}

// ReqParamItem represents a path parameter.
type ReqParamItem struct {
	Name    string `json:"name"`
	Example string `json:"example"`
	Desc    string `json:"desc"`
}

// ReqHeaderItem represents a request header.
type ReqHeaderItem struct {
	Name     string `json:"name"`
	Type     string `json:"type,omitempty"`
	Example  string `json:"example,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Required string `json:"required"`
}

// ReqQueryItem represents a query parameter.
type ReqQueryItem struct {
	Name     string `json:"name"`
	Type     string `json:"type,omitempty"`
	Example  string `json:"example,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Required string `json:"required"`
}

// InterfaceListSummary is the compact list item for interface list output.
type InterfaceListSummary struct {
	ID     int    `json:"_id"`
	Title  string `json:"title"`
	Path   string `json:"path"`
	Method string `json:"method"`
	Status string `json:"status"`
	CatID  int    `json:"catid,omitempty"`
}

// AddCatRequest is the request body for add_cat.
type AddCatRequest struct {
	Name      string `json:"name"`
	ProjectID int    `json:"project_id"`
	Desc      string `json:"desc,omitempty"`
	Token     string `json:"token"`
}

// AddInterfaceRequest is the request body for interface add/up/save.
type AddInterfaceRequest struct {
	Token        string          `json:"token"`
	ProjectID    int             `json:"project_id,omitempty"`
	CatID        int             `json:"catid"`
	Title        string          `json:"title"`
	Path         string          `json:"path"`
	Method       string          `json:"method"`
	Desc         string          `json:"desc,omitempty"`
	ReqBodyType  string          `json:"req_body_type,omitempty"`
	ReqBodyOther string          `json:"req_body_other,omitempty"`
	ResBodyType  string          `json:"res_body_type,omitempty"`
	ResBody      string          `json:"res_body,omitempty"`
	ReqQuery     []ReqQueryItem  `json:"req_query,omitempty"`
	ReqHeaders   []ReqHeaderItem `json:"req_headers,omitempty"`
	ReqBodyForm  []ReqBodyItem   `json:"req_body_form,omitempty"`
	ReqParams    []ReqParamItem  `json:"req_params,omitempty"`
	Status       string          `json:"status,omitempty"`
	SwitchNotice bool            `json:"switch_notice,omitempty"`
	Message      string          `json:"message,omitempty"`
	ID           int             `json:"id,omitempty"` // for update/save
}

// ImportRequest is the request body for import_data.
type ImportRequest struct {
	Token string `json:"token"`
	Type  string `json:"type"`
	JSON  string `json:"json,omitempty"`
	URL   string `json:"url,omitempty"`
	Merge string `json:"merge"`
}

// WriteResult represents a create/update result.
type WriteResult struct {
	OK        int `json:"ok"`
	NModified int `json:"nModified"`
	N         int `json:"n"`
}
