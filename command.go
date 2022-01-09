package main

import (
	"errors"

	cli "github.com/urfave/cli/v2"
)

var flags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:     "sniff-target",
		Aliases:  []string{"i"},
		Usage:    "interface and port to sniff, incompatible with -r; e.g. -i eth0:2379",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "read-file",
		Aliases:  []string{"r"},
		Usage:    "pcap file to parse, incompatible with -i; e.g. -r packet.pcap",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "proto-file",
		Aliases:  []string{"f"},
		Usage:    "proto file to parse http2 frame; e.g. -f rpc.proto",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "method-for-unknown-stream",
		Aliases:  []string{"m"},
		Usage:    "rpc method to parse response frame whose request method is unknown; e.g. -m Watch",
		Required: false,
	},
	&cli.StringFlag{
		Name:        "output-format",
		Aliases:     []string{"o"},
		DefaultText: "text",
		Usage:       "output format including 'text', 'json', 'grpcurl'",
		Required:    false,
	},
	&cli.BoolFlag{
		Name:     "verbose",
		Aliases:  []string{"v"},
		Usage:    "output http2 frames when verbose on",
		Required: false,
	},
}

type ProvideMethod int

const (
	UnknownMethod ProvideMethod = iota
	BySniff
	ByPcapFile
)

type OutputFormat int

const (
	UnknownFormat OutputFormat = iota
	Text
	Json
	Grpcurl
)

type Args struct {
	// provider
	ProvideMethod
	Source string

	// parser
	ProtoFilename string
	GuessMethod   string

	// outputter
	OutputFormat
	Verbose bool
}

func newArgs(ctx *cli.Context) (args *Args, err error) {
	args = &Args{}

	if sniffTarget := ctx.String("sniff-target"); sniffTarget != "" {
		args.ProvideMethod = BySniff
		args.Source = sniffTarget
	}

	if pcapFilename := ctx.String("read-file"); pcapFilename != "" {
		if args.ProvideMethod != UnknownMethod {
			return nil, errors.New("sniff-target and read-file cannot be used together")
		}
		args.ProvideMethod = ByPcapFile
		args.Source = pcapFilename
	}

	args.ProtoFilename = ctx.String("proto-file")
	args.GuessMethod = ctx.String("method-for-unknown-stream")

	switch ctx.String("output-format") {
	case "text":
		args.OutputFormat = Text
	case "json":
		args.OutputFormat = Json
	case "grpcurl":
		args.OutputFormat = Grpcurl
	}

	args.Verbose = ctx.Bool("verbose")
	return args, args.Validate()
}

func (a *Args) Validate() (err error) {
	if a.ProvideMethod == UnknownMethod {
		return errors.New("either sniff-target or read-file must be set")
	}

	if a.OutputFormat == UnknownFormat {
		return errors.New("output-format must be set")
	}

	return nil
}
