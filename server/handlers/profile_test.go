package handlers

import (
	"context"
	"fmt"
	"ftbadge/server/cache"
	"ftbadge/server/ftapi"
	"ftbadge/server/utils"
	"image"
	"image/color"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

type cacheMock struct{}

func (c *cacheMock) Get(ctx context.Context, key string) (string, bool, error) {
	return "", false, nil
}

func (c *cacheMock) BulkSet(ctx context.Context, entries []cache.CacheEntry) error {
	return nil
}

func (c *cacheMock) BulkGet(ctx context.Context, keys ...string) ([]*string, error) {
	return make([]*string, len(keys)), nil
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := `{
		"access_token": "test_access_token",
		"expires_in": 7200
	}`
	w.Write([]byte(response))
}

func getUserHandler(cdnURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := fmt.Sprintf(`{
			"email": "testuser@student.42angouleme.fr",
			"displayname": "testuser",
			"kind": "student",
			"Image": {
				"Versions": {
					"medium": "%s/avatar/testuser"
				}
			},
			"cursus_users": [
				{
					"grade": "Transcender",
					"level": 42.0,
					"cursus": {
						"name": "42cursus"
					}
				}
			]
		}`, cdnURL)
		w.Write([]byte(response))
	}
}

func getAvatarHandler(img []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)

		w.Write(img)
	}
}

func randomImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, g, b := rand.IntN(256), rand.IntN(256), rand.IntN(256)
			c := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			img.Set(x, y, c)
		}
	}

	bytes, err := utils.EncodeToJPEG(img, 70)
	if err != nil {
		panic(fmt.Sprintf("Failed to encode random image: %v", err))
	}

	return bytes
}

func BenchmarkRenderProfile(b *testing.B) {
	cc := &cacheMock{}

	img := randomImage()

	cdnMux := http.NewServeMux()
	cdnMux.HandleFunc("/avatar/testuser", getAvatarHandler(img))
	cdnServer := httptest.NewServer(cdnMux)
	defer cdnServer.Close()

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/oauth/token", oauthHandler)
	apiMux.HandleFunc("/users/testuser", getUserHandler(cdnServer.URL))
	apiServer := httptest.NewServer(apiMux)
	defer apiServer.Close()

	ftc := ftapi.NewClient(apiServer.URL, cdnServer.URL)

	for b.Loop() {
		if _, err := renderProfile(b.Context(), ftc, cc, "testuser"); err != nil {
			b.Fatalf("Failed to render profile: %v", err)
		}
	}
}
