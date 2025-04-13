package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/day0ops/randomise-route-keys/pkg/config"
)

// cache for route list
var routerList RouteStrList

func ReadRouteListFile() error {
	path := config.GetEnv(config.RouteListFilePathEnvVar, config.DefaultRouteListFilePath)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("no route list file found with path %s", path)
	}

	rList, err := readContentAsJson(path)
	if err != nil {
		return fmt.Errorf("error trying to read the content: %w", err)
	}
	routerList = rList
	return nil
}

func GetCachedRouteList() RouteStrList {
	return routerList
}

func ClearCachedRouteList() {
	routerList = RouteStrList{}
}

func readContentAsJson(filename string) (RouteStrList, error) {
	var data RouteStrList
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
