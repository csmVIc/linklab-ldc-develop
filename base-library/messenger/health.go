package messenger

// GetClosed 检查nats是否关闭
func (md *Driver) GetClosed() bool {
	return md.closed
}
