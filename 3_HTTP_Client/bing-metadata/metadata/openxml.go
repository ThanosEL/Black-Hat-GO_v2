package metadata

import (
	"archive/zip"
	"encoding/xml"
	"strings"
)

// OfficeCoreProperty maps to docProps/core.xml inside an Open XML file.
// Contains author/editor info — useful for harvesting real employee usernames.
type OfficeCoreProperty struct {
	XMLName        xml.Name `xml:"coreProperties"`
	Creator        string   `xml:"creator"`
	LastModifiedBy string   `xml:"lastModifiedBy"`
}

// OfficeAppProperty maps to docProps/app.xml inside an Open XML file.
// Contains which application/version created the document — useful for
// identifying what software (and version) the target organization uses.
type OfficeAppProperty struct {
	XMLName     xml.Name `xml:"Properties"`
	Application string   `xml:"Application"`
	Company     string   `xml:"Company"`
	Version     string   `xml:"AppVersion"`
}

// OfficeVersions maps the raw major version number found in AppVersion
// (e.g. "16") to the human-readable release year (e.g. "2016").
var OfficeVersions = map[string]string{
	"16": "2016",
	"15": "2013",
	"14": "2010",
	"12": "2007",
	"11": "2003",
}

// GetMajorVersion extracts the major version number from a raw AppVersion
// string like "15.0300" and looks up the human-readable release year.
func (a *OfficeAppProperty) GetMajorVersion() string {
	tokens := strings.Split(a.Version, ".")
	if len(tokens) < 2 {
		return "Unknown"
	}
	v, ok := OfficeVersions[tokens[0]]
	if !ok {
		return "Unknown"
	}
	return v
}

// NewProperties opens a zip.Reader (i.e. an Open XML file treated as a ZIP
// archive), locates the two metadata XML files inside it, and parses them
// into OfficeCoreProperty and OfficeAppProperty structs.
func NewProperties(r *zip.Reader) (*OfficeCoreProperty, *OfficeAppProperty, error) {
	var coreProps OfficeCoreProperty
	var appProps OfficeAppProperty

	// Loop through every file inside the ZIP archive.
	for _, f := range r.File {
		switch f.Name {
		case "docProps/core.xml":
			if err := process(f, &coreProps); err != nil {
				return nil, nil, err
			}
		case "docProps/app.xml":
			if err := process(f, &appProps); err != nil {
				return nil, nil, err
			}
		default:
			continue // ignore every other file in the archive
		}
	}
	return &coreProps, &appProps, nil
}

// process opens a single file from within the ZIP archive and decodes its
// XML content into whatever struct "prop" points to. Using interface{}
// here means this same function works for BOTH OfficeCoreProperty and
// OfficeAppProperty — no need to duplicate this logic per type.
func process(f *zip.File, prop interface{}) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	if err := xml.NewDecoder(rc).Decode(&prop); err != nil {
		return err
	}
	return nil
}
