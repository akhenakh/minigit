package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var (
	branch   string
	token    string
	format   string
	nopager  bool
	maxCount int
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

	var logCmd = &cobra.Command{
		Use:   "log",
		Short: "git log",
		Long:  "use this command to display commit hashes, git log",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			r, err := git.PlainOpen(path)
			if err != nil {
				log.Fatal(err)
			}

			ref, err := r.Head()
			if err != nil {
				log.Fatal(err)
			}

			// ... retrieves the commit history
			cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
			if err != nil {
				log.Fatal(err)
			}
			count := 0
			// ... just iterates over the commits, printing it
			err = cIter.ForEach(func(c *object.Commit) error {
				if maxCount > 0 && count >= maxCount {
					return nil
				}
				if format != "" {
					if strings.HasPrefix(format, "format:") {
						f := strings.TrimLeft(format[7:], "\"")
						f = strings.TrimRight(f, "\"")
						if f == "%H" {
							fmt.Println(c.Hash)
						}
					}
				} else {
					fmt.Println(c)
				}
				count++
				return nil
			})

		},
	}

	cloneCmd.Flags().StringVarP(&branch, "branch", "b", "master", "branch to clone")
	logCmd.Flags().IntVarP(&maxCount, "max-count", "n", 0, "limit the number of commits to output")
	logCmd.Flags().StringVarP(&format, "pretty", "", "", "format to display logs, only supporting %H")

	var rootCmd = &cobra.Command{Use: "git"}
	rootCmd.PersistentFlags().StringVar(&token, "ghtoken", "", "github access token")
	rootCmd.PersistentFlags().BoolVar(&nopager, "no-pager", true, "this option has no action, output is never paged")
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.Execute()
}
