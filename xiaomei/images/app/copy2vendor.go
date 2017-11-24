package app

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func copy2vendorCmd() *cobra.Command {
	var n bool
	cmd := &cobra.Command{
		Use:   `copy2vendor <package-path> ...`,
		Short: `copy the specified packages to vendor dir.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) <= 0 {
				return errors.New(`need at least one package path.`)
			}
			return copy2Vendor(args, n)
		},
	}
	cmd.Flags().BoolVarP(&n, `no-clobber`, `n`, false, `do not overwrite an existing file.`)
	return cmd
}

func copy2Vendor(pkgDirs []string, noClobber bool) error {
	for _, pkgDir := range pkgDirs {
		if err := checkAndCopy(pkgDir, noClobber); err != nil {
			return err
		}
		if err := copy2Vendor(pkgDeps(pkgDir), noClobber); err != nil {
			return err
		}
	}
	return nil
}

func checkAndCopy(pkgDir string, noClobber bool) error {
	goPath := os.Getenv(`GOPATH`)
	vendorDir := path.Join(release.Root(), `../vendor`)

	pkgSrcDir := path.Join(goPath, `src`, pkgDir)
	// package src dir is exists
	if err := checkDir(pkgSrcDir); err != nil {
		return err
	}
	// package src dir not empty
	if ok, err := fs.IsEmptyDir(pkgSrcDir); err != nil {
		return err
	} else if ok {
		return errors.New(pkgSrcDir + ` exists and is empty.`)
	}
	pkgVendorDir := path.Join(vendorDir, pkgDir)
	if err := checkDir(pkgVendorDir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// package vendor dir not exists, and make all dir
		if err = os.MkdirAll(pkgVendorDir, 0775); err != nil {
			return err
		}
	}
	// copy package
	copyPkg(pkgSrcDir, pkgVendorDir, noClobber)
	return nil
}

func copyPkg(pkgSrcDir, pkgVendorDir string, noClobber bool) {
	var flags string
	if noClobber {
		flags = `-ru`
	} else {
		flags = `-r`
	}
	cmd.Run(cmd.O{}, `rsync`, flags, `--exclude=.*`, `--delete`, pkgSrcDir, path.Dir(pkgVendorDir))
}

func checkDir(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New(dir + ` exist and is not a dir.`)
	}
	return nil
}

func pkgDeps(pkgDir string) (result []string) {
	deps, _ := cmd.Run(cmd.O{Output: true}, `go`, `list`, `-e`, `-f`, `'{{join .Deps "\n" }}'`, pkgDir)
	for _, dep := range strings.Split(deps, "\n") {
		if strings.Contains(dep, `.`) && !strings.Contains(dep, `vendor`) {
			result = append(result, dep)
		}
	}
	return
}
