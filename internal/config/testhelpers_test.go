package config

import "os"

func osMkdirAll(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func osWriteFile(path, contents string) error {
	return os.WriteFile(path, []byte(contents), 0o644)
}
