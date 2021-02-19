// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	"bytes"
	bleve "github.com/blevesearch/bleve"
	mappings "github.com/stackrox/rox/central/imagecomponent/mappings"
	metrics "github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
	"time"
)

const batchSize = 5000

const resourceName = "ImageComponent"

type indexerImpl struct {
	index bleve.Index
}

type imageComponentWrapper struct {
	*storage.ImageComponent `json:"image_component"`
	Type                    string `json:"type"`
}

func (b *indexerImpl) AddImageComponent(imagecomponent *storage.ImageComponent) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "ImageComponent")
	if err := b.index.Index(imagecomponent.GetId(), &imageComponentWrapper{
		ImageComponent: imagecomponent,
		Type:           v1.SearchCategory_IMAGE_COMPONENTS.String(),
	}); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) AddImageComponents(imagecomponents []*storage.ImageComponent) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "ImageComponent")
	batchManager := batcher.New(len(imagecomponents), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(imagecomponents[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(imagecomponents []*storage.ImageComponent) error {
	batch := b.index.NewBatch()
	for _, imagecomponent := range imagecomponents {
		if err := batch.Index(imagecomponent.GetId(), &imageComponentWrapper{
			ImageComponent: imagecomponent,
			Type:           v1.SearchCategory_IMAGE_COMPONENTS.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) Count(q *v1.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "ImageComponent")
	return blevesearch.RunCountRequest(v1.SearchCategory_IMAGE_COMPONENTS, q, b.index, mappings.OptionsMap, opts...)
}

func (b *indexerImpl) DeleteImageComponent(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "ImageComponent")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) DeleteImageComponents(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "ImageComponent")
	batch := b.index.NewBatch()
	for _, id := range ids {
		batch.Delete(id)
	}
	if err := b.index.Batch(batch); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) MarkInitialIndexingComplete() error {
	return b.index.SetInternal([]byte(resourceName), []byte("old"))
}

func (b *indexerImpl) NeedsInitialIndexing() (bool, error) {
	data, err := b.index.GetInternal([]byte(resourceName))
	if err != nil {
		return false, err
	}
	return !bytes.Equal([]byte("old"), data), nil
}

func (b *indexerImpl) Search(q *v1.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "ImageComponent")
	return blevesearch.RunSearchRequest(v1.SearchCategory_IMAGE_COMPONENTS, q, b.index, mappings.OptionsMap, opts...)
}
