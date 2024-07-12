//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"strings"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
)

var (
	app *portapps.App
)

const (
	vmOptionsFile = "webstorm64.vmoptions"
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("webstorm", "WebStorm"); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	webideExe := "webstorm64.exe"
	webideVmOptionsFile := "webstorm64.exe.vmoptions"

	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "bin", webideExe)
	app.WorkingDir = utl.PathJoin(app.AppPath, "bin")

	// override idea.properties
	webidePropContent := strings.Replace(`# DO NOT EDIT! AUTOMATICALLY GENERATED BY PORTAPPS.
	webide.config.path={{ DATA_PATH }}/config
	webide.system.path={{ DATA_PATH }}/system
	webide.plugins.path={{ DATA_PATH }}/plugins
	webide.log.path={{ DATA_PATH }}/log`, "{{ DATA_PATH }}", utl.FormatUnixPath(app.DataPath), -1)

	webidePropPath := utl.PathJoin(app.DataPath, "idea.properties")
	if err := utl.CreateFile(webidePropPath, webidePropContent); err != nil {
		log.Fatal().Err(err).Msg("Cannot write idea.properties")
	}

	// https://www.jetbrains.com/help/webstorm/tuning-the-ide.html#configure-platform-properties
	os.Setenv("WEBIDE_PROPERTIES", webidePropPath)

	// https://www.jetbrains.com/help/webstorm/tuning-the-ide.html#configure-jvm-options
	os.Setenv("WEBIDE_VM_OPTIONS", utl.PathJoin(app.DataPath, vmOptionsFile))
	if !utl.Exists(utl.PathJoin(app.DataPath, vmOptionsFile)) {
		utl.CopyFile(utl.PathJoin(app.AppPath, "bin", webideVmOptionsFile), utl.PathJoin(app.DataPath, vmOptionsFile))
	} else {
		utl.CopyFile(utl.PathJoin(app.DataPath, vmOptionsFile), utl.PathJoin(app.AppPath, "bin", webideVmOptionsFile))
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
