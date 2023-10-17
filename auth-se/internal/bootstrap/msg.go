// Package bootstrap
package bootstrap

import (
	"auth-se/internal/consts"
	"auth-se/pkg/logger"
	"auth-se/pkg/msgx"
)

func RegistryMessage() {
	err := msgx.Setup("msg.yaml", consts.ConfigPath)
	if err != nil {
		logger.Fatal(logger.MessageFormat("file message multi language load error %s", err.Error()))
	}

}
