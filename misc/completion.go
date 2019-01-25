package misc

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func bashCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   `bash-completion`,
		Short: `setup bash completion script.`,
		RunE: release.NoArgCall(func() error {
			dir := getBashCompletionDir()
			if dir == "" {
				return nil
			}
			var buf bytes.Buffer
			if err := rootCmd.GenBashCompletion(&buf); err != nil {
				return err
			}
			cmd.SudoWriteFile(dir+`/xiaomei`, &buf)
			fmt.Printf(`Run the following cmd to make completion take effect immediately:
      source %s
`, strings.TrimSuffix(dir, ".d"))
			return nil
		}),
	}
}

func getBashCompletionDir() string {
	const dir1 = `/etc/bash_completion.d`
	const dir2 = `/usr/local/etc/bash_completion.d`
	if fs.IsDir(dir1) {
		return dir1
	}
	if fs.IsDir(dir2) {
		return dir2
	}
	fmt.Printf("Neither %s nor %s exists.\n", dir1, dir2)
	if runtime.GOOS == `darwin` {
		fmt.Println(`Run the following cmd to install bash completion first:
    brew install bash-completion
then run "xiaomei bash-completion" again.`)
	}
	return ``
}
