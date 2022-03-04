package misc

import (
	"fmt"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func mergeToCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     `merge-to <target branch>`,
		Aliases: []string{`mt`},
		Short:   `merge current branch to target branch, and then push target branch to remote.`,
		RunE: release.Arg1Call("", func(targetBranch string) error {
			return mergeTo(targetBranch)
		}),
	}
	return cmd
}

func mergeTo(targetBranch string) error {
	// "--" after checkout ensure it checkouts a branch, instread of a file or directory.
	_, err := cmd.Run(cmd.O{}, "bash", "-c", fmt.Sprintf(`
set -ex
git checkout %s --
upstream=$(git rev-parse @{upstream} 2>/dev/null || true)
if test -n "$upstream"; then
	git pull --no-edit
elif git rev-parse remotes/origin/%s >/dev/null 2>&1; then
	git branch -u remotes/origin/%s
	git pull --no-edit
	upstream=yes
fi
git merge --no-edit -
if test -n "$upstream"; then
	git push
else
	git push -u origin %s
fi
git checkout - --
`, targetBranch, targetBranch, targetBranch, targetBranch),
	)
	return err
}
