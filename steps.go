package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
)

func ListContributors(repo string, token string) ([]ContributorInfo, error) {
	// Connect endpoint URL
	endpointUrl := fmt.Sprintf("https://api.github.com/repos/%s/contributors", repo)

	// Prepare request
	req, err := http.NewRequest("GET", endpointUrl, nil)
	if err != nil {
		// Failed to create contributors request
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set request headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Do request
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		// Failed to do request
		return nil, fmt.Errorf("sending request: %w", err)
	}

	defer res.Body.Close()

	// Check response
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode response
	var contributors []ContributorInfo
	err = json.NewDecoder(res.Body).Decode(&contributors)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return contributors, nil
}

func FilterUsersOnly(contributors []ContributorInfo) []ContributorInfo {
	var userContributors []ContributorInfo

	for _, contributor := range contributors {
		if contributor.Type == "User" &&
			contributor.Login != "fossabot" { // They are using a user account for bot actions
			userContributors = append(userContributors, contributor)
		}
	}

	return userContributors
}

func DownloadInfo(contributors []ContributorInfo) []ContributorInfoDownload {
	contributorsDl := make([]ContributorInfoDownload, len(contributors))

	for index, contributor := range contributors {
		log.Printf("Download %d of %d...", index+1, len(contributors))

		contributorsDl[index].Login = contributor.Login
		contributorsDl[index].Id = contributor.Id

		// Download image
		avatarImage, err := (&http.Client{}).Get(contributor.AvatarUrl)
		if err != nil {
			log.Printf("Error downloading avatar image %s: %v", contributor.AvatarUrl, err)
			continue
		}

		avatarBytes, err := io.ReadAll(avatarImage.Body)
		if err != nil {
			log.Printf("Error reading avatar image %s: %v", contributor.AvatarUrl, err)
			continue
		}

		avatarImage.Body.Close()

		var img image.Image

		contentType := http.DetectContentType(avatarBytes)
		switch contentType {
		case "image/png":
			img, err = png.Decode(bytes.NewReader(avatarBytes))
		case "image/jpeg":
			img, err = jpeg.Decode(bytes.NewReader(avatarBytes))
		case "image/gif":
			img, err = gif.Decode(bytes.NewReader(avatarBytes))
		default:
			log.Printf("Unknown content type: %s", contentType)
			continue
		}

		cropped := CropImage(img)

		var avatarBuf bytes.Buffer
		err = png.Encode(&avatarBuf, cropped)
		if err != nil {
			log.Printf("Error encoding image %s: %v", contributor.AvatarUrl, err)
			continue
		}

		contributorsDl[index].AvatarBuf = avatarBuf.Bytes()
	}

	return contributorsDl
}

func CropImage(img image.Image) *image.RGBA {

	b := img.Bounds()
	dx := b.Dx()
	dy := b.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, dx, dy))
	// draw.Draw(dst, b, img, b.Min, draw.Src)
	draw.DrawMask(dst, dst.Bounds(), img, image.Point{}, &circle{image.Pt(dx/2, dy/2), min(dx, dy) / 2}, image.Point{}, draw.Over)

	return dst
}

func GenerateGraph(contributors []ContributorInfoDownload) string {
	// Generate SVG with template

	avatarSpace := AvatarSize + 2*AvatarMargin

	startTemplate := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="%d" height="%d">`, AvatarsPerLine*avatarSpace, (len(contributors)/AvatarsPerLine+1)*avatarSpace)
	styleTemplate := `<style></style>`

	internalContent := ""
	for index, contributor := range contributors {
		row := index % AvatarsPerLine
		line := index / AvatarsPerLine
		internalContent += fmt.Sprintf(`<a xlink:href="https://github.com/%s" class="kr-contributors" target="_blank" rel="nofollow" id="github-user-%d"><image x="%d" y="%d" width="%d" height="%d" clip-path="url(#clip)" xlink:href="data:image/png;base64,%s" /><span>@%s</span></a>`,
			contributor.Login, contributor.Id,
			row*avatarSpace+AvatarMargin, line*avatarSpace+AvatarMargin,
			AvatarSize, AvatarSize,
			base64.StdEncoding.EncodeToString(contributor.AvatarBuf), contributor.Login,
		)
	}

	endTemplate := "</svg>"

	return startTemplate + styleTemplate + internalContent + endTemplate
}
