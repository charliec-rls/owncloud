package content_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	conf "github.com/owncloud/ocis/v2/services/search/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	contentMocks "github.com/owncloud/ocis/v2/services/search/pkg/content/mocks"
)

var _ = Describe("Tika", func() {
	Describe("extract", func() {
		var (
			body         string
			fullResponse string
			language     string
			version      string
			srv          *httptest.Server
			tika         *content.Tika
		)

		BeforeEach(func() {
			body = ""
			language = ""
			version = ""
			fullResponse = ""
			srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				out := ""
				switch req.URL.Path {
				case "/version":
					out = version
				case "/language/string":
					out = language
				case "/rmeta/text":
					if fullResponse != "" {
						out = fullResponse
					} else {
						out = fmt.Sprintf(`[{"X-TIKA:content":"%s"}]`, body)
					}
				}

				_, _ = w.Write([]byte(out))
			}))

			cfg := conf.DefaultConfig()
			cfg.Extractor.Tika.TikaURL = srv.URL
			cfg.Extractor.Tika.CleanStopWords = true

			var err error
			tika, err = content.NewTikaExtractor(nil, log.NewLogger(), cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(tika).ToNot(BeNil())

			retriever := &contentMocks.Retriever{}
			retriever.On("Retrieve", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(strings.NewReader(body)), nil)

			tika.Retriever = retriever
		})

		AfterEach(func() {
			srv.Close()
		})

		It("skips non file resources", func() {
			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal(""))
		})

		It("adds content", func() {
			body = "any body"

			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{
				Type: provider.ResourceType_RESOURCE_TYPE_FILE,
				Size: 1,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal(body))
		})

		It("adds audio content", func() {
			fullResponse = `[
				{
					"xmpDM:genre": "Some Genre",
					"xmpDM:album": "Some Album",
					"xmpDM:trackNumber": "7",
					"xmpDM:discNumber": "4",
					"xmpDM:releaseDate": "2004",
					"xmpDM:artist": "Some Artist",
					"xmpDM:albumArtist": "Some AlbumArtist",
					"xmpDM:audioCompressor": "MP3",
					"xmpDM:audioChannelType": "Stereo",
					"version": "MPEG 3 Layer III Version 1",
					"xmpDM:logComment": "some comment",
					"xmpDM:audioSampleRate": "44100",
					"channels": "2",
					"dc:title": "Some Title",
					"xmpDM:duration": "225",
					"Content-Type": "audio/mpeg",
					"samplerate": "44100"
				}
			]`
			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{
				Type: provider.ResourceType_RESOURCE_TYPE_FILE,
				Size: 1,
			})
			Expect(err).ToNot(HaveOccurred())

			audio := doc.Audio
			Expect(audio).ToNot(BeNil())

			Expect(audio.Album).To(Equal(libregraph.PtrString("Some Album")))
			Expect(audio.AlbumArtist).To(Equal(libregraph.PtrString("Some AlbumArtist")))
			Expect(audio.Artist).To(Equal(libregraph.PtrString("Some Artist")))
			// Expect(audio.Bitrate).To(Equal(libregraph.PtrInt64(192)))
			// Expect(audio.Composers).To(Equal(libregraph.PtrString("Some Composers")))
			// Expect(audio.Copyright).To(Equal(libregraph.PtrString("Some Copyright")))
			Expect(audio.Disc).To(Equal(libregraph.PtrInt32(4)))
			// Expect(audio.DiscCount).To(Equal(libregraph.PtrInt32(5)))
			Expect(audio.Duration).To(Equal(libregraph.PtrInt64(225000)))
			Expect(audio.Genre).To(Equal(libregraph.PtrString("Some Genre")))
			// Expect(audio.HasDrm).To(Equal(libregraph.PtrBool(false)))
			// Expect(audio.IsVariableBitrate).To(Equal(libregraph.PtrBool(true)))
			Expect(audio.Title).To(Equal(libregraph.PtrString("Some Title")))
			Expect(audio.Track).To(Equal(libregraph.PtrInt32(7)))
			// Expect(audio.TrackCount).To(Equal(libregraph.PtrInt32(9)))
			Expect(audio.Year).To(Equal(libregraph.PtrInt32(2004)))

		})

		It("removes stop words", func() {
			body = "body to test stop words!!! against almost everyone"
			language = "en"

			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{
				Type: provider.ResourceType_RESOURCE_TYPE_FILE,
				Size: 1,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal("body test stop words!!!"))
		})

		It("keeps stop words", func() {
			body = "body to test stop words!!! against almost everyone"
			language = "en"

			tika.CleanStopWords = false
			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{
				Type: provider.ResourceType_RESOURCE_TYPE_FILE,
				Size: 1,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal(body))
		})
	})
})
