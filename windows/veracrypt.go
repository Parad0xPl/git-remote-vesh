package windows

type VeraCryptWinHandle struct {
	mountLetter string
	valuePath   string
}

const defaultPath = "C:\\Program Files\\VeraCrypt\\VeraCrypt.exe"

func getExecPath() string {
	path, _ := exec.LookPath("VeraCrypt")

	if path == "" {
		path = defaultPath
	}

	return path
}

func createWinVeraCrypt() *VeraCryptWinHandle {

	return &VeraCryptWinHandle{
		mountLetter: "A",
		valuePath:   "Z:\\bit_revolver.crypt",
	}
}

func (s *VeraCryptWinHandle) Start() error {
	path := getExecPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find veracrypt executable: %s", err)
	}

	arguments :=
		[]string{
			"/q",
			"/l",
			s.mountLetter,
			"/a",
			"/v",
			s.valuePath,
		}

	s.cmd = exec.Command(path, arguments...)
	err = s.Run()
	if err != nil {
		return fmt.Errorf("can't start veracrypt executable: %s", err)
	}
	return nil
}

func (s *VeraCryptWinHandle) Stop() error {
	path := getExecPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find veracrypt executable: %s", err)
	}

	arguments :=
		[]string{
			"/q",
			"/d",
			s.mountLetter,
		}

	s.cmd = exec.Command(path, arguments...)
	err = s.Run()
	if err != nil {
		return fmt.Errorf("can't start veracrypt executable: %s", err)
	}
	return nil
}
