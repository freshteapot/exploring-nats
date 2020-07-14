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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	stan "github.com/nats-io/stan.go"
	"github.com/spf13/cobra"
)

// subscriberStreamingCmd represents the subscriberStreaming command
var subscriberStreamingCmd = &cobra.Command{
	Use: "subscriberStreaming",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscriberStreaming called")

		server := "test-cluster"
		clientID := args[0]
		sc, _ := stan.Connect(server, clientID)

		topic := "publish.alist"
		if _, err := sc.Subscribe(topic, processAlistFromBytes, stan.DurableName("my-durable")); err != nil {
			log.Fatalf("Failed to start subscription on '%s': %v", topic, err)
		}

		topic = "publish.alists_by_user"
		if _, err := sc.Subscribe(topic, processAlistsByUserFromBytes, stan.DeliverAllAvailable()); err != nil {
			log.Fatalf("Failed to start subscription on '%s': %v", topic, err)
		}

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		select {
		case <-signals:
		}

		sc.Close()
	},
}

func init() {
	rootCmd.AddCommand(subscriberStreamingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subscriberStreamingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subscriberStreamingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func processAlistFromBytes(stanMsg *stan.Msg) {
	var aList aList
	json.Unmarshal(stanMsg.Data, &aList)
	log.Printf("Content: %s, UUID: %s", aList.Content, aList.UUID)
}

func processAlistsByUserFromBytes(stanMsg *stan.Msg) {
	var lists aListsByUser
	json.Unmarshal(stanMsg.Data, &lists)
	log.Printf("Content: %s, UUID: %s", lists.Content, lists.UserUUID)
}
