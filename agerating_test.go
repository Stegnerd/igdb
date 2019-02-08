package igdb

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

const (
	testAgeRatingGet  string = "test_data/agerating_get.json"
	testAgeRatingList string = "test_data/agerating_list.json"
)

func TestAgeRatingService_Get(t *testing.T) {
	f, err := ioutil.ReadFile(testAgeRatingGet)
	if err != nil {
		t.Fatal(err)
	}

	init := make([]*AgeRating, 1)
	json.Unmarshal(f, &init)

	var tests = []struct {
		name          string
		file          string
		id            int
		opts          []FuncOption
		wantAgeRating *AgeRating
		wantErr       error
	}{
		{"Valid response", testAgeRatingGet, 9644, []FuncOption{SetFields("name")}, init[0], nil},
		{"Invalid ID", testFileEmpty, -1, nil, nil, ErrNegativeID},
		{"Empty response", testFileEmpty, 9644, nil, nil, errInvalidJSON},
		{"Invalid option", testFileEmpty, 9644, []FuncOption{SetOffset(99999)}, nil, ErrOutOfRange},
		{"No results", testFileEmptyArray, 0, nil, nil, ErrNoResults},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c, err := testServerFile(http.StatusOK, test.file)
			if err != nil {
				t.Fatal(err)
			}
			defer ts.Close()

			g, err := c.AgeRatings.Get(test.id, test.opts...)
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			if !reflect.DeepEqual(g, test.wantAgeRating) {
				t.Errorf("got: <%v>, \nwant: <%v>", g, test.wantAgeRating)
			}
		})
	}
}

func TestAgeRatingService_List(t *testing.T) {
	f, err := ioutil.ReadFile(testAgeRatingList)
	if err != nil {
		t.Fatal(err)
	}

	init := make([]*AgeRating, 0)
	json.Unmarshal(f, &init)

	var tests = []struct {
		name           string
		file           string
		ids            []int
		opts           []FuncOption
		wantAgeRatings []*AgeRating
		wantErr        error
	}{
		{"Valid response", testAgeRatingList, []int{9644, 40}, []FuncOption{SetLimit(5)}, init, nil},
		{"Zero IDs", testFileEmpty, nil, nil, nil, ErrEmptyIDs},
		{"Invalid ID", testFileEmpty, []int{-500}, nil, nil, ErrNegativeID},
		{"Empty response", testFileEmpty, []int{9644, 40}, nil, nil, errInvalidJSON},
		{"Invalid option", testFileEmpty, []int{9644, 40}, []FuncOption{SetOffset(99999)}, nil, ErrOutOfRange},
		{"No results", testFileEmptyArray, []int{0, 9999999}, nil, nil, ErrNoResults},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c, err := testServerFile(http.StatusOK, test.file)
			if err != nil {
				t.Fatal(err)
			}
			defer ts.Close()

			g, err := c.AgeRatings.List(test.ids, test.opts...)
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			if !reflect.DeepEqual(g, test.wantAgeRatings) {
				t.Errorf("got: <%v>, \nwant: <%v>", g, test.wantAgeRatings)
			}
		})
	}
}

func TestAgeRatingService_Index(t *testing.T) {
	f, err := ioutil.ReadFile(testAgeRatingList)
	if err != nil {
		t.Fatal(err)
	}

	init := make([]*AgeRating, 0)
	json.Unmarshal(f, &init)

	tests := []struct {
		name           string
		file           string
		opts           []FuncOption
		wantAgeRatings []*AgeRating
		wantErr        error
	}{
		{"Valid response", testAgeRatingList, []FuncOption{SetLimit(5)}, init, nil},
		{"Empty response", testFileEmpty, nil, nil, errInvalidJSON},
		{"Invalid option", testFileEmpty, []FuncOption{SetOffset(99999)}, nil, ErrOutOfRange},
		{"No results", testFileEmptyArray, nil, nil, ErrNoResults},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c, err := testServerFile(http.StatusOK, test.file)
			if err != nil {
				t.Fatal(err)
			}
			defer ts.Close()

			g, err := c.AgeRatings.Index(test.opts...)
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			if !reflect.DeepEqual(g, test.wantAgeRatings) {
				t.Errorf("got: <%v>, \nwant: <%v>", g, test.wantAgeRatings)
			}
		})
	}
}

func TestAgeRatingService_Count(t *testing.T) {
	var tests = []struct {
		name      string
		resp      string
		opts      []FuncOption
		wantCount int
		wantErr   error
	}{
		{"Happy path", `{"count": 100}`, []FuncOption{SetFilter("popularity", OpGreaterThan, "75")}, 100, nil},
		{"Empty response", "", nil, 0, errInvalidJSON},
		{"Invalid option", "", []FuncOption{SetLimit(100)}, 0, ErrOutOfRange},
		{"No results", "[]", nil, 0, ErrNoResults},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c := testServerString(http.StatusOK, test.resp)
			defer ts.Close()

			count, err := c.AgeRatings.Count(test.opts...)
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			if count != test.wantCount {
				t.Fatalf("got: <%v>, want: <%v>", count, test.wantCount)
			}
		})
	}
}

func TestAgeRatingService_Fields(t *testing.T) {
	var tests = []struct {
		name       string
		resp       string
		wantFields []string
		wantErr    error
	}{
		{"Happy path", `["name", "slug", "url"]`, []string{"url", "slug", "name"}, nil},
		{"Dot operator", `["logo.url", "background.id"]`, []string{"background.id", "logo.url"}, nil},
		{"Asterisk", `["*"]`, []string{"*"}, nil},
		{"Empty response", "", nil, errInvalidJSON},
		{"No results", "[]", nil, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, c := testServerString(http.StatusOK, test.resp)
			defer ts.Close()

			fields, err := c.AgeRatings.Fields()
			if errors.Cause(err) != test.wantErr {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(err), test.wantErr)
			}

			ok, err := equalSlice(fields, test.wantFields)
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				t.Fatalf("Expected fields '%v', got '%v'", test.wantFields, fields)
			}
		})
	}
}
