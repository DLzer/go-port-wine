package port

func ScanPort(protocol, hostname string, port int) bool {
	address := hostname + ":" strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)
}