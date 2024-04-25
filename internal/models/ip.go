package models

import (
	"strconv"
	"strings"
)

const IP_SIZE = 4

type IP string

func (ip IP) Size() uint64 {
	return IP_SIZE
}

func (ip IP) Marshal() [IP_SIZE]byte {
	octets := strings.Split(string(ip), ".")

	octet0, _ := strconv.Atoi(octets[0])
	octet1, _ := strconv.Atoi(octets[1])
	octet2, _ := strconv.Atoi(octets[2])
	octet3, _ := strconv.Atoi(octets[3])

	return [IP_SIZE]byte{byte(octet0), byte(octet1), byte(octet2), byte(octet3)}
}

func (ip *IP) Unmarshal(bytes []byte) error {
	str := strconv.Itoa(int(bytes[0])) + "." +
		strconv.Itoa(int(bytes[1])) + "." +
		strconv.Itoa(int(bytes[2])) + "." +
		strconv.Itoa(int(bytes[3]))

	*ip = IP(str)

	return nil
}
