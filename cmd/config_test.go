package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/util"
)

func init() {
	util.AppPath = util.AppPath + "/../" + "test"
}

func createTestDirectory() {
	err := os.RemoveAll(util.AppPath)

	if err != nil {
		logrus.Fatal(err)
	}

	err = os.Mkdir(util.AppPath, 0777)

	if err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createDocFile() {
	content := []byte("")
	err := ioutil.WriteFile(util.GetAbsPath("doc_test_1"), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createSourceFiles() {
	content := []byte("")

	for _, name := range []string{util.GetAbsPath("source_test_1"), util.GetAbsPath("/source_test_2")} {
		err := ioutil.WriteFile(name, content, 0644)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func TestParseConfigAddArgsWithMissingFileDoc(t *testing.T) {
	_, _, _, _, err := parseConfigAddArgs([]string{})

	assert.EqualError(t, err, "Missing file doc", "Must return a missing file doc error")
}

func TestParseConfigAddArgsWithUnexistingFileDoc(t *testing.T) {
	_, _, _, _, err := parseConfigAddArgs([]string{"whatever"})

	assert.EqualError(t, err, "Doc "+util.GetAbsPath("whatever")+" is not a valid existing file, nor a valid URL", "Must return an unexisting file doc error")
}

func TestParseConfigAddArgsWithMissingFileSources(t *testing.T) {
	createTestDirectory()
	createDocFile()

	_, _, _, _, err := parseConfigAddArgs([]string{"doc_test_1"})

	assert.EqualError(t, err, "Missing file/folder sources", "Must return a missing file sources error")
}

func TestParseConfigAddArgsWithUnexistingFileSources(t *testing.T) {
	_, _, _, _, err := parseConfigAddArgs([]string{"doc_test_1", "whatever"})

	assert.EqualError(t, err, "File/folder source whatever doesn't exist", "Must return a unexisting source file error")
}

func TestParseConfigAddArgsWithFile(t *testing.T) {
	createTestDirectory()
	createDocFile()
	createSourceFiles()

	doc, docCat, folderSources, fileSources, err := parseConfigAddArgs([]string{"doc_test_1", "source_test_1,source_test_2"})

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "doc_test_1", doc, "Must return doc file path")
	assert.True(t, 0 == docCat, "Must return a file doc category")
	assert.Equal(t, []string{"source_test_1", "source_test_2"}, fileSources, "Must return sources file path")
	assert.Equal(t, []string{}, folderSources, "Must return empty folder sources")
}

func TestParseConfigAddArgsWithURL(t *testing.T) {
	createTestDirectory()
	createDocFile()
	createSourceFiles()

	doc, docCat, folderSources, fileSources, err := parseConfigAddArgs([]string{"http://google.com", "source_test_1,source_test_2"})

	assert.Equal(t, "http://google.com", doc, "Must return a doc url")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, []string{"source_test_1", "source_test_2"}, fileSources, "Must return sources file path")
	assert.True(t, 1 == docCat, "Must return an URL doc category")
	assert.Equal(t, []string{}, folderSources, "Must return empty folder sources")
}

func TestParseConfigDelArgsWithArgumentNotANumber(t *testing.T) {
	configs := []file.Config{
		file.Config{},
		file.Config{},
		file.Config{},
		file.Config{},
	}

	_, err := parseConfigDelArgs(&configs, "1,2,3,a")

	assert.EqualError(t, err, "a is not a number", "Must return an error")
}

func TestParseConfigDelArgsWithArgumentNotInRange(t *testing.T) {
	configs := []file.Config{
		file.Config{},
		file.Config{},
		file.Config{},
	}

	_, err := parseConfigDelArgs(&configs, "3,4")

	assert.EqualError(t, err, "Value 3 is out of bounds", "Must return an error")
}

func TestParseConfigDelArgs(t *testing.T) {
	configs := []file.Config{
		file.Config{Doc: file.Doc{Identifier: "source_0.php"}},
		file.Config{Doc: file.Doc{Identifier: "source_1.php"}},
		file.Config{Doc: file.Doc{Identifier: "source_2.php"}},
		file.Config{Doc: file.Doc{Identifier: "source_3.php"}},
		file.Config{Doc: file.Doc{Identifier: "source_4.php"}},
	}

	expected := &[]file.Config{
		file.Config{Doc: file.Doc{Identifier: "source_3.php"}},
		file.Config{Doc: file.Doc{Identifier: "source_4.php"}},
	}

	results, err := parseConfigDelArgs(&configs, "3,4")

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, expected, results, "Must return configs")
}
