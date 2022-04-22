package fjira

type textEditor string

const (
	Nano textEditor = "nano"
	Vim  textEditor = "vim"
)

type fjiraSettings struct {
	textEditor     string
	jiraRestUrl    string
	jiraBasicToken string
}

type fjiraWorkspace struct {
	name     string
	settings *fjiraSettings
}
