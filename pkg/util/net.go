package util

import (
	"fmt"
	"net"

	"github.com/go-kit/kit/log/level"

	util_log "github.com/cortexproject/cortex/pkg/util/log"
)

// GetFirstAddressOf returns the first IPv4 address of the supplied interface names.
func GetFirstAddressOf(names []string) (string, error) {
	for _, name := range names {
		inf, err := net.InterfaceByName(name)
		if err != nil {
			level.Warn(util_log.Logger).Log("msg", "error getting interface", "inf", name, "err", err)
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			level.Warn(util_log.Logger).Log("msg", "error getting addresses for interface", "inf", name, "err", err)
			continue
		}
		if len(addrs) <= 0 {
			level.Warn(util_log.Logger).Log("msg", "no addresses found for interface", "inf", name, "err", err)
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if ip := v.IP.To4(); ip != nil {
					return v.IP.String(), nil
				}
				if ip := v.IP.To16(); ip != nil {
					return fmt.Sprintf("[%s]", v.IP.String()), nil
				}
			}
		}
	}

	return "", fmt.Errorf("No address found for %s", names)
}
