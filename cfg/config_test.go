package cfg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyConfiguration(t *testing.T) {
	config, err := loadConfig(bytes.NewReader([]byte("---")))
	assert.NoError(t, err)
	assert.Len(t, config.Servers, 0)
}

func TestSimpleConfiguration(t *testing.T) {
	content := `---
servers:
  # comment
  http:
    listen: http://:8080
  https:
    listen: http://:8443
    certificate: certificate.pem
    privateKey: key.pem
`
	config, err := loadConfig(bytes.NewReader([]byte(content)))
	assert.NoError(t, err)
	assert.Len(t, config.Servers, 2)
}
