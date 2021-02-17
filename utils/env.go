package utils

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//func GetIntEnv(env string, defaultValue int) int {
//	value := strings.Trim(os.Getenv(env), " ")
//	if len(value) == 0 {
//		return defaultValue
//	}
//	v, err := strconv.Atoi(value)
//	if err != nil {
//		return defaultValue
//	}
//
//	return v
//}

// 获取string类型环境变量
func GetStringEnv(env string, defaultValue string) string {
	value := strings.Trim(os.Getenv(env), " ")
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// 获取布尔类型环境变量
//func GetBoolEnv(env string, defaultValue bool) bool {
//	value := strings.Trim(os.Getenv(env), " ")
//	if len(value) == 0 {
//		return defaultValue
//	}
//	value = strings.ToLower(value)
//	switch value {
//	case "1", "true":
//		return true
//	case "0", "false":
//		return false
//	}
//	return defaultValue
//}

func GetSandboxGroupEnv(env string, sandboxId int, defaultValue []string) []string {
	//if err := config.Init(); err != nil {
	//	log.Error().Err(err).Msg("Init Config error")
	//}
	groupList := os.Getenv(env)
	groupConfig := SandboxGroupConfig{}
	err := json.Unmarshal([]byte(groupList), &groupConfig.GroupConfig)
	if err != nil {
		log.Error().Err(err).Msg("Parse json error")
		return nil
	}
	sandboxConfig := &SandboxConfig{}
	length := reflect.ValueOf(sandboxConfig)
	for _, value := range groupConfig.GroupConfig {
		if value.SandboxID != sandboxId {
			continue
		} else {
			result := make([]string, length.Elem().NumField())
			result[0] = strconv.FormatInt(int64(sandboxId), 10)
			result[1] = value.SandboxType
			result[2] = strconv.FormatInt(int64(value.Protect), 10)
			return result
		}
	}
	return defaultValue
}
