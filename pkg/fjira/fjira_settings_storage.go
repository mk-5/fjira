package fjira

type userHomeSettingsStorage struct{}

type settingsStorage interface { //nolint
	write(settings *fjiraSettings) error
	read() (*fjiraSettings, error)
}

func (s *userHomeSettingsStorage) read() (*fjiraSettings, error) {
	return nil, nil
}

func (s *userHomeSettingsStorage) write() (*fjiraSettings, error) {
	return nil, nil
}
