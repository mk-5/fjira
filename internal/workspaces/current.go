package workspaces

func GetCurrent() (string, error) {
	s := NewUserHomeSettingsStorage()
	w, err := s.ReadCurrentWorkspace()
	if err != nil {
		return "", err
	}
	if w == "" {
		return "default", nil
	}
	return w, nil
}
