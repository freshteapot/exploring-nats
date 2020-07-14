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
	"os"
	"os/signal"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Test nats server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscribe called")
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

		// Subscribe
		topic := "publish.alist"
		if _, err := ec.Subscribe(topic, processAlist); err != nil {
			log.Fatalf("Failed to start subscription on '%s': %v", topic, err)
		}

		topic = "publish.alists_by_user"
		if _, err := ec.Subscribe(topic, processAlistsByUser); err != nil {
			log.Fatalf("Failed to start subscription on '%s': %v", topic, err)
		}

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		select {
		case <-signals:
		}
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}

func processAlist(aList aList) {
	log.Printf("Content: %s, UUID: %s", aList.Content, aList.UUID)
}

func processAlistsByUser(lists aListsByUser) {
	log.Printf("Content: %s, UUID: %s", lists.Content, lists.UserUUID)
}
