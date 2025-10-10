package convert

import (
	"fmt"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func DecodePortName(n string) (string, corev1.Protocol, int32) {
	protocol := PortProtocol(n)
	// <portName>:80
	if i := strings.Index(n, ":"); i > 0 {
		hostPort, _ := strconv.ParseInt(n[i+1:], 10, 32)
		return n[0:i], protocol, int32(hostPort)
	}
	return n, protocol, 0
}

func PortProtocol(n string) corev1.Protocol {
	if strings.HasPrefix(n, "udp-") {
		return corev1.ProtocolUDP
	} else if strings.HasPrefix(n, "sctp-") {
		return corev1.ProtocolSCTP
	} else {
		return corev1.ProtocolTCP
	}
}

func FormatPortName(name string, protocol corev1.Protocol, hostOrNotPort int32) string {
	prefix := ""
	switch protocol {
	case corev1.ProtocolUDP:
		prefix = "udp-"
	case corev1.ProtocolSCTP:
		prefix = "stcp-"
	}

	if prefix != "" && !strings.HasPrefix(name, prefix) {
		name = prefix + name
	}

	if hostOrNotPort > 0 {
		name = fmt.Sprintf("%s:%d", name, hostOrNotPort)
	}

	return name
}
