// Copyright © 2019 Martin Tournoij – This file is part of GoatCounter and
// published under the terms of a slightly modified EUPL v1.2 license, which can
// be found in the LICENSE file or at https://license.goatcounter.com

package widgets

import (
	"context"
	"html/template"

	"zgo.at/goatcounter/v2"
	"zgo.at/z18n"
)

type EntryPages struct {
	id     int
	loaded bool
	err    error
	html   template.HTML
	s      goatcounter.WidgetSettings

	Limit  int
	Detail string
	Stats  goatcounter.HitStats
}

func (w EntryPages) Name() string { return "entrypages" }
func (w EntryPages) Type() string { return "hchart" }
func (w EntryPages) Label(ctx context.Context) string {
	return z18n.T(ctx, "label/entry-pages|Entry pages")
}
func (w *EntryPages) SetHTML(h template.HTML)             { w.html = h }
func (w EntryPages) HTML() template.HTML                  { return w.html }
func (w *EntryPages) SetErr(h error)                      { w.err = h }
func (w EntryPages) Err() error                           { return w.err }
func (w EntryPages) ID() int                              { return w.id }
func (w EntryPages) Settings() goatcounter.WidgetSettings { return w.s }

func (w *EntryPages) SetSettings(s goatcounter.WidgetSettings) {
	w.s = s
	if x := s["limit"].Value; x != nil {
		w.Limit = int(x.(float64))
	}
	if x := s["key"].Value; x != nil {
		w.Detail = x.(string)
	}
}

func (w *EntryPages) GetData(ctx context.Context, a Args) (more bool, err error) {
	if w.Detail != "" {
		err = w.Stats.ListEntryPage(ctx, w.Detail, a.Rng, a.PathFilter, w.Limit, a.Offset)
	} else {
		err = w.Stats.ListEntryPages(ctx, a.Rng, a.PathFilter, w.Limit, a.Offset)
	}
	w.loaded = true
	return w.Stats.More, err
}

func (w EntryPages) RenderHTML(ctx context.Context, shared SharedData) (string, interface{}) {
	return "_dashboard_hchart.gohtml", struct {
		Context        context.Context
		ID             int
		RowsOnly       bool
		HasSubMenu     bool
		Loaded         bool
		Err            error
		IsCollected    bool
		Header         string
		TotalUniqueUTC int
		Stats          goatcounter.HitStats
		Detail         string
	}{ctx, w.id, shared.RowsOnly, true, w.loaded, w.err, true,
		z18n.T(ctx, "header/entry-pages|Entry pages"),
		shared.TotalUniqueUTC, w.Stats, w.Detail}
}
