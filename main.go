package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "sniff-interface",
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
			Name:     "proto-filename",
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
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("failed to run app: %+v", err)
	}
}

func run(ctx *cli.Context) (err error) {
	return
}
