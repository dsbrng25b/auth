package main

import (
	"fmt"

	"github.com/dvob/auth"
	"github.com/spf13/cobra"
)

func newTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Create random tokens",
	}
	cmd.AddCommand(
		newTokenCreateCmd(),
	)
	return cmd
}

func newTokenCreateCmd() *cobra.Command {
	var (
		randomBytesN = 32
		ascii85      bool
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "creates a token",
		Long:  "create reads a number of random bytes and encodes them. By default base64 encoding is used.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				token string
				err   error
			)
			if ascii85 {
				token, err = auth.GenerateTokenASCII85(randomBytesN)
			} else {
				token, err = auth.GenerateTokenBase64(randomBytesN)
			}
			if err != nil {
				return err
			}
			cmd.SetOut(nil)
			fmt.Println(token)
			return nil
		},
	}
	cmd.Flags().IntVarP(&randomBytesN, "byte", "b", randomBytesN, "The number of random bytes to read")
	cmd.Flags().BoolVar(&ascii85, "ascii85", ascii85, "Encode bytes in ASCII85 encoding")
	return cmd
}
