package onec

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func loadConfigurationInfo(binPath, connectionString string, login, pass string) (*ConfigurationInfo, error) {
	workDir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, errors.Wrap(err, "create tmp dir")
	}

	listFile := createListFile()
	defer os.Remove(listFile)
	defer os.RemoveAll(workDir)

	var params []string
	params = append(params, "DESIGNER")
	params = append(params, "/DumpConfigToFiles", workDir)
	params = append(params, "-listFile", listFile)
	params = append(params, "/DisableStartupDialogs")
	params = append(params, "/DisableStartupMessages")

	if login != "" {
		params = append(params, "/N", login)
	}
	if pass != "" {
		params = append(params, "/P", pass)
	}

	if strings.HasPrefix(connectionString, "File") {
		re := regexp.MustCompile(`"([^"]*)"`)
		matches := re.FindStringSubmatch(connectionString)
		if len(matches) > 1 {
			params = append(params, "/F", matches[1])
		}
	}
	if strings.HasPrefix(connectionString, "Srvr=") {
		re := regexp.MustCompile(`(?m)Srvr="([^"]*)"[^"]*"([^"]*)`)
		matches := re.FindStringSubmatch(connectionString)
		if len(matches) == 3 {
			params = append(params, "/S", fmt.Sprintf("%s/%s", matches[1], matches[2]))
		}
	}

	err = run(context.Background(), binPath, params)
	if err != nil {
		return nil, errors.Wrap(err, "run error")
	}

	return readConfigurationFile(workDir)
}

func run(ctx context.Context, binPath string, params []string) error {
	runCtx, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	logFile, _ := os.CreateTemp("", "")
	logFilePath := logFile.Name()
	logFile.Close()

	defer os.Remove(logFilePath)

	params = append(params, "/out", logFilePath)

	cmd := exec.CommandContext(runCtx, binPath, params...)
	cmd.Stdout = new(bytes.Buffer)
	cmd.Stderr = new(bytes.Buffer)

	err := cmd.Run()
	if err != nil || cmd.ProcessState.ExitCode() != 0 {
		stderr := cmd.Stderr.(*bytes.Buffer).String()
		if stderr != "" {
			return errors.New(stderr)
		}
		if log := readLogFile(logFilePath); log != "" {
			return errors.New(log)
		}
	}

	return nil
}

func readLogFile(path string) string {
	if data, err := os.ReadFile(path); err == nil {
		reader := transform.NewReader(bytes.NewReader(data), unicode.BOMOverride(unicode.UTF8BOM.NewDecoder()))

		data, _ = io.ReadAll(reader)
		return strings.TrimSpace(string(data))
	}

	return ""
}

func createListFile() string {
	if f, err := os.CreateTemp("", ""); err == nil {
		f.WriteString("Configuration")
		f.Close()
		return f.Name()
	}

	return ""
}

func loadExtensionsInfo(binPath, connectionString, login, pass string) (string, []ConfigurationInfo, error) {
	workDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", nil, errors.Wrap(err, "create tmp dir")
	}

	var params []string
	params = append(params, "DESIGNER")
	params = append(params, "/DumpConfigToFiles", workDir)
	params = append(params, "-AllExtensions")
	params = append(params, "/DisableStartupDialogs")
	params = append(params, "/DisableStartupMessages")

	if login != "" {
		params = append(params, "/N", login)
	}
	if pass != "" {
		params = append(params, "/P", pass)
	}

	if strings.HasPrefix(connectionString, "File") {
		re := regexp.MustCompile(`"([^"]*)"`)
		matches := re.FindStringSubmatch(connectionString)
		if len(matches) > 1 {
			params = append(params, "/F", matches[1])
		}
	}
	if strings.HasPrefix(connectionString, "Srvr=") {
		re := regexp.MustCompile(`(?m)Srvr="([^"]*)"[^"]*"([^"]*)`)
		matches := re.FindStringSubmatch(connectionString)
		if len(matches) == 3 {
			params = append(params, "/S", fmt.Sprintf("%s/%s", matches[1], matches[2]))
		}
	}

	err = run(context.Background(), binPath, params)
	if err != nil {
		return "", nil, errors.Wrap(err, "run error")
	}

	dirs, err := os.ReadDir(workDir)
	if err != nil {
		return "", nil, err
	}

	var result []ConfigurationInfo

	for _, dir := range dirs {
		if info, err := readConfigurationFile(filepath.Join(workDir, dir.Name())); err == nil {
			result = append(result, *info)
		} else {
			log.Println(errors.Wrap(err, "load extensions info"))
		}
	}

	return workDir, result, nil
}
