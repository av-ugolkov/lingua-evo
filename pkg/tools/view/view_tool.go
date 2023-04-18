package view

import (
	"fmt"

	"lingua-evo/internal/config"
)

func GetPathFolder(path string) string {
	return fmt.Sprintf("%s%s", getRoot(), path)
}

func GetPathFile(path string) string {
	return fmt.Sprintf(".%s%s", getRoot(), path)
}

func getRoot() string {
	return config.GetConfig().Front.Root
}
