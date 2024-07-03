package tmpl

import (
	"fmt"
	"regexp"
	"slices"
)

type Template struct {
	content        string
	slotMatchCache map[string][]SlotMatch
}

type ReplaceSlotOption int

const (
	ReplaceSlotOptionAll ReplaceSlotOption = iota
	ReplaceSlotOptionPreserveTags
)

type SlotMatch struct {
	StartIndex      int
	EndIndex        int
	InnerStartIndex int
	InnerEndIndex   int
	Identifier      string
	InnerContent    string
	OuterContent    string
}

func NewTemplate(template string) *Template {
	return &Template{
		content:        template,
		slotMatchCache: make(map[string][]SlotMatch),
	}
}

func (t *Template) GetContent() string {
	return t.content
}

func (t *Template) GetSlotMatches() []SlotMatch {
	if len(t.slotMatchCache) > 0 {
		var matches []SlotMatch
		for _, match := range t.slotMatchCache {
			matches = append(matches, match...)
		}
		return matches
	}

	startPattern := regexp.MustCompile(`/\*<slot\s*(\w+)>\*/`)
	endPattern := regexp.MustCompile(`/\*<\/\s*slot>\*/`)

	tmp := t.content
	var matches []SlotMatch

	offset := 0

	for {
		startMatch := startPattern.FindStringSubmatchIndex(tmp)
		if startMatch == nil {
			break
		}

		endMatch := endPattern.FindStringSubmatchIndex(tmp)
		if endMatch == nil {
			break
		}

		match := SlotMatch{
			StartIndex:      startMatch[0] + offset,
			EndIndex:        endMatch[1] + offset,
			InnerStartIndex: startMatch[1] + offset,
			InnerEndIndex:   endMatch[0] + offset,
			Identifier:      tmp[startMatch[2]:startMatch[3]],
			InnerContent:    tmp[startMatch[1]:endMatch[0]],
			OuterContent:    tmp[startMatch[0]:endMatch[1]],
		}

		matches = append(matches, match)
		t.slotMatchCache[match.Identifier] = append(t.slotMatchCache[match.Identifier], match)

		tmp = tmp[endMatch[1]:]
		offset += endMatch[1]
	}

	return matches
}

func (t *Template) GetSlotMatchesReverseOrder() []SlotMatch {
	matches := t.GetSlotMatches()
	slices.SortFunc(matches, reverseCompare)
	return matches
}

func (t *Template) GetSlotMatch(identifier string) (*SlotMatch, error) {
	if len(t.slotMatchCache) == 0 {
		t.GetSlotMatches()
	}

	if match, ok := t.slotMatchCache[identifier]; ok {
		return &match[0], nil
	}

	return nil, fmt.Errorf("slot with identifier %s not found", identifier)
}

func (t *Template) GetSlotsByIdentifier(identifier string) []SlotMatch {
	if len(t.slotMatchCache) == 0 {
		t.GetSlotMatches()
	}

	if match, ok := t.slotMatchCache[identifier]; ok {
		return match
	}

	return nil
}

// ReplaceSlot replaces the content of a slot with the given identifier.
// It returns a new Template with the slot replaced.
// e.g. content = "Hello, my name is /*<slot name>*/John Doe/*</slot>*/"
//
//	ReplaceSlot("name", "Jane Doe") -> "Hello, my name is Jane Doe"
//
// default behavior is to replace the first slot found.
// options:
//   - ReplaceSlotOptionAll: replace all slots with the given identifier
func (t *Template) ReplaceSlot(identifier, content string, options ...ReplaceSlotOption) (*Template, error) {
	matches := t.GetSlotMatches()
	tmp := t.content

	if slices.Contains(options, ReplaceSlotOptionAll) {
		slices.SortFunc(matches, reverseCompare)
		for _, match := range matches {
			if match.Identifier == identifier {
				tmp = tmp[:match.StartIndex] + content + tmp[match.EndIndex:]
			}
		}
		return NewTemplate(tmp), nil
	} else {
		for _, match := range matches {
			if match.Identifier == identifier {
				tmp = tmp[:match.StartIndex] + content + tmp[match.EndIndex:]
				return NewTemplate(tmp), nil
			}
		}
	}

	return nil, fmt.Errorf("slot with identifier %s not found", identifier)
}

func (t *Template) ReplaceSlotMust(identifier, content string, options ...ReplaceSlotOption) *Template {
	newTemplate, err := t.ReplaceSlot(identifier, content, options...)
	if err != nil {
		panic(err)
	}

	return newTemplate
}

func (t *Template) ReplaceSlotMap(replacementMap map[string]string, options ...ReplaceSlotOption) (*Template, error) {
	matches := t.GetSlotMatches()

	slices.SortFunc(matches, reverseCompare)

	content := t.content

	for _, match := range matches {
		if replacement, ok := replacementMap[match.Identifier]; ok {
			if slices.Contains(options, ReplaceSlotOptionPreserveTags) {
				content = content[:match.InnerStartIndex] + replacement + content[match.InnerEndIndex:]
			} else {
				content = content[:match.StartIndex] + replacement + content[match.EndIndex:]
			}
		}
	}

	return NewTemplate(content), nil
}

func (t *Template) ReplaceSlotFunc(replacementFunc func(identifier, innerContent, outerContent string) string) *Template {
	matches := t.GetSlotMatches()

	slices.SortFunc(matches, reverseCompare)

	content := t.content

	for _, match := range matches {
		replacement := replacementFunc(match.Identifier, match.InnerContent, match.OuterContent)
		content = content[:match.StartIndex] + replacement + content[match.EndIndex:]
	}

	return NewTemplate(content)
}

func (t *Template) ReplaceSlotAny(replacement string, options ...ReplaceSlotOption) *Template {
	matches := t.GetSlotMatches()

	slices.SortFunc(matches, reverseCompare)

	content := t.content

	for _, match := range matches {
		if slices.Contains(options, ReplaceSlotOptionPreserveTags) {
			content = content[:match.InnerStartIndex] + replacement + content[match.InnerEndIndex:]
		} else {
			content = content[:match.StartIndex] + replacement + content[match.EndIndex:]
		}
	}

	return NewTemplate(content)
}

func reverseCompare(i, j SlotMatch) int {
	return j.StartIndex - i.StartIndex
}
