package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var (
	branch string
	token  string
)

func main() {

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
			var auth *http.BasicAuth
			if token != "" {
				auth = &http.BasicAuth{
					Username: "minigit", // anything except an empty string
					Password: token,
				}
			}
			_, err := git.PlainClone(directory, false, &git.CloneOptions{
				URL:               url,
				Auth:              auth,
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

			var auth *http.BasicAuth
			if token != "" {
				auth = &http.BasicAuth{
					Username: "minigit", // anything except an empty string
					Password: token,
				}
			}
			err = w.Pull(&git.PullOptions{
				Auth:          auth,
				RemoteName:    remote,
				ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
			})
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
	rootCmd.PersistentFlags().StringVar(&token, "ghtoken", "", "github access token")
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.Execute()
}
