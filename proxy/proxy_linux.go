package proxy

import (
	"fmt"
	"os/exec"
	"strings"
)

var exclusionListURLs = []string{
	"https://raw.githubusercontent.com/anfragment/zen/main/proxy/exclusions/common.txt",
}

func (p *Proxy) setSystemProxy() error {
	if binaryExists("gsettings") {
		commands := [][]string{
			[]string{"gsettings", "set", "org.gnome.system.proxy", "mode", "manual"},
			[]string{"gsettings", "set", "org.gnome.system.proxy.http", "host", "127.0.0.1"},
			[]string{"gsettings", "set", "org.gnome.system.proxy.http", "port", fmt.Sprint(p.port)},
			[]string{"gsettings", "set", "org.gnome.system.proxy.https", "host", "127.0.0.1"},
			[]string{"gsettings", "set", "org.gnome.system.proxy.https", "port", fmt.Sprint(p.port)},
		}

		for _, command := range commands {
			cmd := exec.Command(command[0], command[1:]...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("%s: %v\n%s", strings.Join(command, " "), err, out)
			}
		}
		return nil
	}
	// TODO: add support for other desktop environments

	return fmt.Errorf("system proxy configuration is currently only supported on GNOME")
}

func (p *Proxy) unsetSystemProxy() error {
	command := []string{"gsettings", "set", "org.gnome.system.proxy", "mode", "none"}
	cmd := exec.Command(command[0], command[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %v\n%s", strings.Join(command, " "), err, out)
	}

	return nil
}

func binaryExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
