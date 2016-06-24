package file

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInsertConfig(t *testing.T) {
	createMocks()

	doc := NewDoc("doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"source1.php", "source2.php"})

	InsertConfig(doc, sources)

	var identifier string

	err := db.QueryRow("select identifier from docs where id = ?", doc.ID).Scan(&identifier)

	assert.NoError(t, err, "Must return no errors")

	assert.Equal(t, doc.Identifier, identifier, "Must record a doc file row")

	results, err := db.Query("select path, fingerprint from sources where doc_id = ? order by path", doc.ID)

	assert.NoError(t, err, "Must return no errors")

	defer func() {
		if err := results.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	var i int

	for results.Next() {
		var path string
		var fingerprint string

		err := results.Scan(&path, &fingerprint)

		assert.NoError(t, err, "Must return no errors")

		assert.Equal(t, (*sources)[i].Path, path, "Must return source path : "+(*sources)[i].Path)
		assert.Equal(t, (*sources)[i].Fingerprint, fingerprint, "Must return fingerprint path : "+(*sources)[i].Fingerprint)

		i++
	}
}

func TestListConfigWithNoResults(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	configs := ListConfig()

	expected := []Config{}

	assert.Equal(t, &expected, configs, "Must return an empty config array")
}

func TestListConfigWithEntries(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"source1.php", "source2.php"})

	InsertConfig(doc, sources)

	doc = NewDoc("doc_file_to_track.txt", FILE)
	sources = NewSources(doc, []string{"source3.php", "source4.php", "source5.php"})

	InsertConfig(doc, sources)

	configs := ListConfig()

	assert.Len(t, *configs, 2, "Must have 2 configs")

	assert.Len(t, (*configs)[0].Sources, 2, "Must have 2 source files")
	assert.Equal(t, "source1.php", (*configs)[0].Sources[0].Path, "Must return correct path")
	assert.Equal(t, "source2.php", (*configs)[0].Sources[1].Path, "Must return correct path")

	assert.Len(t, (*configs)[1].Sources, 3, "Must have 3 source files")
	assert.Equal(t, "source3.php", (*configs)[1].Sources[0].Path, "Must return correct path")
	assert.Equal(t, "source4.php", (*configs)[1].Sources[1].Path, "Must return correct path")
	assert.Equal(t, "source5.php", (*configs)[1].Sources[2].Path, "Must return correct path")
}

func TestRemoveConfigsWithNoResults(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	configs := []Config{}

	RemoveConfigs(&configs)

	result := ListConfig()

	assert.Len(t, *result, 0, "Must have no config")
}

func TestRemoveConfigsWithOneEntry(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"source1.php", "source2.php"})

	InsertConfig(doc, sources)

	configs := ListConfig()

	RemoveConfigs(configs)

	result := ListConfig()

	assert.Len(t, *result, 0, "Must have no config remaining")
}

func TestRemoveConfigsWithSeveralEntries(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"source1.php", "source2.php"})

	InsertConfig(doc, sources)

	doc = NewDoc("doc_file_to_track.txt", FILE)
	sources = NewSources(doc, []string{"source3.php", "source4.php"})

	InsertConfig(doc, sources)

	doc = NewDoc("doc_file_to_track.txt", FILE)
	sources = NewSources(doc, []string{"source3.php", "source5.php"})

	InsertConfig(doc, sources)

	configs := ListConfig()
	expected := (*configs)[1]

	c := append((*configs)[:1], (*configs)[2:]...)

	RemoveConfigs(&c)

	result := ListConfig()

	assert.Len(t, *result, 1, "Must have no config remaining")
	assert.Equal(t, expected, (*result)[0], "Wrong configs deleted")
}