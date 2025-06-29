package net

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strconv"
)

// ip为主机序
// InetNtoA func
func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// InetAtoN func
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

// PrefixToNetmask 如 24 对应的子网掩码地址为 255.255.255.0
func PrefixToNetmask(subnet int) string {
	var buff bytes.Buffer
	for i := 0; i < subnet; i++ {
		buff.WriteString("1")
	}
	for i := subnet; i < 32; i++ {
		buff.WriteString("0")
	}
	masker := buff.String()
	a, _ := strconv.ParseUint(masker[:8], 2, 64)
	b, _ := strconv.ParseUint(masker[8:16], 2, 64)
	c, _ := strconv.ParseUint(masker[16:24], 2, 64)
	d, _ := strconv.ParseUint(masker[24:32], 2, 64)
	resultMask := fmt.Sprintf("%v.%v.%v.%v", a, b, c, d)
	return resultMask
}

func NetmaskToPrefix(mask string) (int, error) {
	nip := net.ParseIP(mask)
	if nip == nil {
		return 0, errors.New("ip not valid")
	}

	ipv4 := nip.To4()
	ipmask := net.IPMask{ipv4[0], ipv4[1], ipv4[2], ipv4[3]}
	len, bit := ipmask.Size()
	if len == 0 || bit == 0 {
		return 0, errors.New("掩码格式错误")
	}
	return len, nil
}
