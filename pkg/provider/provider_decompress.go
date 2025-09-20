package provider

import (
	"fmt"
	"strings"
)

// Decompress gets a provider to decompress zip or tar.gz files
func Decompress(path string) (Provider, error) {
	if strings.HasSuffix(path, ".zip") {
		return &Zip{
			Path: path,
		}, nil
	} else if strings.HasSuffix(path, ".tar.gz") {
		return &Gzip{
			Path: path,
		}, nil
	}
	return nil, fmt.Errorf("provider.Decompress unknown file type for file: %s", path)
}
