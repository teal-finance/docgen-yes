package docgen

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseTemplate(t *testing.T) {
	t.Log("Testing Base Template")

	// check Base Template for expected template fields: {title}, {css}, {favicon.ico}, {intro}, {routes}
	expectedFields := []string{"{title}", "{css}", "{favicon.ico}", "{intro}", "{routes}"}
	bt := BaseTemplate()
	for _, field := range expectedFields {
		assert.True(t, strings.Contains(bt, field))
	}

	t.Log("Base Template Testing Complete")
}

func TestListItem(t *testing.T) {
	t.Log("Testing List Item")

	itemText := "test"
	li := ListItem(itemText)
	assert.True(t, strings.Contains(li, "<li>"))
	assert.True(t, strings.Contains(li, "</li>"))
	assert.True(t, strings.Contains(li, itemText))

	t.Log("List Item Testing Complete")
}

func TestOrderedList(t *testing.T) {
	t.Log("Testing Ordered List")

	li := ListItem("test")
	ol := OrderedList(li)

	assert.True(t, strings.Contains(ol, "<ol>"))
	assert.True(t, strings.Contains(ol, "</ol>"))
	assert.True(t, strings.Contains(ol, li))

	t.Log("Ordered List Testing Complete")
}

func TestUnorderedList(t *testing.T) {
	t.Log("Testing Unordered List")

	li := ListItem("test")
	ul := UnorderedList(li)

	assert.True(t, strings.Contains(ul, "<ul>"))
	assert.True(t, strings.Contains(ul, "</ul>"))
	assert.True(t, strings.Contains(ul, li))

	t.Log("Ordered List Testing Complete")
}

func TestDiv(t *testing.T) {
	t.Log("Testing Div")
	testString := "test"
	div := Div(testString)
	assert.True(t, strings.Contains(div, "<div>"))
	assert.True(t, strings.Contains(div, "</div>"))
	assert.True(t, strings.Contains(div, testString))
	t.Log("Div Testing Complete")
}

func TestP(t *testing.T) {
	t.Log("Testing P")
	testString := "test"
	p := P(testString)
	assert.True(t, strings.Contains(p, "<p>"))
	assert.True(t, strings.Contains(p, "</p>"))
	assert.True(t, strings.Contains(p, testString))
	t.Log("P Testing Complete")
}

func TestHeading(t *testing.T) {
	t.Log("Testing Headings")
	testString := "test"

	// h1 - h6 are valid
	for index := 1; index <= 6; index++ {
		h := Head(index, testString)

		headerOpenTag := fmt.Sprintf("<h%v>", index)
		headerCloseTag := fmt.Sprintf("</h%v>", index)

		assert.True(t, strings.Contains(h, headerOpenTag))
		assert.True(t, strings.Contains(h, headerCloseTag))
		assert.True(t, strings.Contains(h, testString))
	}

	// anything below 1 should return an h1 tag
	hzero := Head(0, "zero | should return h1")
	assert.True(t, strings.Contains(hzero, "<h1>"))
	// anything above 6 should return an h6 tag
	hseven := Head(7, "seven | should return h6")
	assert.True(t, strings.Contains(hseven, "<h6>"))

	t.Log("Heading Testing Complete")
}

func TestCSS(t *testing.T) {
	// MilligramMinCSS
	css := MilligramMinCSS()
	assert.NotNil(t, css)
	// TODO: test CSS better?
}

func TestBaseTemplate1(t *testing.T) {
	tests := []struct {
		name       string
		want       string
		expectFail bool
	}{
		// TODO: Add test cases.
		{"baseTest", BaseTemplate(), false},
		{"failBlank", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BaseTemplate(); got != tt.want && !tt.expectFail {
				t.Errorf("BaseTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnorderedList1(t *testing.T) {
	type args struct {
		listItems string
	}
	tests := []struct {
		name       string
		args       args
		want       string
		expectFail bool
	}{
		{"unorderedListTest", args{"<li>test</li>"}, UnorderedList(ListItem("test")), false},
		{"failBlank", args{"<li>test</li>"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnorderedList(tt.args.listItems); got != tt.want && !tt.expectFail {
				t.Errorf("UnorderedList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedList1(t *testing.T) {
	type args struct {
		listItems string
	}
	tests := []struct {
		name       string
		args       args
		want       string
		expectFail bool
	}{
		{"orderedListTest", args{"<li>test</li>"}, OrderedList(ListItem("test")), false},
		{"failBlank", args{"<li>test</li>"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrderedList(tt.args.listItems); got != tt.want && !tt.expectFail {
				t.Errorf("OrderedList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListItem1(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"baseListItem", args{"test"}, "<li>test</li>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListItem(tt.args.text); got != tt.want {
				t.Errorf("ListItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiv1(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"baseDivTest", args{"test"}, "<div>test</div>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Div(tt.args.text); got != tt.want {
				t.Errorf("Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestP1(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"basePTest", args{"test"}, "<p>test</p>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := P(tt.args.text); got != tt.want {
				t.Errorf("P() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHead(t *testing.T) {
	type args struct {
		level int
		text  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"baseHeadTest", args{1, "test"}, "<h1>test</h1>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Head(tt.args.level, tt.args.text); got != tt.want {
				t.Errorf("Head() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMilligramMinCSS(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"baseCssTest", MilligramMinCSS()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MilligramMinCSS(); got != tt.want {
				t.Errorf("MilligramMinCSS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFaviconIcoData(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"baseFaviconTest", FaviconIcoData()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FaviconIcoData(); got != tt.want {
				t.Errorf("FaviconIcoData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBassCSS(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"baseBassCSS", BassCSS()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BassCSS(); got != tt.want {
				t.Errorf("BassCSS() = %v, want %v", got, tt.want)
			}
		})
	}
}
