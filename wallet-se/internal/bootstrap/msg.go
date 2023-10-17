// Package bootstrap
package bootstrap

import (
	"wallet-se/internal/consts"
	"wallet-se/pkg/logger"
	"wallet-se/pkg/msgx"
)

func RegistryMessage() {
	err := msgx.Setup("msg.yaml", consts.ConfigPath)
	if err != nil {
		logger.Fatal(logger.MessageFormat("file message multi language load error %s", err.Error()))
	}

}
