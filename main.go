package main

import (
	"log"
	"os"

	"github.com/jschwinger23/grpcdump/formatter"
	"github.com/jschwinger23/grpcdump/formatter/grpcurlformatter"
	"github.com/jschwinger23/grpcdump/formatter/jsonformatter"
	"github.com/jschwinger23/grpcdump/formatter/textformatter"
	"github.com/jschwinger23/grpcdump/parser"
	"github.com/jschwinger23/grpcdump/provider"
	"github.com/jschwinger23/grpcdump/provider/pcaprovider"
	"github.com/jschwinger23/grpcdump/provider/sniffprovider"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Flags = flags
	app.Action = func(ctx *cli.Context) (err error) {
		args, err := newArgs(ctx)
		if err != nil {
			return
		}

		var provider provider.Provider
		switch args.ProvideMethod {
		case BySniff:
			provider = sniffprovider.New(args.Source)
		case ByPcapFile:
			provider = pcaprovider.New(args.Source)
		}

		var formatter formatter.Formatter
		switch args.OutputFormat {
		case Text:
			formatter = textformatter.New(args.Verbose)
		case Json:
			formatter = jsonformatter.New(args.Verbose)
		case Grpcurl:
			formatter = grpcurlformatter.New(args.Verbose)
		}

		ch, err := provider.PacketStream()
		if err != nil {
			return
		}
		parser := parser.New(args.ProtoFilename, args.GuessMethod)
		for packet := range ch {
			msg, err := parser.Parse(packet)
			if err != nil {
				return err
			}
			formatter.Format(msg)
		}
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("failed to run app: %+v", err)
	}
}
