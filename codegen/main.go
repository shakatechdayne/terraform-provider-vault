package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/terraform-providers/terraform-provider-vault/codegen/generate"
	"github.com/terraform-providers/terraform-provider-vault/codegen/github"
)

// TODO run via make
// TODO provide way to switch between this or deprecated code
// Main assumes it's being run from this project's home directory.
func main() {
	logger := hclog.Default()

	logger.Info("Checking for latest Vault release")
	githubClient := github.NewClient()
	lastReleaseTag, err := githubClient.LatestTag("hashicorp", "vault")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info(fmt.Sprintf("Lastest release is %s: %s", lastReleaseTag.Name, lastReleaseTag.Commit.SHA))

	logger.Info("Verifying the latest release locally")
	if output, err := exec.Command("vault", "-version").Output(); err != nil {
		logger.Error(err.Error())
		return
	} else {
		if !strings.Contains(string(output), lastReleaseTag.Commit.SHA) {
			msg := fmt.Sprintf("cannot generate code because Vault release %s at commit %s needs to be checked out and built locally", lastReleaseTag.Name, lastReleaseTag.Commit.SHA)
			logger.Error(msg)
			return
		}
	}
	logger.Info("Latest release found locally")

	cmd := exec.Command("sh", "codegen/scripts/gen_openapi.sh")
	wd, err := os.Getwd()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	cmd.Dir = wd
	cmd.Stdout = os.Stdout

	logger.Info("Running Vault locally to generate openapi.json")
	if err := cmd.Run(); err != nil {
		logger.Error(err.Error())
		return
	}
	defer func() {
		logger.Info("Stopping all Vault processes")
		if err := exec.Command("killall", "-9", "vault").Run(); err != nil {
			logger.Error(err.Error())
		}
		logger.Info("Cleaning up openapi.json")
		if err := os.Remove("openapi.json"); err != nil {
			logger.Error(err.Error())
		}
	}()

	doc, err := loadOAS("openapi.json")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("Finished generating openapi.json")

	logger.Info("Generating code")
	if err := generate.Run(logger, doc); err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("Successfully generated code!")
}

func loadOAS(filename string) (*framework.OASDocument, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var d map[string]interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}

	oas, err := framework.NewOASDocumentFromMap(d)
	if err != nil {
		return nil, err
	}

	return oas, nil
}
