package main

import (
	"errors"
	"log"
	"os"

	"github.com/jschwinger23/grpcdump/handler"
	"github.com/jschwinger23/grpcdump/handler/grpcurlhandler"
	"github.com/jschwinger23/grpcdump/handler/jsonhandler"
	"github.com/jschwinger23/grpcdump/handler/texthandler"
	"github.com/jschwinger23/grpcdump/parser"
	"github.com/jschwinger23/grpcdump/parser/grpcparser"
	"github.com/jschwinger23/grpcdump/provider"
	"github.com/jschwinger23/grpcdump/provider/pcaprovider"
	"github.com/jschwinger23/grpcdump/provider/sniffprovider"
	cli "github.com/urfave/cli/v2"
)

func main() {
	var (
		provider provider.Provider
		parser   parser.Parser
		handler  handler.GrpcHandler
	)

	app := cli.NewApp()
	app.Flags = flags
	app.Before = func(ctx *cli.Context) (err error) {
		args, err := newArgs(ctx)
		if err != nil {
			return
		}

		switch args.ProvideMethod {
		case BySniff:
			provider = sniffprovider.New(args.Source)
		case ByPcapFile:
			provider = pcaprovider.New(args.Source)
		default:
			return errors.New("provider not specified")
		}

		parser, err = grpcparser.New(args.ProtoFilename, args.GuessMethod)
		if err != nil {
			return
		}

		switch args.OutputFormat {
		case Text:
			handler = texthandler.New(args.Verbose)
		case Json:
			handler = jsonhandler.New(args.Verbose)
		case Grpcurl:
			handler = grpcurlhandler.New(args.Verbose)
		default:
			return errors.New("output format not specified")
		}

		return
	}
	app.Action = func(ctx *cli.Context) (err error) {
		ch, err := provider.PacketStream()
		if err != nil {
			return
		}
		for packet := range ch {
			messages, err := parser.Parse(packet)
			if err != nil {
				return err
			}
			for _, message := range messages {
				if err := handler.Handle(message); err != nil {
					return err
				}
			}
		}
		return
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("failed to run app: %+v", err)
	}
}
