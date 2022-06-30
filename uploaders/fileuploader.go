package uploaders

import (
	"encoding/json"
	"fmt"
)

// DeployFile ...
func DeployFile(pth, buildURL, token string) (ArtifactURLs, error) {
	return DeployFileWithMeta(pth, buildURL, token, "file", nil)
}

// DeployFileWithMeta ...
func DeployFileWithMeta(pth, buildURL, token, artifactType string, meta interface{}) (ArtifactURLs, error) {
	var metaBytes []byte
	if meta != nil {
		var err error
		metaBytes, err = json.Marshal(meta)
		if err != nil {
			return ArtifactURLs{}, fmt.Errorf("failed to marshal meta: %s", err)
		}
	}
	fmt.Printf("meta: %s\n", string(metaBytes))

	uploadURL, artifactID, err := createArtifact(buildURL, token, pth, artifactType, "")
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to create file artifact, error: %s", err)
	}

	if err := uploadArtifact(uploadURL, pth, ""); err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to upload file artifact, error: %s", err)
	}

	artifactURLs, err := finishArtifact(buildURL, token, artifactID, string(metaBytes), "", "", "no")
	if err != nil {
		return ArtifactURLs{}, fmt.Errorf("failed to finish file artifact, error: %s", err)
	}
	return artifactURLs, nil
}
