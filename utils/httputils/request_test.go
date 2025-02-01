package httputils_test

import (
	"net/url"
	"testing"

	"github.com/ARF-DEV/image-processing-api/utils"
	"github.com/ARF-DEV/image-processing-api/utils/httputils"
)

type BodyTest struct {
	Name string `json:"name"`
}
type TestStruct struct {
	Number  int64    `form:"number"`
	Decimal float64  `form:"decimal"`
	Body    BodyTest `form:"body"`
}

func TestParseFormData(t *testing.T) {
	// r := http.Request{
	// 	Header: http.Header{
	// 		"Content-Type": []string{"multipart/form-data"},
	// 	},
	// 	Body: nil,
	// 	Form: url.Values{},
	// }
	f := url.Values{}
	f.Set("number", "1")
	f.Set("decimal", "100.23")
	f.Set("body", `{"name": "testing 1"}`)

	dest := TestStruct{}
	if err := httputils.ParseURLValues(f, &dest); err != nil {
		t.Fatal(err)
	}

	if dest.Number != 1 {
		t.Fatalf("error expected %v, but got %v", 1, dest.Number)
	}
	if dest.Decimal != 100.23 {
		t.Fatalf("error expected %v, but got %v", 100.23, dest.Decimal)
	}
	if dest.Body.Name != "testing 1" {
		t.Fatalf("error expected %v, but got %v", "testing 1", dest.Body.Name)

	}

	utils.PrintInJSONFormat(dest)
}
