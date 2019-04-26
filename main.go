package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {
	var branch string
	var cloneCmd = &cobra.Command{
		Use:   "clone",
		Short: "clone git repo",
		Long:  "use this command to clone git repo, git clone URL directory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			directory := "."
			if len(args) > 1 {
				directory = args[1]
			}
			_, err := git.PlainClone(directory, false, &git.CloneOptions{
				URL:               url,
				ReferenceName:     plumbing.ReferenceName("refs/heads/" + branch),
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})
			if err != nil {
				log.Fatal(err)
			}

		},
	}
	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "pull git repo",
		Long:  "use this command to pull, git remote branch",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			remote := args[0]
			branch := args[1]
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			r, err := git.PlainOpen(path)
			if err != nil {
				log.Fatal(err)
			}

			w, err := r.Worktree()
			if err != nil {
				log.Fatal(err)
			}

			err = w.Pull(&git.PullOptions{RemoteName: remote, ReferenceName: plumbing.ReferenceName("refs/heads/" + branch)})
			if err == git.NoErrAlreadyUpToDate {
				return
			}
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cloneCmd.Flags().StringVarP(&branch, "branch", "b", "master", "branch to clone")

	var rootCmd = &cobra.Command{Use: "git"}
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.Execute()
}
