package cli

import (
	"github.com/pkg/errors"
	"github.com/wesen/glazed/pkg/formatters"
	"github.com/wesen/glazed/pkg/middlewares"
	"github.com/wesen/glazed/pkg/types"
	"unicode/utf8"
)

type OutputFormatterSettings struct {
	Output          string
	TableFormat     string
	OutputAsObjects bool
	FlattenObjects  bool
	WithHeaders     bool
	CsvSeparator    string
}

func (ofs *OutputFormatterSettings) CreateOutputFormatter() (formatters.OutputFormatter, error) {
	if ofs.Output == "csv" {
		ofs.Output = "table"
		ofs.TableFormat = "csv"
	} else if ofs.Output == "tsv" {
		ofs.Output = "table"
		ofs.TableFormat = "tsv"
	} else if ofs.Output == "markdown" {
		ofs.Output = "table"
		ofs.TableFormat = "markdown"
	} else if ofs.Output == "html" {
		ofs.Output = "table"
		ofs.TableFormat = "html"
	}

	var of formatters.OutputFormatter
	if ofs.Output == "json" {
		of = formatters.NewJSONOutputFormatter(ofs.OutputAsObjects)
	} else if ofs.Output == "yaml" {
		of = formatters.NewYAMLOutputFormatter()
	} else if ofs.Output == "table" {
		if ofs.TableFormat == "csv" {
			csvOf := formatters.NewCSVOutputFormatter()
			csvOf.WithHeaders = ofs.WithHeaders
			r, _ := utf8.DecodeRuneInString(ofs.CsvSeparator)
			csvOf.Separator = r
			of = csvOf
		} else if ofs.TableFormat == "tsv" {
			tsvOf := formatters.NewTSVOutputFormatter()
			tsvOf.WithHeaders = ofs.WithHeaders
			of = tsvOf
		} else {
			of = formatters.NewTableOutputFormatter(ofs.TableFormat)
		}
		of.AddTableMiddleware(middlewares.NewFlattenObjectMiddleware())
	} else {
		return nil, errors.Errorf("Unknown output format: " + ofs.Output)
	}

	return of, nil
}

type TemplateSettings struct {
	RenameSeparator string
	UseRowTemplates bool
	Templates       map[types.FieldName]string
}

func (tf *TemplateSettings) AddMiddlewares(of formatters.OutputFormatter) error {
	if tf.UseRowTemplates && len(tf.Templates) > 0 {
		middleware, err := middlewares.NewRowGoTemplateMiddleware(tf.Templates, tf.RenameSeparator)
		if err != nil {
			return err
		}
		of.AddTableMiddleware(middleware)
	}

	return nil
}

type FieldsFilterSettings struct {
	Filters        []string
	Fields         []string
	SortColumns    bool
	ReorderColumns []string
}

func (fff *FieldsFilterSettings) AddMiddlewares(of formatters.OutputFormatter) {
	of.AddTableMiddleware(middlewares.NewFieldsFilterMiddleware(fff.Fields, fff.Filters))
	if fff.SortColumns {
		of.AddTableMiddleware(middlewares.NewSortColumnsMiddleware())
	}
	if len(fff.ReorderColumns) > 0 {
		of.AddTableMiddleware(middlewares.NewReorderColumnOrderMiddleware(fff.ReorderColumns))
	}

}

type SelectSettings struct {
	SelectField    string
	SelectTemplate string
}

func (ofs *OutputFormatterSettings) UpdateWithSelectSettings(ss *SelectSettings) {
	if ss.SelectField != "" {
		ofs.Output = "table"
		ofs.TableFormat = "tsv"
		ofs.FlattenObjects = true
		ofs.WithHeaders = false
	}
}

func (ffs *FieldsFilterSettings) UpdateWithSelectSettings(ss *SelectSettings) {
	if ss.SelectField != "" {
		ffs.Fields = []string{ss.SelectField}
	}
}

func (tf *TemplateSettings) UpdateWithSelectSettings(ss *SelectSettings) {
	if ss.SelectTemplate != "" {
		tf.UseRowTemplates = true
		tf.Templates = map[types.FieldName]string{
			ss.SelectField: ss.SelectTemplate,
		}
	}
}
