package onec

import (
	"bytes"
	"context"
	"github.com/beevik/etree"
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
	"your-app/internal/utils"
)

func loadConfigurationInfo(binPath, connectionString string) (*ConfigurationInfo, error) {
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

	if strings.HasPrefix(connectionString, "File") {
		re := regexp.MustCompile(`"([^"]*)"`)
		matches := re.FindStringSubmatch(connectionString)
		if len(matches) > 1 {
			params = append(params, "/F", matches[1])
		}
	}

	err = run(context.Background(), binPath, params)
	if err != nil {
		return nil, errors.Wrap(err, "run error")
	}

	return readConfigurationFile(workDir)
}

func readConfigurationFile(dir string) (*ConfigurationInfo, error) {
	f, err := os.Open(filepath.Join(dir, "Configuration.xml"))
	if err != nil {
		return nil, errors.Wrap(err, "open Configuration.xml error")
	}
	defer f.Close()

	doc := etree.NewDocument()
	_, err = doc.ReadFrom(f)
	if err != nil {
		return nil, errors.Wrap(err, "read Configuration.xml error")
	}

	elemName := doc.FindElement("/MetaDataObject/Configuration/Properties/Name")
	elemSynonym := doc.FindElement("/MetaDataObject/Configuration/Properties/Synonym/v8:item/v8:content")
	elemVersion := doc.FindElement("/MetaDataObject/Configuration/Properties/Version")
	elemVendor := doc.FindElement("/MetaDataObject/Configuration/Properties/Vendor")
	elemPurpose := doc.FindElement("/MetaDataObject/Configuration/Properties/ConfigurationExtensionPurpose")
	return &ConfigurationInfo{
		Name:    utils.Ptr(utils.Opt[etree.Element](elemName)).Text(),
		Synonym: utils.Ptr(utils.Opt[etree.Element](elemSynonym)).Text(),
		Version: utils.Ptr(utils.Opt[etree.Element](elemVersion)).Text(),
		Vendor:  utils.Ptr(utils.Opt[etree.Element](elemVendor)).Text(),
		Purpose: utils.Ptr(utils.Opt[etree.Element](elemPurpose)).Text(),
	}, nil
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

func loadExtensionsInfo(binPath, connectionString string) (string, []ConfigurationInfo, error) {
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

	if strings.HasPrefix(connectionString, "File") {
		re := regexp.MustCompile(`"([^"]*)"`)
		matches := re.FindStringSubmatch(connectionString)
		if len(matches) > 1 {
			params = append(params, "/F", matches[1])
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
