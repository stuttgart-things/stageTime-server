/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintBanner(t *testing.T) {
	assert := assert.New(t)
	expectedVersion := "Commit: 1a3a6e2\nDate: 24.0110.1136\nVersion: 0.2.0-ADD_VERSION_CMD"

	date = "24.0110.1136"
	version = "0.2.0-ADD_VERSION_CMD"
	commit = "1a3a6e2"

	versionInformation := PrintBanner()
	fmt.Println(versionInformation)
	fmt.Println(expectedVersion)

	assert.Equal(strings.TrimSpace(expectedVersion), strings.TrimSpace(versionInformation))
}
