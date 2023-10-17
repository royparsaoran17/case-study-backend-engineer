// Package bootstrap
package bootstrap

import (
	"manage-se/internal/consts"
	"manage-se/pkg/logger"
	"manage-se/pkg/msgx"
)

func RegistryMessage() {
	err := msgx.Setup("msg.yaml", consts.ConfigPath)
	if err != nil {
		logger.Fatal(logger.MessageFormat("file message multi language load error %s", err.Error()))
	}

}
