package page

import (
	_ "embed"
	"encoding/base64"
	"sync"
)

var (
	//go:embed assets/favicon-32x32.png
	favicon32 []byte

	faviconOnce    sync.Once
	faviconDataURI string
)

func cronWorkflowFaviconDataURI() string {
	faviconOnce.Do(func() {
		encoded := base64.StdEncoding.EncodeToString(favicon32)
		faviconDataURI = "data:image/png;base64," + encoded
	})
	return faviconDataURI
}
