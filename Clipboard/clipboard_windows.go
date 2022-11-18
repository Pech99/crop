package clipboard

import (
	"fmt"
	"image"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

const (
	cFmtBitmap      = 2 // Win+PrintScreen
	cFmtUnicodeText = 13
	cFmtDIBV5       = 17
	// Screenshot taken from special shortcut is in different format (why??), see:
	// https://jpsoft.com/forums/threads/detecting-clipboard-format.5225/
	cFmtDataObject = 49161 // Shift+Win+s, returned from enumClipboardFormats
	gmemMoveable   = 0x0002
)

var (
	user32                     = syscall.MustLoadDLL("user32")
	openClipboard              = user32.MustFindProc("OpenClipboard")
	closeClipboard             = user32.MustFindProc("CloseClipboard")
	emptyClipboard             = user32.MustFindProc("EmptyClipboard")
	setClipboardData           = user32.MustFindProc("SetClipboardData")
	getClipboardSequenceNumber = user32.MustFindProc("GetClipboardSequenceNumber")
	//	getClipboardData           = user32.MustFindProc("GetClipboardData")
	//	isClipboardFormatAvailable = user32.MustFindProc("IsClipboardFormatAvailable")
	//	enumClipboardFormats       = user32.MustFindProc("EnumClipboardFormats")
	//	registerClipboardFormatA   = user32.MustFindProc("RegisterClipboardFormatA")

	kernel32 = syscall.NewLazyDLL("kernel32")
	gLock    = kernel32.NewProc("GlobalLock")
	gUnlock  = kernel32.NewProc("GlobalUnlock")
	gAlloc   = kernel32.NewProc("GlobalAlloc")
	gFree    = kernel32.NewProc("GlobalFree")
	memMove  = kernel32.NewProc("RtlMoveMemory")
)

type bitmapV5Header struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
	RedMask       uint32
	GreenMask     uint32
	BlueMask      uint32
	AlphaMask     uint32
	CSType        uint32
	Endpoints     struct {
		CiexyzRed, CiexyzGreen, CiexyzBlue struct {
			CiexyzX, CiexyzY, CiexyzZ int32 // FXPT2DOT30
		}
	}
	GammaRed    uint32
	GammaGreen  uint32
	GammaBlue   uint32
	Intent      uint32
	ProfileData uint32
	ProfileSize uint32
	Reserved    uint32
}

func writeImg(img *image.RGBA) (<-chan struct{}, error) {
	errch := make(chan error)
	changed := make(chan struct{}, 1)
	go func() {

		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		for {
			r, _, _ := openClipboard.Call(0)
			if r == 0 {
				continue
			}
			break
		}

		err := writeImage(img)
		if err != nil {
			errch <- err
			closeClipboard.Call()
			return
		}

		closeClipboard.Call()

		cnt, _, _ := getClipboardSequenceNumber.Call()
		errch <- nil
		for {
			time.Sleep(time.Second)
			cur, _, _ := getClipboardSequenceNumber.Call()
			if cur != cnt {
				changed <- struct{}{}
				close(changed)
				return
			}
		}
	}()
	err := <-errch
	if err != nil {
		return nil, err
	}
	return changed, nil
}

func writeImage(img *image.RGBA) error {

	r, _, err := emptyClipboard.Call()
	if r == 0 {
		return fmt.Errorf("failed to clear clipboard: %w", err)
	}

	offset := unsafe.Sizeof(bitmapV5Header{})
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	imageSize := 4 * width * height

	data := make([]byte, int(offset)+imageSize)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := int(offset) + 4*(y*width+x)

			r, g, b, a := img.At(x, height-y).RGBA()
			data[idx+2] = uint8(r)
			data[idx+1] = uint8(g)
			data[idx+0] = uint8(b)
			data[idx+3] = uint8(a)
		}
	}

	info := bitmapV5Header{}
	info.Size = uint32(offset)
	info.Width = int32(width)
	info.Height = int32(height)
	info.Planes = 1
	info.Compression = 0 // BI_RGB
	info.SizeImage = uint32(4 * info.Width * info.Height)
	info.RedMask = 0xff0000 // default mask
	info.GreenMask = 0xff00
	info.BlueMask = 0xff
	info.AlphaMask = 0xff000000
	info.BitCount = 32 // we only deal with 32 bpp at the moment.
	info.CSType = 0x73524742
	info.Intent = 4 // LCS_GM_IMAGES

	infob := make([]byte, int(unsafe.Sizeof(info)))
	for i, v := range *(*[unsafe.Sizeof(info)]byte)(unsafe.Pointer(&info)) {
		infob[i] = v
	}
	copy(data[:], infob[:])

	hMem, _, err := gAlloc.Call(gmemMoveable,
		uintptr(len(data)*int(unsafe.Sizeof(data[0]))))
	if hMem == 0 {
		return fmt.Errorf("failed to alloc global memory: %w", err)
	}

	p, _, err := gLock.Call(hMem)
	if p == 0 {
		return fmt.Errorf("failed to lock global memory: %w", err)
	}
	defer gUnlock.Call(hMem)

	memMove.Call(p, uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)*int(unsafe.Sizeof(data[0]))))

	v, _, err := setClipboardData.Call(cFmtDIBV5, hMem)
	if v == 0 {
		gFree.Call(hMem)
		return fmt.Errorf("failed to set text to clipboard: %w", err)
	}

	return nil
}
