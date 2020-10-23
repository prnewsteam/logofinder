package finder

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"math"

	"github.com/prnewsteam/logofinder/helper"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type Logo struct {
	File *os.File
}

func (l *Logo) Close() {
	l.File.Close()
}

func DefineWidthHeight(filename string, maxWidth uint, maxHeight uint) (uint, uint, error) {

    mw := imagick.NewMagickWand()
    defer mw.Destroy()

    if err := mw.ReadImage(filename); err != nil {
        return 0, 0, err
    }

    origWidth := uint(mw.GetImageWidth())
    origHeight := uint(mw.GetImageHeight())
    horizontal := origWidth > origHeight
    widthRatio := float64(1)
    heightRatio := float64(1)

    if origWidth < maxWidth {
        if horizontal {
            widthRatio = float64(origWidth) / float64(maxWidth)
        } else {
            widthRatio = float64(maxWidth)
        }
    }

    if origHeight < maxHeight {
        if horizontal == false {
            heightRatio = float64(origHeight) / float64(maxHeight)

        } else {
            heightRatio = float64(maxHeight)
        }
    }

    ratio := math.Min(widthRatio, heightRatio)
    return uint(float64(maxWidth) * ratio), uint(float64(maxHeight) * ratio), nil
}

func (l *Logo) Resize(width uint, height uint) (*Logo, error) {

    var err error

    filename := l.File.Name()

    imagick.Initialize()
    defer imagick.Terminate()

    width, height, err = DefineWidthHeight(filename, width, height)

    if err != nil {
        return nil, err
    }

	nFilename := path.Join(
		path.Dir(filename),
		fmt.Sprintf("logo_%dx%d", width, height),
	)

	logo, err := newLogoFromPath(nFilename)
	if err == nil {
		return logo, nil
	}

	nFilename += ".jpg"

	pw := imagick.NewPixelWand()
	pw.SetColor("white")

	mw := imagick.NewMagickWand()
    defer mw.Destroy()

	if path.Ext(filename) == ".svg" {
		rsvg := "rsvg-convert"
		cmd := exec.Command(rsvg, "-d", "300", "-a", "-h", strconv.FormatUint(uint64(height), 10), "-w", strconv.FormatUint(uint64(width), 10), filename)

		file, _ := os.OpenFile(nFilename, os.O_WRONLY|os.O_CREATE, 0777)
		defer file.Close()
		cmd.Stdout = file
		cmd.Run()

	} else {
		imagick.ConvertImageCommand([]string{
			"convert", filename + "[0]",
			"-background", "white",
			"-alpha", "remove",
			"-define", "jpeg:fancy-upsampling=off",
			"-define", "png:compression-filter=5",
            "-define", "png:compression-level=9",
            "-define", "png:compression-strategy=1",
            "-define", "png:exclude-chunk=all",
			"-resize", fmt.Sprintf("%dx%d>", width-5, height-5), nFilename,
		})
	}

	err = mw.ReadImage(nFilename)
	if err != nil {
		panic(err)
	}

	w := int(mw.GetImageWidth())
	h := int(mw.GetImageHeight())
	mw.SetImageBackgroundColor(pw)

	// This centres the original image on the new canvas.
	// Note that the extent's offset is relative to the
	// top left corner of the *original* image, so adding an extent
	// around it means that the offset will be negative
	err = mw.ExtentImage(width, height, -(int(width)-w)/2, -(int(height)-h)/2)
	if err != nil {
		return nil, err
	}

	err = mw.WriteImage(nFilename)
	if err != nil {
		return nil, err
	}

	l.File.Close()
	os.Remove(l.File.Name())

	return newLogoFromPath(nFilename)
}

func (l *Logo) Clear() {
	l.File.Close()
	os.Remove(l.File.Name())
	os.Remove(path.Dir(l.File.Name()));
}

func (l *Logo) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=logo." + path.Ext(l.File.Name()))
	w.WriteHeader(http.StatusOK)
	
	p := make([]byte, 1024)
	for {
		n, err := l.File.Read(p)
		if err == io.EOF{
			break
		}
		w.Write(p[:n])
	}
}

func NewLogoFromUrl(url, domain string) (*Logo, error) {
	ext := strings.Split(path.Ext(url), "?")
	return NewLogoFromUrlWithExtension(url, domain, ext[0])
}

func NewLogoFromUrlWithExtension(url, domain, extension string) (*Logo, error) {
	logoPath := "logo/" + domain + "/logo" + extension
	err := helper.FileDownload(url, domain, logoPath)
	if err != nil {
		return nil, errors.New("Unable to download logo from provided url: " + url)
	}
	return newLogoFromPath(logoPath)
}

func NewLogoFromRaw(domain, content, extension string) (*Logo, error) {
	os.Mkdir("logo/"+domain, 0777)
	file, _ := os.OpenFile("logo/"+domain+"/logo"+extension, os.O_RDWR|os.O_CREATE, 0777)
	file.WriteString(content)
	file.Close()
	file2, _ := os.OpenFile("logo/"+domain+"/logo"+extension, os.O_RDWR, 0777)
	return &Logo{file2}, nil
}

func newLogoFromPath(path string) (*Logo, error) {
	logoPath := helper.FindFileByName(path)
	if logoPath == "" {
		return nil, errors.New("Unable to locate logo by name")
	}

	file, err := os.OpenFile(logoPath, os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	return &Logo{file}, nil
}
