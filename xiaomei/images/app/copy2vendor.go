package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func copy2vendorCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   `copy2vendor [<package-path>] ...`,
		Short: `copy the specified packages to vendor dir.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 && !all {
				return errors.New(`no package path provided.`)
			}
			return copy2Vendor(args)
		},
	}
	cmd.Flags().BoolVarP(&all, `all`, `a`, false, `copy all dependences.`)
	return cmd
}

func copy2Vendor(pkgs []string) error {
	if len(pkgs) == 0 {
		pkgs = getDeps(1)
	}
	goSrcDir := fs.GetGoSrcPath()
	vendorDir := filepath.Join(release.Root(), `../vendor`)
	for _, pkg := range pkgs {
		if err := syncGoFiles(filepath.Join(goSrcDir, pkg), filepath.Join(vendorDir, pkg)); err != nil {
			return err
		}
	}
	return nil
}

func syncGoFiles(srcDir, destDir string) error {
	srcFiles, err := filepath.Glob(srcDir + "/*.go")
	if err != nil {
		return err
	}
	if fs.Exist(destDir) {
		if err := removeRedundantDestFiles(srcFiles, destDir); err != nil {
			return err
		}
	} else if err = os.MkdirAll(destDir, 0775); err != nil {
		return err
	}
	for _, src := range srcFiles {
		dst := strings.Replace(src, srcDir, destDir, 1)
		if err := fs.Copy(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func removeRedundantDestFiles(srcFiles []string, destDir string) error {
	goFiles := map[string]bool{}
	for _, srcFile := range srcFiles {
		goFiles[filepath.Base(srcFile)] = true
	}

	destFiles, err := filepath.Glob(destDir + "/*.go")
	if err != nil {
		return err
	}

	for _, destFile := range destFiles {
		if !goFiles[filepath.Base(destFile)] {
			if err := os.Remove(destFile); err != nil {
				return err
			}
		}
	}
	return nil
}
