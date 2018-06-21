package dependency

// DepType identifies the type of a dependency file in the repository
type DepType string

func (id *DepLanguageID) String() string {
	return string(*id)
}

// DepLanguageID identifies the language used for a dependency file
type DepLanguageID string

func (id *DepType) String() string {
	return string(*id)
}

// DepFileTypeInfo is the struct used to extract DepType Information from a JSON
type DepTypeInfo struct {
	Language   DepLanguageID `json:"language"`
	Type       DepType       `json:"type"`
	FileName   []string      `json:"fileNames"`
	FilePrefix []string      `json:"filePrefixes"`
	FileSuffix []string      `json:"fileSuffixes"`
}

const (
	// DepFileTypeSource tags a file as being source code in a given language
	DepFileTypeSource DepType = "Source"
	// DepFileTypeGenerator tags a file as being part of a source file generator tool
	DepFileTypeGenerator DepType = "Generator"
	// DepFileTypeBuildSystem tags a file as being part of a build system tool
	DepFileTypeBuildSystem DepType = "BuildSystem"
	// DepFileTypeEnvConfig tags a file as being used to configure some development tools
	DepFileTypeEnvConfig DepType = "EnvConfig"
	// DepFileTypeLicense tags is used to tag license files
	DepFileTypeLicense DepType = "License"
	// DepFileTypeUnknown tags unrecognized files
	DepFileTypeUnknown DepType = ""
	// LanguageUnknown is used to identify unknown languages
	LanguageUnknown LanguageID = ""
)
