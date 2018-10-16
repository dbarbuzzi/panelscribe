package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildQuoteGroupsSuccess(t *testing.T) {
	// this will test successful calls
	type fields struct {
		groups string
		flip   bool
	}
	var tests = []struct {
		desc     string
		input    fields
		expected string
	}{
		{
			desc:     "short, simple",
			input:    fields{"21", false},
			expected: "\n---\n\n> \n\n> \n\n---\n\n> \n",
		},
		{
			desc:     "short, simple, flipped",
			input:    fields{"21", true},
			expected: "\n> \n\n> \n\n---\n\n> \n\n---\n",
		},
		{
			desc:     "medium, simple",
			input:    fields{"2143", false},
			expected: "\n---\n\n> \n\n> \n\n---\n\n> \n\n---\n\n> \n\n> \n\n> \n\n> \n\n---\n\n> \n\n> \n\n> \n",
		},
	}

	for _, test := range tests {
		actual, err := buildQuoteGroups(test.input.groups, test.input.flip)
		if err != nil {
			t.Errorf("%s: received error '%s', expected success", test.desc, err.Error())
		}
		if diff := cmp.Diff(test.expected, actual); diff != "" {
			t.Errorf("%s: quote group template differs: (-want +got)\n%s", test.desc, diff)
		}
	}
}

func TestBuildQuoteGroupsFail(t *testing.T) {
	// this will test error-producing calls
	type fields struct {
		groups string
		flip   bool
	}
	var tests = []struct {
		desc     string
		input    fields
		expected string
	}{
		{
			desc:     "short, simple",
			input:    fields{" 1", false},
			expected: `strconv.Atoi: parsing " ": invalid syntax`,
		},
		{
			desc:     "short, simple, flipped",
			input:    fields{"2.1", true},
			expected: `strconv.Atoi: parsing ".": invalid syntax`,
		},
	}

	for _, test := range tests {
		_, err := buildQuoteGroups(test.input.groups, test.input.flip)
		if err == nil {
			t.Errorf("%s: received error was nil, expected failure", test.desc)
		} else if diff := cmp.Diff(test.expected, err.Error()); diff != "" {
			t.Errorf("%s: returned error differs: (-want +got)\n%s", test.desc, diff)
		}
	}
}

func TestValidateQuoteGroupArgFail(t *testing.T) {
	// test cases where error is expected
	var tests = []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "non-integer number (float)",
			input:    "1.2",
			expected: "arg must be an integer, got 1.2",
		},
		{
			desc:     "non-integer number (trailing space)",
			input:    "12 ",
			expected: "arg must be an integer, got 12 ",
		},
		{
			desc:     "non-integer number (alpha)",
			input:    "1e2",
			expected: "arg must be an integer, got 1e2",
		},
		{
			desc:     "zero at beginning",
			input:    "059",
			expected: "arg should only contain positive integers, had '0' at position 0",
		},
		{
			desc:     "zero at end",
			input:    "340",
			expected: "arg should only contain positive integers, had '0' at position 2",
		},
		{
			desc:     "zero in the middle",
			input:    "307",
			expected: "arg should only contain positive integers, had '0' at position 1",
		},
	}

	for _, test := range tests {
		err := validateQuoteGroupArg(test.input)
		if diff := cmp.Diff(test.expected, err.Error()); diff != "" {
			t.Errorf("%s: returned error differs: (-want +got)\n%s", test.desc, diff)
		}
	}
}
func TestValidateQuoteGroupArgSuccess(t *testing.T) {
	// test cases where validation should pass
	var tests = []struct {
		desc     string
		input    string
		expected error
	}{
		{
			desc:     "short value",
			input:    "3",
			expected: nil,
		},
		{
			desc:     "medium value",
			input:    "382",
			expected: nil,
		},
		{
			desc:     "large value",
			input:    "26642",
			expected: nil,
		},
	}

	for _, test := range tests {
		err := validateQuoteGroupArg(test.input)
		if diff := cmp.Diff(test.expected, err); diff != "" {
			t.Errorf("%s: returned error differs: (-want +got)\n%s", test.desc, diff)
		}
	}
}
