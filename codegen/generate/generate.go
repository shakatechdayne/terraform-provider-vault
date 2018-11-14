package generate

import (
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/logical/framework"
)

var secretEngines = []string{
		"ad",
		"alicloud",
		"aws",
		"azure",
		"cassandra",
		"consul",
		"cubbyhole",
		"database",
		"gcp",
		"kv",
		"mongodb",
		"mssql",
		"mysql",
		"nomad",
		"pki",
		"postgresql",
		"rabbitmq",
		"ssh",
		"secret",
		"totp",
		"transit",
	}

func Run(logger hclog.Logger, doc *framework.OASDocument) error {
	for pathName := range doc.Paths {
		if strings.Contains(pathName, "config") {
			logger.Debug("Skipping " + pathName + " because it's a config path and may contain sensitive fields")
			continue
		}

		if err := mkDirs(pathName); err != nil {
			return err
		}
	}
	return nil
}

func mkDirs(pathName string) error {
	if strings.HasPrefix(pathName, "/") {
		pathName = pathName[1:]
	}
	if !strings.HasPrefix(pathName, "auth") {
		for _, secretEngine := range secretEngines {
			if strings.HasPrefix(pathName, secretEngine) {
				pathName = "secrets/" + pathName
				break
			}
		}
	}
	pathHeirarchy := strings.Split(pathName, "/")
	for i := 0; i < len(pathHeirarchy)-1; i++ {
		joinedPath := strings.Join(pathHeirarchy[:i+1], "/")
		if _, err := os.Stat(joinedPath); err != nil {
			if err := os.Mkdir(joinedPath, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}
