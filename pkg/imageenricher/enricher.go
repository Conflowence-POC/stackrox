package imageenricher

import (
	"time"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/logging"
	"github.com/karlseguin/ccache"
	"golang.org/x/time/rate"
)

var (
	logger = logging.LoggerForModule()
)

// ImageEnricher provides functions for enriching images with integrations.
type ImageEnricher interface {
	EnrichImage(image *v1.Image) bool

	IntegrationSet() IntegrationSet
}

// New returns a new ImageEnricher instance. You should use the singleton in singleton.go instead.
func New() ImageEnricher {
	return &enricherImpl{
		integrations: NewIntegrationSet(),

		metadataLimiter: rate.NewLimiter(rate.Every(5*time.Second), 3),
		metadataCache:   ccache.New(ccache.Configure().MaxSize(maxCacheSize).ItemsToPrune(itemsToPrune)),
		scanLimiter:     rate.NewLimiter(rate.Every(5*time.Second), 3),
		scanCache:       ccache.New(ccache.Configure().MaxSize(maxCacheSize).ItemsToPrune(itemsToPrune)),
	}
}
