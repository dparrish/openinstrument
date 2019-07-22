package stream_container

import (
	"sync"

	"golang.org/x/net/context"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

type Container interface {
	// AddStream adds a single stream to the container.
	AddStream(context.Context, *oproto.ValueStream)

	// AddStreams adds all the streams provided on the supplied channel to the container.
	// It will return when the channel is closed, or when the context is Done.
	// The streams written may not be written to disk until after Flush() has been called.
	AddStreams(context.Context, <-chan *oproto.ValueStream)

	// Reader gets a set of streams for a supplied variable match on the returned channel.
	// It will return once all matching streams have been sent to the channel, or when the context is Done.
	// No specific order of the streams is implied.
	Reader(context.Context, *variable.Variable) <-chan *oproto.ValueStream

	// GetAllStreams gets all the streams in the container.
	// It will return once all matching streams have been sent to the channel, or when the context is Done.
	// No specific order of the streams is implied.
	// The streams returned should not be modified unless the object returned by RWLocker is locked.
	GetAllStreams(context.Context) (<-chan *oproto.ValueStream, error)

	// GetUnloggedStreams gets all streams that have not yet been flushed to disk in any way.
	// It will return once all matching streams have been sent to the channel, or when the context is Done.
	// No specific order of the streams is implied.
	// The streams returned are considered at risk.
	// The streams returned must not be modified unless the object returned by RWLocker is locked.
	GetUnloggedStreams(context.Context) (<-chan *oproto.ValueStream, error)

	// GetLoggedStreams gets all streams that have been flushed to disk in a temporary capacity.
	// It will return once all matching streams have been sent to the channel, or when the context is Done.
	// No specific order of the streams is implied.
	// The streams returned are considered safe but are not available in the fastest way.
	// The streams returned must not be modified unless the object returned by RWLocker is locked.
	GetLoggedStreams(context.Context) (<-chan *oproto.ValueStream, error)

	// GetIndexedStreams gets all streams that have been flushed to disk and indexing has completed.
	// It will return once all matching streams have been sent to the channel, or when the context is Done.
	// No specific order of the streams is implied.
	// The streams returned are considered safe and highly available.
	// The streams returned must not be modified unless the object returned by RWLocker is locked.
	GetIndexedStreams(context.Context) (<-chan *oproto.ValueStream, error)

	// Flush writes all at-risk streams to disk.
	// If this returns no error, the all the streams held by this container have been written to persistent storage.
	// If there is an error writing to storage, or the context is cancelled during the flush, an error will be returned and the flush aborted.
	Flush(context.Context) error

	// NumStreams gets the number of streams managed by this container.
	NumStreams() uint32

	// NumValues gets the total number of values in all streams combined.
	NumValues() uint32

	// Locker returns an that is locked while background data manipulation is being performed.
	// The object returned can be used to keep a consistent view of data for read access.
	Locker() sync.Locker
}
