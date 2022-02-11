package main

import (
	"i2pgit.org/idk/i2p.plugins.tor-manager/get"
	"os"
	"path/filepath"
)

func GenerateAppArmor() error {
	return nil
}

func CreateShortcuts() error {
	if err := CreateShortcut("torbrowser"); err != nil {
		return err
	}
	if err := CreateShortcut("i2pbrowser"); err != nil {
		return err
	}
	if err := CreateShortcut("i2pconfig"); err != nil {
		return err
	}
	return nil
}

func CreateShortcut(linkname string) error {
	// check if there is a symlink in the $HOME/Desktop Directory
	// if not, create one
	// if there is, check if it points to the correct location
	// if not, delete it and create a new one
	// if there is, do nothing
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	absolutepath, err := filepath.Abs(filepath.Join(filepath.Dir(exe), exe))
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, "Desktop", linkname)
	if tbget.FileExists(linkname) {
		if originfile, err := os.Readlink(path); err != nil || originfile != absolutepath {
			if err := os.Remove(path); err != nil {
				return err
			}
			return os.Symlink(absolutepath, path)
		}
	}
	return os.Symlink(absolutepath, path)
	return nil
}
