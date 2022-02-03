package main

import (
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"os/user"
	"path/filepath"
)

func GenerateAppArmor() error {
	return nil
}

func DesktopDirectory() (string, error) {
	myself, error := user.Current()
	if error != nil {
		return "", error
	}
	homedir := myself.HomeDir
	desktop := filepath.Join(homedir, "Desktop")
	return desktop, nil
}

func CreateShortcuts() error {
	desktopDir, err := DesktopDirectory()
	if err != nil {
		return err
	}
	desktop := filepath.Join(desktopDir, "torbrowser.lnk")
	torBrowserPath, err := pathToMe()
	if err != nil {
		return err
	}
	torbrowserShortcut := torBrowserPath + " -torbrowser"
	if err := makeLink(torbrowserShortcut, desktop); err != nil {
		return err
	}
	i2pBrowserPath := torBrowserPath + " -i2pbrowser"
	if err := makeLink(i2pBrowserPath, desktop); err != nil {
		return err
	}
	return nil
}

func makeLink(src, dst string) error {
	log.Println("Creating desktop shortcut:", dst)
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", dst)
	if err != nil {
		return err
	}
	idispatch := cs.ToIDispatch()
	oleutil.PutProperty(idispatch, "TargetPath", src)
	oleutil.CallMethod(idispatch, "Save")
	return nil
}
