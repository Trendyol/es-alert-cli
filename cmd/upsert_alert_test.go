package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetFlagVariables_ValidFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("cluster", "http://example.com:9200", "")
	cmd.Flags().String("filename", "example.yaml", "")

	cluster, filename, ok := getFlags(cmd)
	assert.Equal(t, "http://example.com:9200", cluster)
	assert.Equal(t, "example.yaml", filename)
	assert.True(t, ok)
}

func TestGetFlagVariables_MissingClusterFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("filename", "example.yaml", "")
	_ = cmd.Flags().Set("cluster", "") // Simulate missing cluster flag
	_, _, ok := getFlags(cmd)
	assert.False(t, ok)
}

func TestGetFlagVariables_MissingFilenameFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("cluster", "http://example.com:9200", "")
	_ = cmd.Flags().Set("filename", "") // Simulate missing filename flag
	_, _, ok := getFlags(cmd)
	assert.False(t, ok)
}
