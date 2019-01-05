package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"

	"github.com/jakobvarmose/crypta/ipfs"
	"github.com/jakobvarmose/crypta/server"
	"github.com/jakobvarmose/crypta/userstore"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Offline bool `long:"offline" description:"Run in offline mode"`
	Dev     bool `long:"dev" description:"Run in development mode"`
	Docker  bool `long:"inside-docker" description:"Special settings for running inside docker"`
}

var opts Options

func main() {
	parser := flags.NewNamedParser("Crypta", flags.Default)
	parser.AddGroup("Options", "Options", &opts)
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		return
	}
	repoPath, err := ipfs.BestKnownPath()
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err := ipfs.NewNode(repoPath, opts.Offline)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer n.Close()
	us, err := userstore.New(path.Join(repoPath, "userstore"))
	if err != nil {
		fmt.Println(err)
		return
	}

	db := NewDatabase(us, n)
	app, err := server.New(n, us, db)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		err = server.Run(app, opts.Dev, opts.Docker)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	fmt.Println()
}
