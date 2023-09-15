package issues

import (
	"bufio"
	"fmt"
	"github.com/mk-5/fjira/internal/workspaces"
	"io"
	"os"
	"strings"
)

type jqlStorage struct {
}

func (s *jqlStorage) addNew(jql string) error {
	lines, err := s.readAll()
	if err != nil {
		return err
	}
	normalizedJql := s.normalizeJql(jql)

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, normalizedJql)
	if len(lines) > 0 {
		newLines = append(newLines, lines...)
	}

	err = s.writeAll(newLines)

	return err
}

func (s *jqlStorage) remove(jql string) error {
	lines, _ := s.readAll()
	newLines := make([]string, 0, len(lines))

	for i := range lines {
		if lines[i] == jql {
			continue
		}
		newLines = append(newLines, lines[i])
	}

	err := s.writeAll(newLines)

	return err
}

func (s *jqlStorage) readAll() ([]string, error) {
	jqlFile, err := s.jqlsFile()
	defer func(jqlFile *os.File) {
		_ = jqlFile.Close()
	}(jqlFile)
	if err != nil {
		return nil, err
	}
	lines := make([]string, 0, MaxJqlLines)
	// it's not that many jqls stored - we can load them into memory
	scanner := bufio.NewScanner(jqlFile)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, string(line))
	}
	return lines, nil
}

func (s *jqlStorage) writeAll(jqls []string) error {
	jqlFile, err := s.jqlsFile()
	defer func(jqlFile *os.File) {
		_ = jqlFile.Close()
	}(jqlFile)
	if err != nil {
		return err
	}
	if len(jqls) == 0 {
		err := os.Truncate(jqlFile.Name(), 0)
		return err
	}
	if len(jqls) > MaxJqlLines {
		jqls = jqls[:MaxJqlLines]
	}
	_, err = jqlFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	for _, line := range jqls {
		if strings.TrimSpace(line) == "" {
			continue
		}
		_, err = jqlFile.WriteString(fmt.Sprintln(line))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *jqlStorage) jqlsFile() (*os.File, error) {
	userHomeStorage := workspaces.NewUserHomeSettingsStorage()
	configDir, err := userHomeStorage.ConfigDir()
	if err != nil {
		return nil, err
	}
	currentWorkspace, err := workspaces.GetCurrent()
	if err != nil {
		return nil, err
	}
	jqlFilePath := fmt.Sprintf("%s/%s.jqls", configDir, currentWorkspace)
	jqlFile, err := os.OpenFile(jqlFilePath, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return nil, err
	}
	return jqlFile, nil
}

func (s *jqlStorage) normalizeJql(jql string) string {
	jql = strings.TrimSpace(jql)
	jql = strings.ReplaceAll(jql, "\n", " ")
	jql = strings.ReplaceAll(jql, "\t", " ")
	return jql
}
