package saasbackend

// SInfo 配置信息
type SInfo struct {
	UserSite UserSiteInfo `json:"usersite"`
	Site     SiteInfo     `json:"siteinfo"`
}

// UserSiteInfo 用户租户信息
type UserSiteInfo struct {
	URL string `json:"url"`
}

// SiteInfo 租户信息
type SiteInfo struct {
	URL string `json:"url"`
}

// SaasResponse 回复报文
type SaasResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SiteResult 租户查询结果
type SiteResult struct {
	ID   int
	Name string
}
