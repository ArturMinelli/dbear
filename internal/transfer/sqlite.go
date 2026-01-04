package transfer

import (
	"bytes"
	"fmt"
	"os/exec"

	"dbear/internal/config"
)

func DumpSQLite(conn config.Connection) ([]byte, error) {
	cmd := exec.Command("sqlite3", conn.Database, ".dump")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("sqlite dump failed: %w, stderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

func RestoreSQLite(conn config.Connection, dumpData []byte) error {
	cmd := exec.Command("sqlite3", conn.Database)
	cmd.Stdin = bytes.NewReader(dumpData)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sqlite restore failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

