package overlay

import (
	"fmt"
	"github.com/hpcng/warewulf/internal/pkg/config"
	"github.com/hpcng/warewulf/internal/pkg/errors"
	"github.com/hpcng/warewulf/internal/pkg/node"
	"github.com/hpcng/warewulf/internal/pkg/util"
	"github.com/hpcng/warewulf/internal/pkg/vnfs"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func FindAllRuntimeOverlays() ([]string, error) {
	config := config.New()
	var ret []string

	wwlog.Printf(wwlog.DEBUG, "Looking for runtime overlays...")
	files, err := ioutil.ReadDir(config.RuntimeOverlayDir())
	if err != nil {
		return ret, err
	}

	for _, file := range files {
		wwlog.Printf(wwlog.DEBUG, "Evaluating runtime overlay: %s\n", file.Name())
		if file.IsDir() == true {
			ret = append(ret, file.Name())
		}
	}

	return ret, nil
}


func RuntimeOverlayInit(name string) error {
	config := config.New()

	path := config.RuntimeOverlaySource(name)

	if util.IsDir(path) == true {
		return errors.New("Runtime overlay already exists: "+name)
	}

	err := os.MkdirAll(path, 0755)

	return err
}


func RuntimeBuild(nodeList []node.NodeInfo, force bool) error {
	config := config.New()

	wwlog.SetIndent(4)

	for _, node := range nodeList {
		if node.RuntimeOverlay != "" {
			OverlayDir := config.RuntimeOverlaySource(node.RuntimeOverlay)
			OverlayFile := config.RuntimeOverlayImage(node.Fqdn)
			vnfs := vnfs.New(node.Vnfs)

			vnfsDir := config.VnfsChroot(vnfs.NameClean())

			wwlog.Printf(wwlog.VERBOSE, "Building Runtime Overlay for: %s\n", node.Fqdn)

			tmpDir, err := ioutil.TempDir(os.TempDir(), ".wwctl-runtime-overlay-")
			if err != nil {
				return err
			}

			if util.IsDir(OverlayDir) == false {
				wwlog.Printf(wwlog.WARN, "%-35s: Skipped (runtime overlay template not found)\n", node.Fqdn)
				continue
			}

			if util.IsDir(vnfsDir) == false {
				wwlog.Printf(wwlog.WARN, "%-35s: Skipped (VNFS not imported)\n", node.Fqdn)
				continue
			}

			err = os.MkdirAll(path.Dir(OverlayFile), 0755)
			if err != nil {
				return err
			}

			if force == false {
				wwlog.Printf(wwlog.DEBUG, "Checking if overlay is required\n")
			}
			if util.PathIsNewer(OverlayDir, OverlayFile) {
				if force == false {
					wwlog.Printf(wwlog.INFO, "%-35s: Skipping, overlay is current\n", node.Fqdn)
					continue
				}
			}

			wwlog.Printf(wwlog.DEBUG, "Changing directory to OverlayDir: %s\n", OverlayDir)
			err = os.Chdir(OverlayDir)
			if err != nil {
				wwlog.Printf(wwlog.ERROR, "Could not chdir() to OverlayDir: %s\n", OverlayDir)
				continue
			}

			wwlog.Printf(wwlog.DEBUG, "Walking the file system: %s\n", OverlayDir)
			err = filepath.Walk(".", func(location string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					wwlog.Printf(wwlog.DEBUG, "Found directory: %s\n", location)

					err := os.MkdirAll(path.Join(tmpDir, location), info.Mode())
					if err != nil {
						wwlog.Printf(wwlog.ERROR, "%s\n", err)
						return err
					}

				} else if filepath.Ext(location) == ".ww" {
					wwlog.Printf(wwlog.DEBUG, "Found template file: %s\n", location)

					destFile := strings.TrimSuffix(location, ".ww")

					tmpl, err := template.New(path.Base(location)).Funcs(template.FuncMap{"Include": templateFileInclude, "IncludeFromVnfs": templateVnfsFileInclude}).ParseGlob(path.Join(OverlayDir, destFile + ".ww*"))
					if err != nil {
						wwlog.Printf(wwlog.ERROR, "%s\n", err)
						return err
					}

					w, err := os.OpenFile(path.Join(tmpDir, destFile), os.O_RDWR|os.O_CREATE, info.Mode())
					if err != nil {
						wwlog.Printf(wwlog.ERROR, "%s\n", err)
						return err
					}
					defer w.Close()

					err = tmpl.Execute(w, node)
					if err != nil {
						wwlog.Printf(wwlog.ERROR, "%s\n", err)
						return err
					}

				} else if b, _ := regexp.MatchString(`\.ww[a-zA-Z0-9\-\._]*$`, location); b == true {
					wwlog.Printf(wwlog.DEBUG, "Ignoring WW template file: %s\n", location)
				} else {
					wwlog.Printf(wwlog.DEBUG, "Found file: %s\n", location)

					err := util.CopyFile(path.Join(OverlayDir, location), path.Join(tmpDir, location))
					if err != nil {
						wwlog.Printf(wwlog.ERROR, "%s\n", err)
						return err
					}

				}

				return nil
			})
			wwlog.Printf(wwlog.VERBOSE, "Finished generating overlay directory for: %s\n", node.Fqdn)

			cmd := fmt.Sprintf("cd \"%s\"; find . | cpio --quiet -o -H newc -F \"%s\"", tmpDir, OverlayFile)
			wwlog.Printf(wwlog.DEBUG, "RUNNING: %s\n", cmd)
			err = exec.Command("/bin/sh", "-c", cmd).Run()
			if err != nil {
				wwlog.Printf(wwlog.ERROR, "Could not generate runtime image overlay: %s\n", err)
				continue
			}
			wwlog.Printf(wwlog.INFO, "%-35s: Done\n", node.Fqdn)

			wwlog.Printf(wwlog.DEBUG, "Removing temporary directory: %s\n", tmpDir)
			os.RemoveAll(tmpDir)
		}
	}

	wwlog.SetIndent(0)
	return nil
}