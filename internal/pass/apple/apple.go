package apple

import _ "image/png"
import _ "image/gif"
import _ "image/jpeg"

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ghophp/pkpass-go"
	"github.com/google/uuid"
	"github.com/llumus/x-plugin-wallet-passes/internal/components"
	"github.com/nfnt/resize"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

type Pass struct {
	webServiceURL  string
	teamIdentifier string
	basePath       string
	tmpPath        string
}

func NewPassbook(webServiceURL string, teamIdentifier string, basePath string, tmpPath string) *Pass {
	return &Pass{
		webServiceURL:  webServiceURL,
		teamIdentifier: teamIdentifier,
		basePath:       basePath,
		tmpPath:        tmpPath,
	}
}

func (a *Pass) CreatePass(_ context.Context, request components.CreatePassbookRequestObject) (string, string, error) {
	var passId = uuid.New().String()
	var passName = "WalletPass_" + passId
	var authToken = uuid.New().String()

	tempDir := filepath.Join(a.tmpPath, passName+".pass")
	err := os.MkdirAll(tempDir, 0777)
	if err != nil {
		return "", "", err
	}

	if request.Body.IconImage != "" {
		err := a.downloadImage(request.Body.IconImage, filepath.Join(tempDir, "icon.png"), 90)
		if err != nil {
			return "", "", fmt.Errorf("error downloading icon image (max 3MB type PNG): %s", err)
		}
	}

	if request.Body.LogoImage != nil && *request.Body.LogoImage != "" {
		err := a.downloadImage(*request.Body.LogoImage, filepath.Join(tempDir, "logo.png"), 160)
		if err != nil {
			return "", "", fmt.Errorf("error downloading logo image (max 3MB type PNG): %s", err)
		}
	}

	if request.Body.ThumbnailImage != nil && *request.Body.ThumbnailImage != "" {
		err := a.downloadImage(*request.Body.ThumbnailImage, filepath.Join(tempDir, "thumbnail.png"), 180)
		if err != nil {
			return "", "", fmt.Errorf("error downloading thumbnail image (max 3MB type PNG): %s", err)
		}
	}

	if request.Body.StripImage != nil && *request.Body.StripImage != "" {
		err := a.downloadImage(*request.Body.StripImage, filepath.Join(tempDir, "strip.png"), 375)
		if err != nil {
			return "", "", fmt.Errorf("error downloading strip image (max 3MB type PNG): %s", err)
		}
	}

	if request.Body.ForegroundColor != "" && !isValidRGBColor(request.Body.ForegroundColor) {
		return "", "", fmt.Errorf("invalid foreground color %s", request.Body.ForegroundColor)
	}

	if request.Body.LabelColor != "" && !isValidRGBColor(request.Body.LabelColor) {
		return "", "", fmt.Errorf("invalid label color %s", request.Body.LabelColor)
	}

	if request.Body.BackgroundColor != "" && !isValidRGBColor(request.Body.BackgroundColor) {
		return "", "", fmt.Errorf("invalid background color %s", request.Body.BackgroundColor)
	}

	passData := map[string]interface{}{
		"formatVersion":       1,
		"passTypeIdentifier":  "pass.com.wallet-passes",
		"serialNumber":        passId,
		"webServiceURL":       a.webServiceURL,
		"authenticationToken": authToken,
		"teamIdentifier":      a.teamIdentifier,
		"organizationName":    request.Body.OrganizationName,
		"description":         request.Body.Description,
		"foregroundColor":     request.Body.ForegroundColor,
		"labelColor":          request.Body.LabelColor,
		"backgroundColor":     request.Body.BackgroundColor,
		request.Body.Type:     request.Body.Fields,
	}

	if request.Body.Barcode != nil {
		passData["barcode"] = request.Body.Barcode
	}

	if request.Body.LogoText != nil && *request.Body.LogoText != "" {
		passData["logoText"] = *request.Body.LogoText
	}

	// Convert the constructed map to JSON bytes
	jsonBytes, err := json.Marshal(passData)
	if err != nil {
		return "", "", err
	}

	err = os.WriteFile(filepath.Join(tempDir, "pass.json"), jsonBytes, 0644)
	if err != nil {
		return "", "", err
	}

	cert, err := os.Open(filepath.Join(a.basePath, "assets/certificates", "pass.p12"))
	if err != nil {
		return "", "", err
	}
	defer cert.Close()

	r, err := pkpass.New(tempDir, "", cert)
	if err != nil {
		return "", "", fmt.Errorf("error generating pass %s", err)
	}

	passPath := filepath.Join(a.basePath, "tmp", passName+".pkpass")

	f, err := os.Create(passPath)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return "", "", err
	}

	return passId, passPath, nil
}

const MaxFileSize = 3 * 1024 * 1024 // 3MB

func (a *Pass) downloadImage(url, filePath string, maxSize int) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the Content-Length header to prevent downloading large files
	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		size, err := strconv.Atoi(contentLength)
		if err == nil && size > MaxFileSize {
			return fmt.Errorf("file size %d exceeds maximum allowed size", size)
		}
	}

	contentType := resp.Header.Get("Content-Type")
	var img image.Image

	if strings.Contains(contentType, "image/svg+xml") {
		// Decode SVG
		icon, err := oksvg.ReadIconStream(resp.Body)
		if err != nil {
			return err
		}
		w, h := icon.ViewBox.W, icon.ViewBox.H
		rgba := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
		icon.Draw(rasterx.NewDasher(int(w), int(h), rasterx.NewScannerGV(int(w), int(h), rgba, rgba.Bounds())), 1)
		img = rgba
	} else {
		// Decode other image types
		img, _, err = image.Decode(resp.Body)
		if err != nil {
			return err
		}
	}

	bounds := img.Bounds()
	width := bounds.Dx()

	var resizedImg image.Image
	if width > maxSize {
		resizedImg = resize.Resize(uint(maxSize), 0, img, resize.Lanczos3)
	} else {
		resizedImg = img
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = png.Encode(outFile, resizedImg)
	if err != nil {
		return err
	}

	return nil
}

var regexRgb = regexp.MustCompile(`^rgb\((\d+),\s*(\d+),\s*(\d+)\)$`)

// isValidRGBColor validates if a string is a color in the format "rgb(0,0,0)"
func isValidRGBColor(color string) bool {
	// Regular expression to match the rgb format
	matches := regexRgb.FindStringSubmatch(color)

	// If the regular expression doesn't match, return false
	if matches == nil {
		return false
	}

	// Validate that each component is within 0-255
	for i := 1; i <= 3; i++ {
		num, err := strconv.Atoi(matches[i])
		if err != nil {
			return false
		}
		if num < 0 || num > 255 {
			return false
		}
	}

	return true
}
