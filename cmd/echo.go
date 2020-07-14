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
	"runtime"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("echo called")
		opts := []nats.Option{nats.Name("NATS Echo Service")}
		opts = setupConnOptions(opts)

		urls := "0.0.0.0"
		nc, err := nats.Connect(urls, opts...)
		if err != nil {
			log.Fatal(err)
		}

		subj, i := args[0], 0

		nc.QueueSubscribe(subj, "echo", func(msg *nats.Msg) {
			i++
			if msg.Reply != "" {
				printMsg(msg, i)
				nc.Publish(msg.Reply, msg.Data)
			}
		})
		nc.Flush()

		if err := nc.LastError(); err != nil {
			log.Fatal(err)
		}

		log.Printf("Echo Service listening on [%s]\n", subj)

		// Now handle signal to terminate so we cam drain on exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)

		go func() {
			// Wait for signal
			<-c
			log.Printf("<caught signal - draining>")
			nc.Drain()
		}()

		runtime.Goexit()
	},
}

func init() {
	rootCmd.AddCommand(echoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// echoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// echoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		if !nc.IsClosed() {
			log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
		}
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		if !nc.IsClosed() {
			log.Fatal("Exiting: no servers available")
		} else {
			log.Fatal("Exiting")
		}
	}))
	return opts
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Echoing to [%s]: %q", i, m.Reply, m.Data)
}
