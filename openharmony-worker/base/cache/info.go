package cache

// AddressInfo 地址信息
type AddressInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// ClientInfo 客户端信息
type ClientInfo struct {
	Address  []AddressInfo `json:"address"`
	PassWord string        `json:"password"`
}

// DistributedLockInfo 分布式锁信息
type DistributedLockInfo struct {
	TimeOut     int `json:"timeout"`
	MaxRetry    int `json:"maxretry"`
	RIntervalMs int `json:"rintervalms"`
}

// CInfo 缓存信息
type CInfo struct {
	Client          ClientInfo          `json:"client"`
	DistributedLock DistributedLockInfo `json:"distributedlock"`
}
