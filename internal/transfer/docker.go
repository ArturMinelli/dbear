package transfer

import (
	"bytes"
	"fmt"
	"os/exec"

	"dbear/internal/config"
)

func DumpDatabase(conn config.Connection, dockerImage string, schemas []string) ([]byte, error) {
	var dumpCmd *exec.Cmd

	if conn.Type == config.TypePostgreSQL {
		dumpCmd = buildPostgreSQLDumpCommand(conn, dockerImage, schemas)
	} else if conn.Type == config.TypeMySQL {
		dumpCmd = buildMySQLDumpCommand(conn, dockerImage)
	} else {
		return nil, fmt.Errorf("unsupported database type for docker dump: %s", conn.Type)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	dumpCmd.Stdout = &stdout
	dumpCmd.Stderr = &stderr

	if err := dumpCmd.Run(); err != nil {
		return nil, fmt.Errorf("docker dump failed: %w, stderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

func RestoreDatabase(conn config.Connection, dockerImage string, dumpData []byte) error {
	var restoreCmd *exec.Cmd

	if conn.Type == config.TypePostgreSQL {
		restoreCmd = buildPostgreSQLRestoreCommand(conn, dockerImage, dumpData)
	} else if conn.Type == config.TypeMySQL {
		restoreCmd = buildMySQLRestoreCommand(conn, dockerImage, dumpData)
	} else {
		return fmt.Errorf("unsupported database type for docker restore: %s", conn.Type)
	}

	var stderr bytes.Buffer
	restoreCmd.Stderr = &stderr

	if err := restoreCmd.Run(); err != nil {
		return fmt.Errorf("docker restore failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

func buildPostgreSQLDumpCommand(conn config.Connection, dockerImage string, schemas []string) *exec.Cmd {
	env := []string{
		fmt.Sprintf("PGHOST=%s", conn.Host),
		fmt.Sprintf("PGPORT=%d", conn.Port),
		fmt.Sprintf("PGUSER=%s", conn.Username),
		fmt.Sprintf("PGPASSWORD=%s", conn.Password),
		fmt.Sprintf("PGDATABASE=%s", conn.Database),
	}

	args := []string{
		"run",
		"--rm",
		"--network", "host",
	}

	for _, envVar := range env {
		args = append(args, "-e", envVar)
	}

	args = append(args, dockerImage, "pg_dump", "--no-owner", "--no-acl", "-Fc")
	if len(schemas) > 0 {
		for _, schema := range schemas {
			args = append(args, "-n", schema)
		}
	}

	return exec.Command("docker", args...)
}

func buildPostgreSQLRestoreCommand(conn config.Connection, dockerImage string, dumpData []byte) *exec.Cmd {
	env := []string{
		fmt.Sprintf("PGHOST=%s", conn.Host),
		fmt.Sprintf("PGPORT=%d", conn.Port),
		fmt.Sprintf("PGUSER=%s", conn.Username),
		fmt.Sprintf("PGPASSWORD=%s", conn.Password),
		fmt.Sprintf("PGDATABASE=%s", conn.Database),
	}

	args := []string{
		"run",
		"--rm",
		"--network", "host",
		"-i",
	}

	for _, e := range env {
		args = append(args, "-e", e)
	}

	args = append(args, dockerImage, "pg_restore", "--no-owner", "--no-acl", "--disable-triggers", "-d", conn.Database)

	cmd := exec.Command("docker", args...)
	cmd.Stdin = bytes.NewReader(dumpData)

	return cmd
}

func buildMySQLDumpCommand(conn config.Connection, dockerImage string) *exec.Cmd {
	args := []string{
		"run",
		"--rm",
		"--network", "host",
		dockerImage,
		"mysqldump",
		"-h", conn.Host,
		"-P", fmt.Sprintf("%d", conn.Port),
		"-u", conn.Username,
		fmt.Sprintf("-p%s", conn.Password),
		conn.Database,
	}

	return exec.Command("docker", args...)
}

func buildMySQLRestoreCommand(conn config.Connection, dockerImage string, dumpData []byte) *exec.Cmd {
	args := []string{
		"run",
		"--rm",
		"--network", "host",
		"-i",
		dockerImage,
		"mysql",
		"-h", conn.Host,
		"-P", fmt.Sprintf("%d", conn.Port),
		"-u", conn.Username,
		fmt.Sprintf("-p%s", conn.Password),
		conn.Database,
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdin = bytes.NewReader(dumpData)

	return cmd
}

