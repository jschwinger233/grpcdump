package main

import (
	"errors"
	"fmt"
	"strings"

	cli "github.com/urfave/cli/v2"
)

var flags []cli.Flag = []cli.Flag{
	&cli.IntFlag{
		Name:     "service-port",
		Aliases:  []string{"p"},
		Usage:    "grpc service port; e.g. -p 2379",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "interface",
		Aliases:  []string{"i"},
		Usage:    "listen on interface; e.g. -i eth0",
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
		Name:     "guess-path",
		Aliases:  []string{"m"},
		Usage:    "e.g. -m /pb.CoreRPC/WatchServiceStatus,/pb.CoreRPC/SetWorkloadsStatus or -m AUTO",
		Required: false,
	},
	&cli.StringFlag{
		Name:        "output-format",
		Aliases:     []string{"o"},
		DefaultText: "text",
		Usage:       "output format including 'text', 'json', 'grpcurl'",
		Required:    false,
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
	ServicePort   int
	ProtoFilename string
	GuessPaths    []string

	// outputter
	OutputFormat
}

func newArgs(ctx *cli.Context) (args *Args, err error) {
	args = &Args{}

	args.ServicePort = ctx.Int("service-port")

	if iface := ctx.String("interface"); iface != "" {
		args.ProvideMethod = BySniff
		args.Source = fmt.Sprintf("%s:%d", iface, args.ServicePort)
	}

	if pcapFilename := ctx.String("read-file"); pcapFilename != "" {
		if args.ProvideMethod != UnknownMethod {
			return nil, errors.New("sniff-target and read-file cannot be used together")
		}
		args.ProvideMethod = ByPcapFile
		args.Source = pcapFilename
	}

	args.ProtoFilename = ctx.String("proto-file")
	args.GuessPaths = strings.Split(ctx.String("guess-path"), ",")
	if args.GuessPaths[0] == "" {
		args.GuessPaths = []string{}
	}

	switch ctx.String("output-format") {
	case "text":
		args.OutputFormat = Text
	case "json":
		args.OutputFormat = Json
	case "grpcurl":
		args.OutputFormat = Grpcurl
	default:
		args.OutputFormat = Text
	}

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
