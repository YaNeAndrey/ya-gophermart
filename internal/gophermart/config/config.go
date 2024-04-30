package config

import (
	"fmt"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/consterror"
	"strconv"
	"strings"
)

type Config struct {
	srvAddr     string
	srvPort     int
	accrualAddr string
}

func (c *Config) SetSrvAddr(srvAddrStr string) error {
	srvAddr, srvPort, err := parseEndpoint(srvAddrStr)
	if err != nil {
		return err
	}
	if srvPort > 65535 || srvPort < 0 {
		return consterror.ErrIncorrectPortNumber
	} else {
		c.srvPort = srvPort
	}
	c.srvAddr = srvAddr
	return nil
}

func (c *Config) SetAccrualAddr(accrualAddrStr string) error {
	c.accrualAddr = accrualAddrStr
	return nil
}

func (c *Config) GetSrvAddr() string {
	return fmt.Sprintf("%s:%d", c.srvAddr, c.srvPort)
}

func (c *Config) GetAccrualAddr() string {
	return c.accrualAddr
}

func parseEndpoint(endpointStr string) (string, int, error) {
	hp := strings.Split(endpointStr, ":")
	if len(hp) != 2 {
		return "", 0, consterror.ErrIncorrectEndpointFormat
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return "", 0, err
	}
	return hp[0], port, nil
}
