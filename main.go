package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jschwinger233/grpcdump/grpchelper/grpcurl"
	"github.com/jschwinger233/grpcdump/handler"
	"github.com/jschwinger233/grpcdump/handler/jsonhandler"
	"github.com/jschwinger233/grpcdump/handler/texthandler"
	"github.com/jschwinger233/grpcdump/parser"
	"github.com/jschwinger233/grpcdump/parser/grpcparser"
	"github.com/jschwinger233/grpcdump/provider"
	"github.com/jschwinger233/grpcdump/provider/pcaprovider"
	"github.com/jschwinger233/grpcdump/provider/sniffprovider"
	"github.com/jschwinger233/grpcdump/version"
	cli "github.com/urfave/cli/v2"
)

func main() {
	var (
		provider provider.Provider
		parser   parser.Parser
		handler  handler.GrpcHandler
	)

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(version.Version())
	}

	app := cli.NewApp()
	app.Flags = flags
	app.Usage = "A tool to sniff and decode gRPC frames"
	app.Version = "--version/-v"
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

		parser, err = grpcparser.New(args.ProtoFilenames, args.ServicePort, args.GuessPaths)
		if err != nil {
			return
		}

		var grpcurlManager *grpcurl.Manager
		if args.WithGrpcurl {
			grpcurlManager = grpcurl.New(args.ProtoFilenames)
		}

		switch args.OutputFormat {
		case Text:
			handler = texthandler.New(grpcurlManager)
		case Json:
			handler = jsonhandler.New(grpcurlManager)
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
