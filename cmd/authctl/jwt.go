package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func newJWTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jwt",
		Short: "Managed JWT tokens",
	}
	cmd.AddCommand(
		newJWTCreateCmd(),
		newJWTVerifyCmd(),
		newJWTShowCmd(),
	)
	return cmd
}

func newJWTCreateCmd() *cobra.Command {
	var (
		claims = jwt.Claims{
			Issuer:   "authctl",
			Audience: []string{"authctl"},
		}
		groups  []string
		keyFile string
		key     string
		file    string
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a jwt token",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				rawKey []byte
				output io.Writer
				err    error
			)
			// key
			if key != "" {
				rawKey = []byte(key)
			} else if keyFile != "" {
				rawKey, err = ioutil.ReadFile(keyFile)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("no key provided. use either --key oder --key-file")
			}

			if file != "" {
				f, err := os.Create(file)
				if err != nil {
					return err
				}
				f.Close()
				output = f
			} else {
				output = cmd.OutOrStdout()
			}

			claims.Expiry = jwt.NewNumericDate(time.Now().Add(time.Second * 60 * 24 * 30))
			claims.NotBefore = jwt.NewNumericDate(time.Now())

			privateCl := struct {
				Groups []string `json:"groups,omitempty"`
			}{
				groups,
			}

			signingKey := jose.SigningKey{Algorithm: jose.HS256, Key: rawKey}
			sig, err := jose.NewSigner(signingKey, (&jose.SignerOptions{}).WithType("JWT"))
			if err != nil {
				return err
			}
			raw, err := jwt.Signed(sig).Claims(claims).Claims(privateCl).CompactSerialize()
			if err != nil {
				return err
			}

			_, err = output.Write([]byte(raw))
			return err
		},
	}
	cmd.Flags().StringVar(&claims.Subject, "subject", claims.Subject, "Subject of the JWT token")
	//cmd.Flags().StringSliceVar(&aud, "audience", claims.Audience, "Audiences of the JWT token")
	cmd.Flags().StringVar(&claims.Issuer, "issuer", claims.Issuer, "Issuer of the JWT token")
	cmd.Flags().StringVarP(&file, "file", "f", file, "File to write the JWT token to")
	cmd.Flags().StringVarP(&key, "key", "k", key, "HMAC 256 secret")
	cmd.Flags().StringVar(&keyFile, "key-file", keyFile, "File to read the HMAC 256 key from")
	return cmd
}

func newJWTVerifyCmd() *cobra.Command {
	var (
		file    string
		keyFile string
		key     string
	)
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "verify a jwt token",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				rawToken []byte
				rawKey   []byte
				err      error
			)

			// key
			if key != "" {
				rawKey = []byte(key)
			} else if keyFile != "" {
				rawKey, err = ioutil.ReadFile(keyFile)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("no key provided. use either --key oder --key-file")
			}

			// token
			if len(args) > 0 {
				rawToken = []byte(args[0])
			} else if file != "" {
				rawToken, err = ioutil.ReadFile(file)
				if err != nil {
					return err
				}
			} else {
				rawToken, err = ioutil.ReadAll(cmd.InOrStdin())
				if err != nil {
					return err
				}
			}

			token, err := jwt.ParseSigned(string(rawToken))
			if err != nil {
				return err
			}
			out := map[string]interface{}{}
			if err := token.Claims(rawKey, &out); err != nil {
				return err
			}

			json, err := json.MarshalIndent(out, "", "  ")
			cmd.Println(string(json))
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", file, "File to read the JWT token from")
	cmd.Flags().StringVarP(&key, "key", "k", key, "HMAC 256 secret")
	cmd.Flags().StringVar(&keyFile, "key-file", keyFile, "File to read the HMAC 256 key from")
	return cmd
}

func newJWTShowCmd() *cobra.Command {
	var (
		file string
	)
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show a jwt token",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				rawToken []byte
				err      error
			)
			// token
			if len(args) > 0 {
				rawToken = []byte(args[0])
			} else if file != "" {
				rawToken, err = ioutil.ReadFile(file)
				if err != nil {
					return err
				}
			} else {
				rawToken, err = ioutil.ReadAll(cmd.InOrStdin())
				if err != nil {
					return err
				}
			}

			parts := bytes.Split(rawToken, []byte("."))
			if len(parts) != 3 {
				return fmt.Errorf("failed to split JWT token")
			}
			payload, err := base64.RawStdEncoding.DecodeString(string(parts[1]))
			if err != nil {
				return err
			}

			out := map[string]interface{}{}
			err = json.Unmarshal(payload, &out)
			if err != nil {
				return err
			}
			json, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				return err
			}
			cmd.Println(string(json))
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", file, "File to read the JWT token from")
	return cmd
}
