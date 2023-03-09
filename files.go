package chglog

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Parse parse a changelog.yml into ChangeLogEntries.
func Parse(file string) (entries ChangeLogEntries, err error) {
	var body []byte
	body, err = os.ReadFile(file) // nolint: gosec,gocritic
	switch {
	case os.IsNotExist(err):
		return make(ChangeLogEntries, 0), nil
	case err != nil:
		return nil, fmt.Errorf("error parsing %s: %w", file, err)
	}

	if err = yaml.Unmarshal(body, &entries); err != nil {
		return entries, fmt.Errorf("error parsing %s: %w", file, err)
	}

	return entries, nil
}

// Save save ChangeLogEntries to a yml file.
func (c *ChangeLogEntries) Save(filename string) (err error) {
	fout, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer fout.Close()

	// Work around for https://github.com/go-yaml/yaml/issues/643 (pending fix
	// https://github.com/go-yaml/yaml/pull/864).
	//
	// YAML indentation is broken when parsing a signed commit with a signed
	// annotated tag. The block is indented 4 spaces, but the second signature
	// starts with only 2. This is not a flaw in the original C libyaml, this
	// was a change specific to the go implementation.
	//
	// Changing the default indentation to 2 hides the bug. The only way to do
	// this is to switch from Marshal to using an Encoder.

	e := yaml.NewEncoder(fout)
	e.SetIndent(2)

	return e.Encode(c)
}
