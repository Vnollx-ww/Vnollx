package rpc

import "Vnollx/pkg/viper"

func init() {
	userConfig := viper.Init("user")
	InitUser(&userConfig)
	fileConfig := viper.Init("file")
	InitFile(&fileConfig)
	folderConfig := viper.Init("folder")
	InitFolder(&folderConfig)
}
