package query

import (
	"errors"
	"fmt"
	"time"

	"github.com/Thibault-Van-Win/thehive4go/utils/tlp"
)

// Struct to hold query params
type Options struct {
	listing   string
	startDate time.Time
	company   string
	sortField string
	sortOrder string
	maxTLP    tlp.TLP
}

// Function type used to set options
type Option func(*Options)

// Create new options with default config
func initOptions() Options {
	return Options{
		sortField: "severity",
		sortOrder: "desc",
		maxTLP:    tlp.AMBER,
	}
}

// Private option getting listings
func WithListing(listing string) Option {
	return func(fo *Options) {
		fo.listing = listing
	}
}

// Option to set the start date
func WithStartDate(startDate time.Time) Option {
	return func(fo *Options) {
		fo.startDate = startDate
	}
}

// Option to fetch artifacts of a certain company
func WithCompany(company string) Option {
	return func(fo *Options) {
		fo.company = company
	}
}

// Option to sort the result
func WithSort(field, order string) Option {
	return func(fo *Options) {
		fo.sortField = field
		fo.sortOrder = order
	}
}

// Option for setting a max TLP
func WithMaxTLP(tlp tlp.TLP) Option {
	return func(fo *Options) {
		fo.maxTLP = tlp
	}
}

// Parse the options to build a query
func (opts *Options) buildQuery() (map[string]any, error) {
	// Validate the options
	if err := opts.validate(); err != nil {
		return nil, fmt.Errorf("validation error: %v", err)
	}

	var startTimestamp int64
	if !opts.startDate.IsZero() {
		startTimestamp = opts.startDate.UnixNano() / int64(time.Millisecond)
	}

	// This could be translated to a private query option
	query := map[string]any{
		"query": []any{
			map[string]any{
				"_name": opts.listing,
			},
		},
	}

	// Append filters dynamically based on the options that are set
	if startTimestamp > 0 {
		query["query"] = append(query["query"].([]any), map[string]any{
			"_name": "filter",
			"_gt": map[string]any{
				"_field": "newDate",
				"_value": startTimestamp,
			},
		})
	}

	if opts.company != "" {
		query["query"] = append(query["query"].([]any), map[string]any{
			"_name": "filter",
			"_eq": map[string]any{
				"_field": "customFields.company.string",
				"_value": opts.company,
			},
		})
	}

	if opts.maxTLP >= 0 {
		query["query"] = append(query["query"].([]any), map[string]any{
			"_name": "filter",
			"_lte": map[string]any{
				"_field": "tlp",
				"_value": opts.maxTLP,
			},
		})
	}

	query["query"] = append(query["query"].([]any), map[string]any{
		"_name": "sort",
		"_fields": []map[string]any{
			{
				opts.sortField: opts.sortOrder,
			},
		},
	})

	return query, nil
}

func (opts *Options) validate() error {
	if opts.listing == "" {
		return errors.New("listing is empty")
	}

	if err := tlp.Validate(opts.maxTLP); err != nil {
		return fmt.Errorf("invalid TLP value: %v", err)
	}

	if !(opts.sortOrder == "asc" || opts.sortOrder == "desc") {
		return fmt.Errorf("invalid sorting order: %s, options are: ascending (asc) or descending (desc)", opts.sortOrder)
	}

	return nil
}
