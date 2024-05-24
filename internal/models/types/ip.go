package types

import (
	"strconv"
	"strings"
)

const IP_SIZE uint64 = 4

type Ip string

func (ip Ip) Size() uint64 {
	return IP_SIZE
}

func (ip Ip) Marshal() [IP_SIZE]byte {
	octets := strings.Split(string(ip), ".")

	octet0, _ := strconv.Atoi(octets[0])
	octet1, _ := strconv.Atoi(octets[1])
	octet2, _ := strconv.Atoi(octets[2])
	octet3, _ := strconv.Atoi(octets[3])

	return [IP_SIZE]byte{byte(octet0), byte(octet1), byte(octet2), byte(octet3)}
}

func (ip *Ip) Unmarshal(bytes []byte) error {
	str := strconv.Itoa(int(bytes[0])) + "." +
		strconv.Itoa(int(bytes[1])) + "." +
		strconv.Itoa(int(bytes[2])) + "." +
		strconv.Itoa(int(bytes[3]))

	*ip = Ip(str)

	return nil
}
