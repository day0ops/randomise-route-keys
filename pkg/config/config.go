package config

import "os"

const DefaultRouteListFilePath = "/etc/config/route-list.json"
const RouteListFilePathEnvVar = "ROUTE_LIST_FILE_PATH"

var LogLevel = os.Getenv("LOG_LEVEL")

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
