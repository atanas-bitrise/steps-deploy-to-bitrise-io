package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/ziputil"
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
	IsDir  bool   `json:"is_dir"`
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

		fileInfo, err := os.Stat(pth)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			tmpDir, err := pathutil.NormalizedOSTempDirPath("pipeline_file_share")
			if err != nil {
				return err
			}

			name := strings.TrimSuffix(filepath.Base(pth), filepath.Ext(pth))
			targetPth := filepath.Join(tmpDir, name+".zip")

			if err := ziputil.ZipDir(pth, targetPth, true); err != nil {
				return fmt.Errorf("failed to zip output dir, error: %s", err)
			}

			pth = targetPth
			meta.IsDir = true
		}

		b, err := json.Marshal(meta)
		if err != nil {
			return err
		}

		fakeMeta := map[string]interface{}{
			"scheme": string(b),
		}

		_, err = uploaders.DeployFileWithMeta(pth, buildURL, buildAPIToken, "ios-xcarchive", fakeMeta)
		if err != nil {
			return fmt.Errorf("failed to push pipeline intermediate file (%s): %s", pth, err)
		}
	}

	return nil
}
