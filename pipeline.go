package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-deploy-to-bitrise-io/uploaders"
)

func parsePipelineIntermediateFiles(s string) (map[string]string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	intermediateFiles := map[string]string{}

	list := strings.Split(s, "\n")
	for _, item := range list {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		idx := strings.LastIndex(item, ":")
		if idx == -1 {
			return nil, fmt.Errorf("invalid item (%s): doesn't contain ':' character", item)
		}

		path := item[:idx]
		if path == "" {
			return nil, fmt.Errorf("invalid item (%s): doesn't specify file path", item)
		}

		key := item[idx+1:]
		if key == "" {
			return nil, fmt.Errorf("invalid item (%s): doesn't specify key", item)
		}

		intermediateFiles[path] = key
	}

	return intermediateFiles, nil
}

type PipelineIntermediateFileMeta struct {
	EnvKey string `json:"env_key"`
}

func PushPipelineIntermediateFiles(fileList, buildURL, buildAPIToken string) error {
	intermediateFiles, err := parsePipelineIntermediateFiles(fileList)
	if err != nil {
		return err
	}

	for pth, key := range intermediateFiles {
		fmt.Println()
		log.Donef("Pushing pipeline intermediate file: %s", pth)

		var err error
		pth, err = filepath.Abs(pth)
		if err != nil {
			return fmt.Errorf("failed to push pipeline intermediate file (%s): %s", pth, err)
		}

		meta := PipelineIntermediateFileMeta{
			EnvKey: key,
		}

		fmt.Printf("meta: %v\n", meta)

		_, err = uploaders.DeployFileWithMeta(pth, buildURL, buildAPIToken, meta)
		if err != nil {
			return fmt.Errorf("failed to push pipeline intermediate file (%s): %s", pth, err)
		}
	}

	return nil
}
