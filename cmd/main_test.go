package main

import (
	"fmt"
	"testing"

	"github.com/LordMathis/GitEcho/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {

	configPath := "../config.yaml"
	config, err := config.ReadConfig(configPath)

	fmt.Println(config)

	assert.NoError(t, err)

}

func cleanup() {
	return
}
