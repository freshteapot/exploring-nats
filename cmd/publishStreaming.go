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
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"gopkg.in/square/go-jose.v2/json"

	stan "github.com/nats-io/stan.go"
)

// publishStreamingCmd represents the publishStreaming command
var publishStreamingCmd = &cobra.Command{
	Use:   "publishStreaming",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("publishStreaming called")
		server := "test-cluster"
		clientID := "chris"
		sc, err := stan.Connect(server, clientID)
		fmt.Println(err)

		// Define the object

		// Publish the message
		inputA, _ := json.Marshal(&aList{Content: "I am a list", UUID: "fake-list-123"})
		if err := sc.Publish("publish.alist", inputA); err != nil {
			log.Fatal(err)
		}

		inputB, _ := json.Marshal(&aListsByUser{Content: "[]", UserUUID: "fake-user-123"})
		if err := sc.Publish("publish.alists_by_user", inputB); err != nil {
			log.Fatal(err)
		}

		// Simple Synchronous Publisher
		sc.Publish("foo", []byte("Hello World")) // does not return until an ack has been received from NATS Streaming

		// Close connection
		sc.Close()
	},
}

func init() {
	rootCmd.AddCommand(publishStreamingCmd)
}
