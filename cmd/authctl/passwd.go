package main

import (
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func newPasswdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passwd",
		Short: "Managed user password file",
	}
	cmd.AddCommand(
		newPasswdCreateCmd(),
	)
	return cmd
}

func newPasswdCreateCmd() *cobra.Command {
	var (
		file     string
		groups   []string
		password string
		stdin    bool
	)
	cmd := &cobra.Command{
		Use:   "create <user>",
		Args:  cobra.ExactArgs(1),
		Short: "creates a user entryp",
		Long: `creates a bcrypt hash for the password and returns a line in the form:
<user> bcrypt(<password>) [group1,group2,group3,...]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			user := args[0]
			if password != "" {
				// password set by flag
			} else if stdin {
				pwBytes, err := ioutil.ReadAll(cmd.InOrStdin())
				if err != nil {
					return err
				}
				password = string(pwBytes)
			} else {
				cmd.PrintErr("password: ")
				pwBytes, err := terminal.ReadPassword(0)
				cmd.Println("")
				if err != nil {
					return err
				}
				password = string(pwBytes)
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
			if err != nil {
				return err
			}
			cmd.Println(user, string(hash), strings.Join(groups, ","))
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", file, "password file")
	cmd.Flags().StringVarP(&password, "password", "p", password, "password")
	cmd.Flags().StringSliceVar(&groups, "group", groups, "groups")
	cmd.Flags().BoolVar(&stdin, "stdin", stdin, "read password from stdin")
	return cmd
}
