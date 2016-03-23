package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
)

const (
	lvmPluginSocketPath  = "/run/docker/plugins/lvm.sock"
	vgConfigPath         = "/etc/docker/docker-lvm-plugin"
	lvmHome              = "/var/lib/docker-lvm-plugin"
	lvmVolumesConfigPath = "/var/lib/docker-lvm-plugin/lvmVolumesConfig.json"
	lvmCountConfigPath   = "/var/lib/docker-lvm-plugin/lvmCountConfig.json"
)

var (
	flVersion *bool
	flDebug   *bool
)

func init() {
	flVersion = flag.Bool("version", false, "Print version information and quit")
	flDebug = flag.Bool("debug", false, "Enable debug logging")
}

func main() {

	flag.Parse()

	if *flVersion {
		fmt.Fprint(os.Stdout, "docker lvm plugin version: 1.0\n")
		return
	}

	if *flDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if _, err := os.Stat(lvmHome); err != nil {
		if !os.IsNotExist(err) {
			logrus.Fatal(err)
		}
		logrus.Debugf("Created home dir at %s", lvmHome)
		if err := os.MkdirAll(lvmHome, 0700); err != nil {
			logrus.Fatal(err)
		}
	}

	lvm := newDriver(lvmHome, vgConfigPath)

	// Call loadFromDisk only if config file exists.
	if _, err := os.Stat(lvmVolumesConfigPath); err == nil {
		if err := loadFromDisk(lvm); err != nil {
			logrus.Fatal(err)
		}
	}

	h := volume.NewHandler(lvm)
	if err := h.ServeUnix("root", lvmPluginSocketPath); err != nil {
		logrus.Fatal(err)
	}
}
