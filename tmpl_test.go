package tmpl_test

import (
	"fmt"
	"github.com/kanamone/tmpl"
	"strings"
	"testing"
)

func TestTemplate(t *testing.T) {
	tmp := tmpl.NewTemplate("/*<slot name>*/content/*</slot>*/")

	if tmp == nil {
		t.Errorf("Expected template, got nil")
		return
	}
}

func TestGetSlotMatchesSingle(t *testing.T) {
	tmp := tmpl.NewTemplate("/*<slot name>*/content/*</slot>*/")

	if tmp == nil {
		t.Errorf("Expected template, got nil")
		return
	}

	matches := tmp.GetSlotMatches()

	if len(matches) == 0 {
		t.Errorf("Expected matches, got empty")
		return
	}

	if matches[0].StartIndex != 0 {
		t.Errorf("Expected startByte 0, got %d", matches[0].StartIndex)
	}

	if matches[0].EndIndex != 33 {
		t.Errorf("Expected endByte 8, got %d", matches[0].EndIndex)
	}

	if matches[0].Identifier != "name" {
		t.Errorf("Expected identifier name, got %s", matches[0].Identifier)
	}

	if matches[0].InnerContent != "content" {
		t.Errorf("Expected innerContent content, got %s", matches[0].InnerContent)
	}

	if matches[0].OuterContent != "/*<slot name>*/content/*</slot>*/" {
		t.Errorf("Expected outerContent /*<slot name>*/content/*</slot>*/, got %s", matches[0].OuterContent)
	}
}

// TestGetSlotMatchesSingleMultiple tests that the GetSlotMatches function
// e.g. /*<slot a>*/content1/*</slot>*/ /*<slot b>*/content2/*</slot>*/
func TestGetSlotMatchesMultiple(t *testing.T) {
	tmp := tmpl.NewTemplate("/*<slot a>*/content1/*</slot>*/ /*<slot b>*/content2/*</slot>*/")

	if tmp == nil {
		t.Errorf("Expected template, got nil")
		return
	}

	matches := tmp.GetSlotMatches()

	if len(matches) == 0 {
		t.Errorf("Expected matches, got empty")
		return
	}

	if matches[0].StartIndex != 0 {
		t.Errorf("Expected startByte 0, got %d", matches[0].StartIndex)
	}

	if matches[0].EndIndex != 31 {
		t.Errorf("Expected endByte 31, got %d", matches[0].EndIndex)
	}

	if matches[0].Identifier != "a" {
		t.Errorf("Expected identifier a, got %s", matches[0].Identifier)
	}

	if matches[0].InnerContent != "content1" {
		t.Errorf("Expected innerContent content1, got %s", matches[0].InnerContent)
	}

	if matches[0].OuterContent != "/*<slot a>*/content1/*</slot>*/" {
		t.Errorf("Expected outerContent /*<slot a>*/content1/*</slot>*/, got %s", matches[0].OuterContent)
	}

	if matches[1].StartIndex != 32 {
		t.Errorf("Expected startByte 32, got %d", matches[1].StartIndex)
	}

	if matches[1].EndIndex != 63 {
		t.Errorf("Expected endByte 63, got %d", matches[1].EndIndex)
	}

	if matches[1].Identifier != "b" {
		t.Errorf("Expected identifier b, got %s", matches[1].Identifier)
	}

	if matches[1].InnerContent != "content2" {
		t.Errorf("Expected innerContent content2, got %s", matches[1].InnerContent)
	}

	if matches[1].OuterContent != "/*<slot b>*/content2/*</slot>*/" {
		t.Errorf("Expected outerContent /*<slot b>*/content2/*</slot>*/, got %s", matches[1].OuterContent)
	}
}

func TestGetSlotMatchesMultilinePreserveTags(t *testing.T) {
	tmp := tmpl.NewTemplate("/*<slot name>*/\ncontent\n/*</slot>*/")

	if tmp == nil {
		t.Errorf("Expected template, got nil")
		return
	}

	matches := tmp.GetSlotMatches()

	if len(matches) == 0 {
		t.Errorf("Expected matches, got empty")
		return
	}

	if matches[0].StartIndex != 0 {
		t.Errorf("Expected startIndex 0, got %d", matches[0].StartIndex)
	}

	if matches[0].EndIndex != 35 {
		t.Errorf("Expected endIndex 8, got %d", matches[0].EndIndex)
	}

	if matches[0].Identifier != "name" {
		t.Errorf("Expected identifier name, got %s", matches[0].Identifier)
	}

	if matches[0].InnerContent != "\ncontent\n" {
		t.Errorf("Expected innerContent content, got %s", matches[0].InnerContent)
	}

	if matches[0].OuterContent != "/*<slot name>*/\ncontent\n/*</slot>*/" {
		t.Errorf("Expected outerContent /*<slot name>*/content/*</slot>*/, got %s", matches[0].OuterContent)
	}
}

// TestGetSlotMatchesMultilineOuter tests that the GetSlotMatches function
// correctly returns the start and end indices of the slot delimiters when the
// inner content spans multiple lines. e.g. "Big\n/*<slot name>*/Band/*</slot>*/\nBeat"
func TestGetSlotMatchesMultilineOuter(t *testing.T) {
	tmp := tmpl.NewTemplate("Big\n/*<slot name>*/Band/*</slot>*/\nBeat")

	if tmp == nil {
		t.Errorf("Expected template, got nil")
		return
	}

	matches := tmp.GetSlotMatches()

	if len(matches) == 0 {
		t.Errorf("Expected matches, got empty")
		return
	}

	if matches[0].StartIndex != 4 {
		t.Errorf("Expected startIndex 4, got %d", matches[0].StartIndex)
	}

	if matches[0].EndIndex != 34 {
		t.Errorf("Expected endIndex 34, got %d", matches[0].EndIndex)
	}

	if matches[0].Identifier != "name" {
		t.Errorf("Expected identifier name, got %s", matches[0].Identifier)
	}

	if matches[0].InnerContent != "Band" {
		t.Errorf("Expected innerContent Band, got %s", matches[0].InnerContent)
	}

	if matches[0].OuterContent != "/*<slot name>*/Band/*</slot>*/" {
		t.Errorf("Expected outerContent /*<slot name>*/Band/*</slot>*/, got %s", matches[0].OuterContent)
	}
}

func TestGetSlotMatchesUnicode(t *testing.T) {
	tmp := tmpl.NewTemplate("/*<slot name>*/🐱🐶🐭🐹🐰/*</slot>*/")

	if tmp == nil {
		t.Errorf("Expected template, got nil")
		return
	}

	matches := tmp.GetSlotMatches()

	if len(matches) == 0 {
		t.Errorf("Expected matches, got empty")
		return
	}

	if matches[0].Identifier != "name" {
		t.Errorf("Expected identifier name, got %s", matches[0].Identifier)
	}

	if matches[0].InnerContent != "🐱🐶🐭🐹🐰" {
		t.Errorf("Expected innerContent content, got %s", matches[0].InnerContent)
	}

	if matches[0].OuterContent != "/*<slot name>*/🐱🐶🐭🐹🐰/*</slot>*/" {
		t.Errorf("Expected outerContent /*<slot name>*/content/*</slot>*/, got %s", matches[0].OuterContent)
	}

	if matches[0].StartIndex != 0 {
		t.Errorf("Expected startIndex 0, got %d", matches[0].StartIndex)
	}

	if matches[0].EndIndex != 46 {
		t.Errorf("Expected endIndex 46, got %d", matches[0].EndIndex)
	}
}

func TestGetSlotMatchesJapanese(t *testing.T) {
	tmp := tmpl.NewTemplate("/*<slot name>*/こんにちは/*</slot>*/")

	if tmp == nil {
		t.Errorf("Expected template, got nil")
		return
	}

	matches := tmp.GetSlotMatches()

	if len(matches) == 0 {
		t.Errorf("Expected matches, got empty")
		return
	}

	if matches[0].Identifier != "name" {
		t.Errorf("Expected identifier name, got %s", matches[0].Identifier)
	}

	if matches[0].InnerContent != "こんにちは" {
		t.Errorf("Expected innerContent content, got %s", matches[0].InnerContent)
	}

	if matches[0].OuterContent != "/*<slot name>*/こんにちは/*</slot>*/" {
		t.Errorf("Expected outerContent /*<slot name>*/content/*</slot>*/, got %s", matches[0].OuterContent)
	}

	if matches[0].StartIndex != 0 {
		t.Errorf("Expected startIndex 0, got %d", matches[0].StartIndex)
	}

	if matches[0].EndIndex != 41 {
		t.Errorf("Expected endIndex 41, got %d", matches[0].EndIndex)
	}
}

func TestGetSlotByIdentifier(t *testing.T) {
	tmp := tmpl.NewTemplate("/*<slot name>*/content/*</slot>*/")

	if tmp == nil {
		t.Errorf("Expected template, got nil")
		return
	}

	match, err := tmp.GetSlotMatch("name")

	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
		return
	}

	if match == nil {
		t.Errorf("Expected match, got nil")
		return
	}

	if match.Identifier != "name" {
		t.Errorf("Expected identifier name, got %s", match.Identifier)
	}

	if match.InnerContent != "content" {
		t.Errorf("Expected innerContent content, got %s", match.InnerContent)
	}

	if match.OuterContent != "/*<slot name>*/content/*</slot>*/" {
		t.Errorf("Expected outerContent /*<slot name>*/content/*</slot>*/, got %s", match.OuterContent)
	}
}

func TestReplaceSlot(t *testing.T) {
	tmp := tmpl.NewTemplate("Hello, my name is /*<slot name>*/John Doe/*</slot>*/")

	replaced, err := tmp.ReplaceSlot("name", "Jane Doe")

	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
		return
	}

	if replaced.GetContent() != "Hello, my name is Jane Doe" {
		t.Errorf("Expected replaced template, got %s", replaced.GetContent())
	}
}

func TestReplaceSlotAllOccurrence(t *testing.T) {
	tmp := tmpl.NewTemplate("Hello, my name is /*<slot name>*/John Doe/*</slot>*/ and I am a /*<slot name>*/developer/*</slot>*/")

	replaced, err := tmp.ReplaceSlot("name", "Jane Doe", tmpl.ReplaceSlotOptionAll)

	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
		return
	}

	if replaced.GetContent() != "Hello, my name is Jane Doe and I am a Jane Doe" {
		t.Errorf("Expected replaced template, got '%s'", replaced.GetContent())
	}
}

func TestReplaceSlotMapSingle(t *testing.T) {
	tmp := tmpl.NewTemplate("Hello, my name is /*<slot name>*/John Doe/*</slot>*/ and I am a /*<slot role>*/developer/*</slot>*/")

	replacements := map[string]string{
		"name": "Jane Doe",
		"role": "designer",
	}

	replaced, err := tmp.ReplaceSlotMap(replacements)

	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
		return
	}

	if replaced.GetContent() != "Hello, my name is Jane Doe and I am a designer" {
		t.Errorf("Expected replaced template, got '%s'", replaced.GetContent())
	}
}

func TestReplaceSlotMapMultiple(t *testing.T) {
	tmp := tmpl.NewTemplate(strings.Join([]string{
		"Hello, my name is /*<slot name>*/John Doe/*</slot>*/ and I am a /*<slot role>*/developer/*</slot>*/",
		"Hello, my name is /*<slot name>*/Claire Doe/*</slot>*/ and I am a /*<slot role>*/designer/*</slot>*/",
	}, "\n"))

	replacements := map[string]string{
		"name": "Dane Doe",
		"role": "architect",
	}

	replaced, err := tmp.ReplaceSlotMap(replacements)

	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
		return
	}

	if replaced.GetContent() != strings.Join([]string{
		"Hello, my name is Dane Doe and I am a architect",
		"Hello, my name is Dane Doe and I am a architect",
	}, "\n") {
		t.Errorf("Expected replaced template, got '%s'", replaced.GetContent())
	}
}

func TestReplaceSlotMapSinglePreserveTags(t *testing.T) {
	tmp := tmpl.NewTemplate("Hello, my name is /*<slot name>*/John Doe/*</slot>*/ and I am a /*<slot role>*/developer/*</slot>*/")

	replacements := map[string]string{
		"name": "Jane Doe",
		"role": "designer",
	}

	replaced, err := tmp.ReplaceSlotMap(replacements, tmpl.ReplaceSlotOptionPreserveTags)

	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
		return
	}

	if replaced.GetContent() != "Hello, my name is /*<slot name>*/Jane Doe/*</slot>*/ and I am a /*<slot role>*/designer/*</slot>*/" {
		t.Errorf("Expected replaced template, got '%s'", replaced.GetContent())
	}
}

func TestReplaceSlotAny(t *testing.T) {
	tmp := tmpl.NewTemplate("Hello, my name is /*<slot name>*/John Doe/*</slot>*/ and I am a /*<slot name>*/developer/*</slot>*/")

	replaced := tmp.ReplaceSlotAny("Jane Doe")

	if replaced.GetContent() != "Hello, my name is Jane Doe and I am a Jane Doe" {
		t.Errorf("Expected replaced template, got '%s'", replaced.GetContent())
	}
}

func TestReplaceSlotAnyPreserveTags(t *testing.T) {
	tmp := tmpl.NewTemplate("Hello, my name is /*<slot name>*/John Doe/*</slot>*/ and I am a /*<slot name>*/developer/*</slot>*/")

	replaced := tmp.ReplaceSlotAny("Jane Doe", tmpl.ReplaceSlotOptionPreserveTags)

	if replaced.GetContent() != "Hello, my name is /*<slot name>*/Jane Doe/*</slot>*/ and I am a /*<slot name>*/Jane Doe/*</slot>*/" {
		t.Errorf("Expected replaced template, got '%s'", replaced.GetContent())
	}
}

func TestReplaceSlotFunc(t *testing.T) {
	tmp := tmpl.NewTemplate("Hello, my name is /*<slot name>*/John Doe/*</slot>*/ and I am a /*<slot name>*/developer/*</slot>*/")

	replaced := tmp.ReplaceSlotFunc(func(identifier, innerContent, outerContent string) string {
		return fmt.Sprintf("[%s=%s]", identifier, innerContent)
	})

	if replaced.GetContent() != "Hello, my name is [name=John Doe] and I am a [name=developer]" {
		t.Errorf("Expected replaced template, got '%s'", replaced.GetContent())
	}
}
