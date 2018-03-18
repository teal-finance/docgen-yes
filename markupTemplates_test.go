package docgen

import "testing"

func TestBaseTemplate(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BaseTemplate(); got != tt.want {
				t.Errorf("BaseTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnorderedList(t *testing.T) {
	type args struct {
		listItems string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnorderedList(tt.args.listItems); got != tt.want {
				t.Errorf("UnorderedList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedList(t *testing.T) {
	type args struct {
		listItems string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrderedList(tt.args.listItems); got != tt.want {
				t.Errorf("OrderedList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListItem(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListItem(tt.args.text); got != tt.want {
				t.Errorf("ListItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Div(tt.args.text); got != tt.want {
				t.Errorf("Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestP(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
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
	// TODO: Add test cases.
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
	// TODO: Add test cases.
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FaviconIcoData(); got != tt.want {
				t.Errorf("FaviconIcoData() = %v, want %v", got, tt.want)
			}
		})
	}
}
