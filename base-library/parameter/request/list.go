package request

// DeviceListQuery 获取设备列表
type DeviceListQuery struct {
	BoardName string `form:"boardname" json:"boardname" binding:"required"`
}
