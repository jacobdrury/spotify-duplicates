package pagination

import (
	"github.com/zmb3/spotify/v2"
)

type PageOptions struct {
	Offset int
	Limit  int
}

func NewPageOptions(limit int, offset int) *PageOptions {
	return &PageOptions{
		Offset: offset,
		Limit:  limit,
	}
}

func NewDefaultPageOptions() *PageOptions {
	return &PageOptions{
		Offset: 0,
		Limit:  50,
	}
}

type PaginatedEndpoint[T any] interface {
	RequestData(*spotify.Client, *PageOptions, chan T) bool
}

func ConsumePaginatedEndpoint[T any](client *spotify.Client, paginatedEndpoint PaginatedEndpoint[T], options *PageOptions) <-chan T {
	ch := make(chan T)

	go func() {
		for {
			hasNext := paginatedEndpoint.RequestData(client, options, ch)
			options.Offset += options.Limit

			if !hasNext {
				close(ch)
				return
			}
		}
	}()

	return ch
}
