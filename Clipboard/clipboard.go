package clipboard

import (
	"fmt"
	"image"
	"os"
	"sync"
)

var (
	// activate only for running tests.
	debug = false
	//errUnavailable = errors.New("clipboard unavailable")
	//errUnsupported = errors.New("unsupported format")
)

// Format represents the format of clipboard data.
type Format int

// All sorts of supported clipboard data
const (
	// FmtText indicates plain text clipboard format
	FmtText Format = iota
	// FmtImage indicates image/png clipboard format
	FmtImage
)

// Due to the limitation on operating systems (such as darwin),
// concurrent read can even cause panic, use a global lock to
// guarantee one read at a time.
var lock = sync.Mutex{}

// Write writes a given buffer to the clipboard in a specified format.
// Write returned a receive-only channel can receive an empty struct
// as a signal, which indicates the clipboard has been overwritten from
// this write.
// If format t indicates an image, then the given buf assumes
// the image data is PNG encoded.
func WriteImmage(img *image.RGBA) <-chan struct{} {
	lock.Lock()
	defer lock.Unlock()

	changed, err := writeImg(img)
	if err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "write to clipboard err: %v\n", err)
		}
		return nil
	}
	return changed
}
