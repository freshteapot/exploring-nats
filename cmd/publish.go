/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		server := "0.0.0.0"
		nc, err := nats.Connect(server)
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()

		ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
		if err != nil {
			log.Fatal(err)
		}
		defer ec.Close()

		// Define the object

		// Publish the message
		if err := ec.Publish("publish.alist", &aList{Content: "I am a list", UUID: "fake-list-123"}); err != nil {
			log.Fatal(err)
		}

		if err := ec.Publish("publish.alists_by_user", &aListsByUser{Content: "[]", UserUUID: "fake-user-123"}); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// publishCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// publishCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type aList struct {
	Content string `json:"content"`
	UUID    string `json:"uuid"`
}

type aListsByUser struct {
	Content  string `json:"content"`
	UserUUID string `json:"user_uuid"`
}
