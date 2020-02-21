package supply

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Stager interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/stager.go
	BuildDir() string
	DepDir() string
	DepsIdx() string
	DepsDir() string
	WriteProfileD(string, string) error
	LinkDirectoryInDepDir(string, string) error
	WriteEnvFile(string, string) error
}

type Manifest interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/manifest.go
	AllDependencyVersions(string) []string
	DefaultVersion(string) (libbuildpack.Dependency, error)
}

type Installer interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/installer.go
	InstallDependency(libbuildpack.Dependency, string) error
	InstallOnlyVersion(string, string) error
}

type Command interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/command.go
	Execute(string, io.Writer, io.Writer, string, ...string) error
	Output(dir string, program string, args ...string) (string, error)
}

type Supplier struct {
	Manifest  Manifest
	Installer Installer
	Stager    Stager
	Command   Command
	Log       *libbuildpack.Logger
}

func (s *Supplier) Run() error {
	s.Log.BeginStep("Supplying osgeo")

	var dep libbuildpack.Dependency
	dep, err := s.Manifest.DefaultVersion("osgeo")
	if err != nil {
		return err
	}
	OsgeoInstallDir := filepath.Join(s.Stager.DepDir(), "osgeo")
	if err := s.Installer.InstallDependency(dep, OsgeoInstallDir); err != nil {
		return err
	}

	if err := s.Stager.LinkDirectoryInDepDir(filepath.Join(s.Stager.DepDir(), "osgeo", "bin"), "bin"); err != nil {
		return err
	}
	if err := s.Stager.LinkDirectoryInDepDir(filepath.Join(s.Stager.DepDir(), "osgeo", "lib"), "lib"); err != nil {
		return err
	}
	if err := s.Stager.LinkDirectoryInDepDir(filepath.Join(s.Stager.DepDir(), "osgeo", "include"), "include"); err != nil {
		return err
	}

	var environmentVars = map[string]string{
		"GDAL_DATA":          filepath.Join(OsgeoInstallDir, "share/gdal"),
		"PROJ_LIB":           filepath.Join(OsgeoInstallDir, "share/proj"),
		"LDFLAGS":            "-L" + filepath.Join(OsgeoInstallDir, "lib"),
		"CPLUS_INCLUDE_PATH": filepath.Join(OsgeoInstallDir, "include"),
		"C_INCLUDE_PATH":     filepath.Join(OsgeoInstallDir, "include"),
	}

	for envVar, envValue := range environmentVars {
		fmt.Print(envValue)
		if err := s.Stager.WriteEnvFile(envVar, envValue); err != nil {
			return err
		}
	}

	scriptContents := fmt.Sprintf(`export GDAL_DATA=$DEPS_DIR/%s/osgeo/share/gdal
export PROJ_LIB=$DEPS_DIR/%s/osgeo/share/proj
export LDFLAGS=-L$DEPS_DIR/%s/osgeo/lib
export CPLUS_INCLUDE_PATH=$DEPS_DIR/%s/osgeo/include/
export C_INCLUDE_PATH=$DEPS_DIR/%s/osgeo/include/
`, s.Stager.DepsIdx(), s.Stager.DepsIdx(), s.Stager.DepsIdx(),
		s.Stager.DepsIdx(), s.Stager.DepsIdx())

	return s.Stager.WriteProfileD("osgeo.sh", scriptContents)
	return nil
}
