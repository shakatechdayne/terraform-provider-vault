package generate

import (
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/logical/framework"
)

const secretEngineTag = "secrets"

func Run(logger hclog.Logger, doc *framework.OASDocument) error {
	for pathName, pathItem := range doc.Paths {
		if strings.Contains(pathName, "config") {
			logger.Info("Skipping " + pathName + " because it's a config path and may contain sensitive fields")
			continue
		}
		info, err := mkDirs(pathName, pathItem)
		if err != nil {
			return err
		}
		if err := mkResource(pathName, pathItem, info); err != nil {
			return err
		}
	}
	return nil
}

// TODO we're getting a top-level "secret" one, and under "secrets" we're getting both "secret" and "kv".
// Given a pathName like "/auth/ldap/duo/access", mkDirs will make any of
// the the four following nested folders that don't yet exist: "resources/auth/ldap/duo".
func mkDirs(pathName string, pathItem *framework.OASPathItem) (*resourceInfo, error) {
	// Trim any leading slash.
	if strings.HasPrefix(pathName, "/") {
		pathName = pathName[1:]
	}
	if isSecretEngine(pathItem) {
		// Prepend the path with secrets to sort secret engines into their
		// own directory.
		pathName = "secrets/" + pathName
	}
	// Place all generated code under a directory called "resources".
	pathName = "resources/" + pathName

	pathHeirarchy := strings.Split(pathName, "/")

	// Strip out any fields that are actually a path parameter, and anything below them.
	for i, pathField := range pathHeirarchy {
		if strings.Contains(pathField, "{") {
			pathHeirarchy = pathHeirarchy[:i]
			break
		}
	}

	info := &resourceInfo{
		fileName: pathHeirarchy[len(pathHeirarchy)-1],
	}
	for i := 0; i < len(pathHeirarchy)-1; i++ {
		joinedPath := strings.Join(pathHeirarchy[:i+1], "/")

		// The values these receive on the last iteration of the loop will be
		// correct.
		info.packageName = pathHeirarchy[i]
		info.pathToPackage = joinedPath

		if _, err := os.Stat(joinedPath); err != nil {
			if err := os.Mkdir(joinedPath, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}
	return info, nil
}

type resourceInfo struct {
	// ex. "resources/secrets/alicloud"
	pathToPackage string

	// ex. "alicloud"
	packageName string

	// ex. "role"
	fileName string
}

func mkResource(pathName string, pathItem *framework.OASPathItem, info *resourceInfo) error {
	// TODO
	return nil
}

func isSecretEngine(pathItem *framework.OASPathItem) bool {
	if pathItem.Get != nil {
		for _, tag := range pathItem.Get.Tags {
			if tag == secretEngineTag {
				return true
			}
		}
	}
	if pathItem.Post != nil {
		for _, tag := range pathItem.Post.Tags {
			if tag == secretEngineTag {
				return true
			}
		}
	}
	if pathItem.Delete != nil {
		for _, tag := range pathItem.Delete.Tags {
			if tag == secretEngineTag {
				return true
			}
		}
	}
	return false
}
