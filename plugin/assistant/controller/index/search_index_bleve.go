package index

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	"github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
	"github.com/blevesearch/bleve/analysis/language/en"
)

// InitIndex initializes the search index at the specified path
func InitIndex(filepath string) (bleve.Index, error) {
	index, err := bleve.Open(filepath)

	// Doesn't yet exist (or error opening) so create a new one
	if err != nil {
		index, err = bleve.New(filepath, buildIndexMapping())
		if err != nil {
			return nil, err
		}
	}
	return index, nil
}

func buildIndexMapping() *bleve.IndexMapping {
	simpleTextFieldMapping := bleve.NewTextFieldMapping()
	simpleTextFieldMapping.Analyzer = simple_analyzer.Name

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword_analyzer.Name

	m := bleve.NewDocumentMapping()
	m.AddFieldMappingsAt("Name", simpleTextFieldMapping)
	m.AddFieldMappingsAt("FullName", simpleTextFieldMapping)
	m.AddFieldMappingsAt("Description", englishTextFieldMapping)
	m.AddFieldMappingsAt("Language", keywordFieldMapping)
	m.AddFieldMappingsAt("Tags.Name", keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Search", starMapping)
	// indexMapping.AddDocumentMapping("Starred", starMapping)

	return indexMapping
}
