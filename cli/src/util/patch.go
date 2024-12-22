package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/SKevo18/gopatch"
	cp "github.com/otiai10/copy"
)

func PatchModFiles(rootDir string, outputDir string, patchFilePaths []string) error {
	patchLines, err := gopatch.ReadPatchFiles(patchFilePaths)
	if err != nil {
		return err
	}

	// TODO: apply config, if found

	if len(patchLines) == 0 {
		return fmt.Errorf("no valid patch files found in mod paths for root `%s`", rootDir)
	}
	if err := gopatch.PatchDir(rootDir, outputDir, patchLines); err != nil {
		return err
	}

	return nil
}

// Copies mod files from `modRootPaths` (looks up "source" directory here)
// to `outputDir`. Patch files are copied to `<outputDir>/_patches`.
// If a `config.json` file is found in the mod root, it will be used to render
// the `mod.gopatch` file as a template.
func CopyModFiles(modRootPaths []string, outputDir string) error {
	for order, modPath := range modRootPaths {
		if _, err := os.Stat(modPath); os.IsNotExist(err) {
			return fmt.Errorf("mod path `%s` does not exist", modPath)
		}

		// copy patch file
		patchFile := filepath.Join(modPath, "mod.gopatch")
		if _, err := os.Stat(patchFile); !os.IsNotExist(err) {
			// copy patch file to output dir
			patchDir := filepath.Join(outputDir, "_patches")
			if err := os.MkdirAll(patchDir, 0o755); err != nil {
				return fmt.Errorf("fatal error while creating root dir for patches: %s", err)
			}
			patchFileOut := filepath.Join(patchDir, fmt.Sprintf("%d_%s.gopatch", order, filepath.Base(modPath)))
			if err := cp.Copy(patchFile, patchFileOut); err != nil {
				return fmt.Errorf("fatal error while copying patch file: %s", err)
			}

			configFile := filepath.Join(modPath, "config.json")
			if _, err := os.Stat(configFile); !os.IsNotExist(err) {
				// render original patch from config, if found
				configData, err := readConfigFile(configFile)
				if err != nil {
					return fmt.Errorf("fatal error while reading config file: %s", err)
				}
				renderedPatch, err := renderTemplate(patchFile, configData)
				if err != nil {
					return fmt.Errorf("fatal error while rendering patch file: %s", err)
				}

				// write rendered to out patch
				if err := os.WriteFile(patchFileOut, []byte(renderedPatch), 0o644); err != nil {
					return fmt.Errorf("fatal error while writing rendered patch file: %s", err)
				}
			}
		}

		// source dir is optional (if there are only patches)
		sourceDir := filepath.Join(modPath, "source")
		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			continue
		}
		// copy
		if err := cp.Copy(sourceDir, outputDir); err != nil {
			return fmt.Errorf("fatal error while copying mods: %s", err)
		}
	}

	return nil
}

func readConfigFile(configPath string) (configData map[string]interface{}, err error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(configFile, &configData); err != nil {
		return nil, err
	}

	return configData, nil
}

func renderTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", err
	}

	return rendered.String(), nil
}
