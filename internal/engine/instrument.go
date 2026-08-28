package engine

import (
	"fmt"
	"sync"

	"github.com/jb843051627/strata-forge/internal/model"
)

type InstrumentCatalog struct {
	mu    sync.RWMutex
	items map[string]*model.Instrument
}

func NewInstrumentCatalog() *InstrumentCatalog {
	return &InstrumentCatalog{items: make(map[string]*model.Instrument)}
}

func (c *InstrumentCatalog) Register(item model.Instrument) error {
	item = item.Normalize()
	if !item.Usable() {
		return fmt.Errorf("%w: unusable instrument", model.ErrInvalidInput)
	}
	c.mu.Lock()
	copy := item
	c.items[item.Kind] = &copy
	c.mu.Unlock()
	return nil
}

func (c *InstrumentCatalog) Find(kind string) (*model.Instrument, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item := c.items[kind]
	if item == nil {
		return nil, nil
	}
	copy := *item
	return &copy, nil
}

func (c *InstrumentCatalog) List() []model.Instrument {
	c.mu.RLock()
	defer c.mu.RUnlock()
	items := make([]model.Instrument, 0, len(c.items))
	for _, item := range c.items {
		if item != nil {
			items = append(items, *item)
		}
	}
	return items
}

func InstrumentLabel(item *model.Instrument) string {
	if item == nil {
		return "unassigned instrument"
	}
	return item.Kind + " / " + item.Serial
}
