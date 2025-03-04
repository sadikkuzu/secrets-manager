/*
|    Protect your secrets, protect your sensitive data.
:    Explore VMware Secrets Manager docs at https://vsecm.com/
</
<>/  keep your secrets… secret
>/
<>/' Copyright 2023–present VMware, Inc.
>/'  SPDX-License-Identifier: BSD-2-Clause
*/

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/akamensky/argparse"
	"github.com/vmware-tanzu/secrets-manager/app/sentinel/internal/safe"
)

func main() {
	parser := argparse.NewParser("safe", "Assigns secrets to workloads.")

	list := parseList(parser)
	useKubernetes := parseUseKubernetes(parser)
	deleteSecret := parseDeleteSecret(parser)
	appendSecret := parseAppendSecret(parser)
	namespace := parseNamespace(parser)
	inputKeys := parseInputKeys(parser)
	backingStore := parseBackingStore(parser)
	workload := parseWorkload(parser)
	secret := parseSecret(parser)
	template := parseTemplate(parser)
	format := parseFormat(parser)
	encrypt := parseEncrypt(parser)

	err := parser.Parse(os.Args)
	if err != nil {
		printUsage(parser)
		return
	}

	if *list {
		safe.Get()
		return
	}

	if *namespace == "" {
		*namespace = "default"
	}

	if inputValidationFailure(workload, encrypt, inputKeys, secret, deleteSecret) {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case <-c:
			fmt.Println("Operation was cancelled.")
			cancel()
		}
	}()

	safe.Post(
		ctx,
		*workload, *secret, *namespace, *backingStore, *useKubernetes,
		*template, *format, *encrypt, *deleteSecret, *appendSecret, *inputKeys,
	)
}
