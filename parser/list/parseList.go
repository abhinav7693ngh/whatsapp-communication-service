package list

import (
	"multiBot/constants"
	"multiBot/types/list"

	"github.com/google/uuid"
)

type ListHeader struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ListBody struct {
	Text string `json:"text"`
}

type ListFooter struct {
	Text string `json:"text"`
}

type ListSectionRows struct {
	Id          string `json:"id"`
	Title       string `json:"title" validate:"required" message:"title is required field"`
	Description string `json:"description"`
}

type ListSection struct {
	Title string            `json:"title" validate:"required" message:"title is required field"`
	Rows  []ListSectionRows `json:"rows" validate:"required" message:"rows is required field"`
}

type ListAction struct {
	Button   string        `json:"button"`
	Sections []ListSection `json:"sections"`
}

type ListParsedBody struct {
	Type   string     `json:"type"`
	Header ListHeader `json:"header"`
	Body   ListBody   `json:"body"`
	Footer ListFooter `json:"footer"`
	Action ListAction `json:"action"`
}

func (l *ListParsedBody) fill_defaults() {
	l.Type = "list"
}

func makeListSectionRows(rows []list.ListSectionRows) []ListSectionRows {
	var sectionRows []ListSectionRows
	for _, r := range rows {
		var row ListSectionRows
		row.Id = uuid.New().String()
		row.Title = r.Title
		row.Description = r.Description
		sectionRows = append(sectionRows, row)
	}
	return sectionRows
}

func makeSections(sec []list.ListSection) []ListSection {
	var sections []ListSection
	for _, s := range sec {
		var section ListSection
		section.Title = s.Title
		section.Rows = makeListSectionRows(s.Rows)
		sections = append(sections, section)
	}
	return sections
}

func ParseList(list list.ListDataStruct) ListParsedBody {
	var parsedBody ListParsedBody
	parsedBody.fill_defaults()
	parsedBody.Header.Type = string(constants.LIST_MSG_TYPE_TEXT)
	parsedBody.Header.Text = list.Header
	parsedBody.Body.Text = list.Body
	parsedBody.Footer.Text = list.Footer
	parsedBody.Action.Button = list.Button
	parsedBody.Action.Sections = makeSections(list.Sections)
	return parsedBody
}
