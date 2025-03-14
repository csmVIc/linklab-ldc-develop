package boardcmd

// BoardCommand 设备操作接口
type BoardCommand interface {
	Burn(devport string, burnfile string, burncmd string) (string, error)
	BurnEmptyProgram(devport string, burncmd string) error
	SetEmptyProgram(emptyprogram string)
	Scan() (map[string]bool, error)
	SetScanCmd(scancmd string)
}
