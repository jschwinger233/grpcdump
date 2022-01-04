package formatter

import "golang.org/x/net/http2"

type Formatter interface {
	Format(http2.DataFrame)
}
