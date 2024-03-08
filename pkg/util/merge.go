package util

import (
	"fmt"
	"os"
	"path/filepath"

	gitadd "github.com/ldez/go-git-cmd-wrapper/v2/add"
	gitcheckout "github.com/ldez/go-git-cmd-wrapper/v2/checkout"
	gitcommit "github.com/ldez/go-git-cmd-wrapper/v2/commit"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	gitinit "github.com/ldez/go-git-cmd-wrapper/v2/init"
	gitmerge "github.com/ldez/go-git-cmd-wrapper/v2/merge"
	cp "github.com/otiai10/copy"
)

// Merges two directories together using git as a middleman
func MergeRecursively(dest string, srcs []string) error {
	// create `dest`
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		return err
	}
	if err := os.Mkdir(dest, 0o755); err != nil {
		return err
	}

	// init
	if _, err := git.Init(
		gitinit.Bare, gitinit.Quiet,
		gitinit.Directory(dest),
	); err != nil {
		return err
	}

	// create main branch
	if _, err := git.Checkout(
		gitcheckout.NewBranchForce("main"),
	); err != nil {
		return err
	}

	// merge
	for _, src := range srcs {
		// copy branch
		branchName := fmt.Sprintf("src-%s", filepath.Base(src))
		if err := copyAsBranch(src, dest, branchName); err != nil {
			return err
		}

		// checkout main
		if _, err := git.Checkout(
			gitcheckout.Branch("main"),
		); err != nil {
			return err
		}

		// merge
		if _, err := git.Merge(
			gitmerge.StrategyOption("theirs"),
			gitmerge.AllowUnrelatedHistories,
		); err != nil {
			return err
		}
	}

	// cleanup `.git` repo
	if err := os.RemoveAll(filepath.Join(dest, ".git")); err != nil {
		return err
	}

	return nil
}

// Copies all files from `src` to `dest` as a new branch and commit them.
func copyAsBranch(src string, dest string, branchName string) error {
	// create branch
	if _, err := git.Checkout(
		gitcheckout.NewBranchForce(branchName),
	); err != nil {
		return err
	}

	// copy
	if err := cp.Copy(src, dest); err != nil {
		return err
	}

	// add & commit
	if _, err := git.Add(
		gitadd.All,
	); err != nil {
		return err
	}
	if _, err := git.Commit(
		gitcommit.Message(fmt.Sprintf("merge `%s`", branchName)),
	); err != nil {
		return err
	}

	return nil
}
