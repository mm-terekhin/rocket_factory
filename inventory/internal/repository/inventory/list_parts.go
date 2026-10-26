package inventory

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/repository/converter"
)

type preparedFilter struct {
	uuids      map[string]struct{}
	names      map[string]struct{}
	categories map[model.Category]struct{}
	countries  map[string]struct{}
	tags       map[string]struct{}
}

func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pf := prepareFilter(&filter)

	out := make([]model.Part, 0, len(r.parts))
	for _, part := range r.parts {
		partModel := converter.PartToModel(part)

		if matches(&partModel, pf) {
			out = append(out, partModel)
		}
	}

	return out, nil
}

func prepareFilter(f *model.PartsFilter) preparedFilter {
	pf := preparedFilter{}

	if len(f.UUIDs) > 0 {
		pf.uuids = toStringSet(f.UUIDs)
	}

	if len(f.Names) > 0 {
		pf.names = toStringSet(f.Names)
	}

	if len(f.Categories) > 0 {
		pf.categories = toCategorySet(f.Categories)
	}

	if len(f.ManufacturerCountries) > 0 {
		pf.countries = toStringSet(f.ManufacturerCountries)
	}

	if len(f.Tags) > 0 {
		pf.tags = toStringSet(f.Tags)
	}

	return pf
}

func toStringSet(items []string) map[string]struct{} {
	m := make(map[string]struct{}, len(items))

	for _, s := range items {
		m[s] = struct{}{}
	}

	return m
}

func toCategorySet(items []model.Category) map[model.Category]struct{} {
	m := make(map[model.Category]struct{}, len(items))

	for _, s := range items {
		m[s] = struct{}{}
	}

	return m
}

func matches(part *model.Part, pf preparedFilter) bool {
	if len(pf.uuids) > 0 {
		if _, ok := pf.uuids[part.UUID]; !ok {
			return false
		}
	}

	if len(pf.names) > 0 {
		if _, ok := pf.names[part.Info.Name]; !ok {
			return false
		}
	}

	if len(pf.categories) > 0 {
		if _, ok := pf.categories[part.Info.Category]; !ok {
			return false
		}
	}

	if len(pf.countries) > 0 {
		if _, ok := pf.countries[part.Info.Manufacturer.Country]; !ok {
			return false
		}
	}

	if len(pf.tags) > 0 {
		matched := false
		for _, t := range part.Info.Tags {
			if _, ok := pf.tags[t]; ok {
				matched = true
				break
			}
		}

		if !matched {
			return false
		}
	}

	return true
}
