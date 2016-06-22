package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSourcesFingerprint(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("/tmp/doc-hunt/doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"/tmp/doc-hunt/source1.php", "/tmp/doc-hunt/source2.php"})

	InsertConfig(doc, sources)

	configsBefore := ListConfig()

	err := ioutil.WriteFile(fmt.Sprint("/tmp/doc-hunt/source1.php"), []byte("<?php echo 'Hello world !';"), 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	UpdateSourcesFingeprint()
	configsAfter := ListConfig()

	assert.Equal(t, (*configsAfter)[0].Sources[0].ID, (*configsBefore)[0].Sources[0].ID, "Must return same id")
	assert.NotEqual(t, (*configsAfter)[0].Sources[0].Fingerprint, (*configsBefore)[0].Sources[0].Fingerprint, "Must return a different fingerprint")
	assert.True(t, (*configsAfter)[0].Sources[0].UpdatedAt.After((*configsAfter)[0].Sources[0].CreatedAt), "Must change updated_at date")

	assert.Equal(t, (*configsAfter)[0].Sources[1].ID, (*configsBefore)[0].Sources[1].ID, "Must return same id")
	assert.Equal(t, (*configsAfter)[0].Sources[1].Fingerprint, (*configsBefore)[0].Sources[1].Fingerprint, "Must return same fingerprint")
	assert.True(t, (*configsAfter)[0].Sources[0].UpdatedAt.After((*configsAfter)[0].Sources[0].CreatedAt), "Must change updated_at date")
}

func TestDeleteSources(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("/tmp/doc-hunt/doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"/tmp/doc-hunt/source1.php", "/tmp/doc-hunt/source2.php"})

	InsertConfig(doc, sources)

	configsBefore := ListConfig()

	err := os.Remove("/tmp/doc-hunt/source1.php")

	if err != nil {
		logrus.Fatal(err)
	}

	DeleteSources(&[]string{"/tmp/doc-hunt/source1.php"})
	configsAfter := ListConfig()

	assert.Len(t, (*configsBefore)[0].Sources, 2, "Must contains 2 elements")

	assert.Len(t, (*configsAfter)[0].Sources, 1, "Must contains 1 element")
	assert.Equal(t, "/tmp/doc-hunt/source2.php", (*configsAfter)[0].Sources[0].Path, "Must keep element not deleted")
}

func TestDeleteSourcesWithOnlyOneSourceRemaining(t *testing.T) {
	createMocks()
	createDocFiles()
	deleteDatabase()
	createTables()

	doc := NewDoc("/tmp/doc-hunt/doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"/tmp/doc-hunt/source1.php"})

	InsertConfig(doc, sources)

	doc = NewDoc("/tmp/doc-hunt/doc_file_to_track_2.txt", FILE)
	sources = NewSources(doc, []string{"/tmp/doc-hunt/source2.php"})

	InsertConfig(doc, sources)

	configsBefore := ListConfig()

	err := os.Remove("/tmp/doc-hunt/source1.php")

	if err != nil {
		logrus.Fatal(err)
	}

	DeleteSources(&[]string{"/tmp/doc-hunt/source1.php"})
	configsAfter := ListConfig()

	rows, err := db.Query("select d.id from docs d where id = ?", (*configsBefore)[0].Doc.ID)

	if err != nil {
		logrus.Fatal(err)
	}

	assert.Len(t, (*configsBefore)[0].Sources, 1, "Must contains 1 element")
	assert.Len(t, (*configsAfter), 1, "Must contains remaining config only")
	assert.False(t, rows.Next(), "Must have deleted doc id")
}

func TestUpdateFilenameSources(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("/tmp/doc-hunt/doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"/tmp/doc-hunt/source1.php"})

	InsertConfig(doc, sources)

	configsBefore := ListConfig()

	UpdateFilenameSources(&map[string]string{"/tmp/doc-hunt/source1.php": "/tmp/doc-hunt/source2.php"})

	configsAfter := ListConfig()

	assert.Equal(t, (*configsBefore)[0].Sources[0].Path, "/tmp/doc-hunt/source1.php", "Must contains original path")
	assert.Equal(t, (*configsAfter)[0].Sources[0].Path, "/tmp/doc-hunt/source2.php", "Must contains renamed path")
}
