package dependency

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"

	"github.com/sniperkit/hub/plugin/assitant/core"
)

// DepFileTypeRegistry is used to find the type of a file from its name
type DepFileTypeRegistry struct {
	nameAndSuffixMap map[string]depFileLanguageType
	prefixMap        map[string]depFileLanguageType
}

// NewDepFileTypeRegistry create a new fileTypeRegistry filled with default file types support
func NewDepFileTypeRegistry() *DepFileTypeRegistry {
	r := &DepFileTypeRegistry{
		nameAndSuffixMap: make(map[string]depFileLanguageType),
		prefixMap:        make(map[string]depFileLanguageType),
	}
	if err := r.Load(GetDefaultDepFileTypes()); err != nil {
		panic(err) // not supposed to happen
	}
	return r
}

type depFileLanguageType struct {
	Language DepLanguageID
	Type     DepFileType
}

// LoadFromJSONFile loads the types definition from a JSON file
func (r *DepFileTypeRegistry) LoadFromJSONFile(filePath string) error {
	var types []DepTypeInfo
	if err := core.ReadJSONFile(filePath, types); err != nil {
		return err
	}
	return r.Load(types)
}

// Load loads the given types definition into the registry
func (r *DepFileTypeRegistry) Load(types []DepTypeInfo) error {
	for _, value := range types {
		// process fileNames
		for _, name := range value.FileName {
			if name == "" || strings.ToLower(name) != name {
				return errors.New("invalid fileName value: must be in lower case and not empty")
			}
			err := r.findPossibleKeyConflicts(name)
			if err != nil {
				return errors.Wrap(err, "conflict was found while creating fileTypes maps")
			}
			r.nameAndSuffixMap[strings.ToLower(name)] = depFileLanguageType{value.Language, value.Type}
		}

		// process fileSuffixes
		for _, suffix := range value.FileSuffix {
			if suffix == "" || strings.ToLower(suffix) != suffix {
				return errors.New("invalid fileSuffix value: must be in lower case and not empty")
			}
			err := r.findPossibleKeyConflicts(suffix)
			if err != nil {
				return errors.Wrap(err, "conflict was found while creating fileTypes maps")
			}
			r.nameAndSuffixMap[strings.ToLower(suffix)] = depFileLanguageType{value.Language, value.Type}
		}

		// process filePrefixes
		for _, prefix := range value.FilePrefix {
			if prefix == "" || strings.ToLower(prefix) != prefix {
				return errors.New("invalid filePrefix value: must be in lower case and not empty")
			}
			err := r.findPossibleKeyConflicts(prefix)
			if err != nil {
				return errors.Wrap(err, "conflict was found while creating fileTypes maps")
			}
			r.prefixMap[strings.ToLower(prefix)] = depFileLanguageType{value.Language, value.Type}
		}
	}
	return r.findPossiblePrefixConflicts()
}

// GetDepFileTypeAndLanguage tries to identify the type and eventual language from a given file name
func (r *DepFileTypeRegistry) GetDepFileTypeAndLanguage(fileName string) (DepFileType, LanguageID) {
	fileName = strings.ToLower(fileName)

	// try first with file name
	if _, exist := r.nameAndSuffixMap[fileName]; exist {
		return DepFileType(r.nameAndSuffixMap[fileName].Type), LanguageID(r.nameAndSuffixMap[fileName].Language)
	}

	for key, value := range r.prefixMap {
		if strings.HasPrefix(fileName, key) {
			return DepFileType(value.Type), LanguageID(value.Language)
		}
	}

	for i := 0; i < len(fileName); i++ {
		if _, ok := r.nameAndSuffixMap[fileName[i:]]; ok {
			return DepFileType(r.nameAndSuffixMap[fileName[i:]].Type), LanguageID(r.nameAndSuffixMap[fileName[i:]].Language)
		}
	}

	return DepFileTypeUnknown, LanguageUnknown
}

func (r *DepFileTypeRegistry) findPossibleKeyConflicts(key string) error {
	if _, exist := r.nameAndSuffixMap[key]; exist {
		return errors.New("duplicate entry error on value " + key)
	}
	if _, exist := r.prefixMap[key]; exist {
		return errors.New("duplicate entry error on value " + key)
	}
	return nil
}

func (r *DepFileTypeRegistry) findPossiblePrefixConflicts() error {
	for key := range r.prefixMap {
		for key2 := range r.prefixMap {
			if (strings.HasPrefix(key, key2) || strings.HasPrefix(key2, key)) && key != key2 {
				return errors.Errorf("%s conflicts with %s", key, key2)
			}
		}
	}
	return nil
}
