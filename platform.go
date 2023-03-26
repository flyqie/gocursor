package gocursor

// PlatformStatus 获取实例的平台状态
func (c *Cursor) PlatformStatus() interface{} {
	if !c.usable {
		return nil
	}
	return c.platImpl.Status()
}
