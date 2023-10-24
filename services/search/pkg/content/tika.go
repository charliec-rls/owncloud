package content

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/google/go-tika/tika"

	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
)

// Tika is used to extract content from a resource,
// it uses apache tika to retrieve all the data.
type Tika struct {
	*Basic
	Retriever
	tika                       *tika.Client
	ContentExtractionSizeLimit uint64
	CleanStopWords             bool
}

// NewTikaExtractor creates a new Tika instance.
func NewTikaExtractor(gatewaySelector pool.Selectable[gateway.GatewayAPIClient], logger log.Logger, cfg *config.Config) (*Tika, error) {
	basic, err := NewBasicExtractor(logger)
	if err != nil {
		return nil, err
	}

	tk := tika.NewClient(nil, cfg.Extractor.Tika.TikaURL)
	tkv, err := tk.Version(context.Background())
	if err != nil {
		return nil, err
	}
	logger.Info().Msgf("Tika version: %s", tkv)

	return &Tika{
		Basic:                      basic,
		Retriever:                  newCS3Retriever(gatewaySelector, logger, cfg.Extractor.CS3AllowInsecure),
		tika:                       tika.NewClient(nil, cfg.Extractor.Tika.TikaURL),
		ContentExtractionSizeLimit: cfg.ContentExtractionSizeLimit,
		CleanStopWords:             cfg.Extractor.Tika.CleanStopWords,
	}, nil
}

// Extract loads a resource from its underlying storage, passes it to tika and processes the result into a Document.
func (t Tika) Extract(ctx context.Context, ri *provider.ResourceInfo) (Document, error) {
	doc, err := t.Basic.Extract(ctx, ri)
	if err != nil {
		return doc, err
	}

	if ri.Size == 0 {
		return doc, nil
	}

	if ri.Size > t.ContentExtractionSizeLimit {
		t.logger.Info().Interface("ResourceID", ri.Id).Str("Name", ri.Name).Msg("file exceeds content extraction size limit. skipping.")
		return doc, nil
	}

	if ri.Type != provider.ResourceType_RESOURCE_TYPE_FILE {
		return doc, nil
	}

	data, err := t.Retrieve(ctx, ri.Id)
	if err != nil {
		return doc, err
	}
	defer data.Close()

	metas, err := t.tika.MetaRecursive(ctx, data)
	if err != nil {
		return doc, err
	}

	for _, meta := range metas {
		if title, err := getFirstValue(meta, "title"); err == nil {
			doc.Title = strings.TrimSpace(fmt.Sprintf("%s %s", doc.Title, title))
		}

		if content, err := getFirstValue(meta, "X-TIKA:content"); err == nil {
			doc.Content = strings.TrimSpace(fmt.Sprintf("%s %s", doc.Content, content))
		}

		if contentType, err := getFirstValue(meta, "Content-Type"); err == nil && strings.HasPrefix(contentType, "audio/") {
			audio := libregraph.Audio{}

			if v, err := getFirstValue(meta, "xmpDM:album"); err == nil {
				audio.SetAlbum(v)
			}

			if v, err := getFirstValue(meta, "xmpDM:albumArtist"); err == nil {
				audio.SetAlbumArtist(v)
			}

			if v, err := getFirstValue(meta, "xmpDM:artist"); err == nil {
				audio.SetArtist(v)
			}

			// TODO: audio.Bitrate: not provided by tika
			// TODO: audio.Composers: not provided by tika
			// TODO: audio.Copyright: not provided by tika for audio files?

			if v, err := getFirstValue(meta, "xmpDM:discNumber"); err == nil {
				if i, err := strconv.ParseInt(v, 10, 32); err == nil {
					audio.SetDisc(int32(i))
				}

			}

			//  TODO: audio.DiscCount: not provided by tika

			if v, err := getFirstValue(meta, "xmpDM:duration"); err == nil {
				if i, err := strconv.ParseInt(v, 10, 64); err == nil {
					audio.SetDuration(i * 1000)
				}
			}

			if v, err := getFirstValue(meta, "xmpDM:genre"); err == nil {
				audio.SetGenre(v)
			}

			// TODO: audio.HasDrm: not provided by tika
			// TODO: audio.IsVariableBitrate: not provided by tika

			if v, err := getFirstValue(meta, "dc:title"); err == nil {
				audio.SetTitle(v)
			}

			if v, err := getFirstValue(meta, "xmpDM:trackNumber"); err == nil {
				if i, err := strconv.ParseInt(v, 10, 32); err == nil {
					audio.SetTrack(int32(i))
				}
			}

			// TODO: audio.TrackCount: not provided by tika

			if v, err := getFirstValue(meta, "xmpDM:releaseDate"); err == nil {
				if i, err := strconv.ParseInt(v, 10, 32); err == nil {
					audio.SetYear(int32(i))
				}
			}

			doc.Audio = &audio
		}
	}

	if langCode, _ := t.tika.LanguageString(ctx, doc.Content); langCode != "" && t.CleanStopWords {
		doc.Content = CleanString(doc.Content, langCode)
	}

	return doc, nil
}
