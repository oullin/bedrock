package inertia

import (
	"github.com/bedrock/packages/inertia/props"
	"github.com/bedrock/packages/inertia/response"
)

func toResponseScrollProps(in map[string]props.ScrollMeta) map[string]response.Scroll {
	if len(in) == 0 {
		return nil
	}

	out := make(map[string]response.Scroll, len(in))

	for key, meta := range in {
		out[key] = response.Scroll{
			PageName:     meta.PageName,
			PreviousPage: meta.PreviousPage,
			NextPage:     meta.NextPage,
			CurrentPage:  meta.CurrentPage,
			Reset:        meta.Reset,
		}
	}

	return out
}

func toResponseOnceProps(in map[string]props.OnceMeta) map[string]response.Once {
	if len(in) == 0 {
		return nil
	}

	out := make(map[string]response.Once, len(in))

	for key, meta := range in {
		out[key] = response.Once{
			Prop:      meta.Prop,
			ExpiresAt: meta.ExpiresAt,
		}
	}

	return out
}
