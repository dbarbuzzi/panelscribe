package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateBaseToml(t *testing.T) {
	var tests = []struct {
		desc     string
		input    int
		expected SmudgeToml
	}{
		{
			desc:  "single page",
			input: 1,
			expected: SmudgeToml{
				Tags:   []string{},
				Chars:  map[string]string{},
				Scenes: map[string][][]string{"p1": [][]string{}}},
		},
		{
			desc:  "two pages",
			input: 2,
			expected: SmudgeToml{
				Tags:   []string{},
				Chars:  map[string]string{},
				Scenes: map[string][][]string{"p1": [][]string{}, "p2": [][]string{}}},
		},
	}

	for _, test := range tests {
		actual := createBaseToml(test.input)
		if diff := cmp.Diff(test.expected, actual); diff != "" {
			t.Errorf("%s: base toml differs: (-want +got)\n%s", test.desc, diff)
		}
	}
}

func TestCreateBaseMd(t *testing.T) {
	type fields struct {
		name  string
		pages int
	}
	var tests = []struct {
		desc     string
		input    fields
		expected string
	}{
		{
			desc:     "simple name, single page",
			input:    fields{"comic", 1},
			expected: "# comic\n\n\nupper_text\n\n---\n\nlower_text\n\n## Page 1 <page1.jpg>\n",
		},
		{
			desc:     "simple name, two pages",
			input:    fields{"comic", 2},
			expected: "# comic\n\n\nupper_text\n\n---\n\nlower_text\n\n## Page 1 <page1.jpg>\n\n\nupper_text\n\n---\n\nlower_text\n\n## Page 2 <page2.jpg>\n",
		},
		{
			desc:     "simple name, five pages",
			input:    fields{"Hero Man Stories – Hero Man 1: Hero Man Enters", 5},
			expected: "# Hero Man Stories – Hero Man 1: Hero Man Enters\n\n\nupper_text\n\n---\n\nlower_text\n\n## Page 1 <page1.jpg>\n\n\nupper_text\n\n---\n\nlower_text\n\n## Page 2 <page2.jpg>\n\n\nupper_text\n\n---\n\nlower_text\n\n## Page 3 <page3.jpg>\n\n\nupper_text\n\n---\n\nlower_text\n\n## Page 4 <page4.jpg>\n\n\nupper_text\n\n---\n\nlower_text\n\n## Page 5 <page5.jpg>\n",
		},
	}

	for _, test := range tests {
		actual := createBaseMd(test.input.name, test.input.pages)
		if diff := cmp.Diff(test.expected, actual); diff != "" {
			t.Errorf("%s: base md differs: (-want +got)\n%s", test.desc, diff)
		}
	}
}
