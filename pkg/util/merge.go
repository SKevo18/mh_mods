package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gitadd "github.com/ldez/go-git-cmd-wrapper/v2/add"
	gitcheckout "github.com/ldez/go-git-cmd-wrapper/v2/checkout"
	gitcommit "github.com/ldez/go-git-cmd-wrapper/v2/commit"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	gitinit "github.com/ldez/go-git-cmd-wrapper/v2/init"
	gitmerge "github.com/ldez/go-git-cmd-wrapper/v2/merge"
	"github.com/ldez/go-git-cmd-wrapper/v2/types"
	cp "github.com/otiai10/copy"
)

func MergeModsRecursivelyGit(dest string, srcs []string, debug bool) error {
	if _, err := git.Init(gitinit.Directory(dest)); err != nil {
		return err
	}
    log.Printf("Initialized empty git repository in `%s`...", dest)

    runGitInOpt := runGitIn(dest, debug)
	if _, err := git.Checkout(gitcheckout.NewBranch("main"), runGitInOpt); err != nil {
		log.Printf("Main branch already exists or error creating it: %v", err)
	}
	if _, err := git.Commit(
		gitcommit.Message("Initial commit"),
		gitcommit.AllowEmpty,
		runGitInOpt,
	); err != nil {
		log.Printf("Error creating initial commit: %v", err)
	}

	for i, src := range srcs {
		branchName := fmt.Sprintf("src-%d", i + 1)

		// new branch for each modified src
		if _, err := git.Checkout(
			gitcheckout.NewBranch(branchName),
			runGitInOpt,
		); err != nil {
			return err
		}

		// copy mod files
		if err := cp.Copy(src, dest); err != nil {
			return err
		}

		// add & commit mod
		if _, err := git.Add(gitadd.All, runGitInOpt); err != nil {
			return err
		}
		if _, err := git.Commit(
			gitcommit.Message(fmt.Sprintf("Apply changes from path `%s`", src)),
			runGitInOpt,
		); err != nil {
			return err
		}

		// checkout main
		if _, err := git.Checkout(
			gitcheckout.Branch("main"),
			runGitInOpt,
		); err != nil {
			return err
		}

		// merge
		if _, err := git.Merge(
			gitmerge.Commits(branchName),
			gitmerge.StrategyOption("theirs"),
			gitmerge.NoFf,
			gitmerge.AllowUnrelatedHistories,
			runGitInOpt,
		); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(filepath.Join(dest, ".git")); err != nil {
		return fmt.Errorf("failed to remove `.git` directory: %v", err)
	}

	return nil
}

func runGitIn(path string, debug bool) types.Option {
	return git.CmdExecutor(
		func(ctx context.Context, name string, _ bool, args ...string) (string, error) {
			if debug {
				log.Printf("(%s) git %s", path, strings.Join(args, " "))
			}

			cmd := exec.CommandContext(ctx, name, args...)
			cmd.Dir = path
			output, err := cmd.CombinedOutput()

			if debug {
				log.Printf("git output: %s", output)
			}
			return string(output), err
		},
	)
}
