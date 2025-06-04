package srctotxt

import "testing"
import "github.com/tree-sitter/go-tree-sitter"
import "slices"

//import "fmt"

func TestCSGetFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `public int Add(int x, int y) {
                    return x + y;
                }`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	err := parser.SetLanguage(lang)
	if err != nil {
		t.Errorf("TestCSGetFunction failed to set language: %v", err)
	}
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if len(slice) != 1 {
		t.Errorf("TestCSGetFunction failed length wrong")
	}
	if slice[0] != "public int Add(int x, int y)" {
		t.Errorf(
			"TestCSGetFunction failed want %q got %v",
			"public int Add(int x, int y)",
			slice[0],
		)
	}
}

func TestCSGetMultipleFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `public int Subtract(int x, int y) {
                    return x - y;
                }
                
                public string Add(string x, int y) {
                    return x + y.ToString();
                }`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	err := parser.SetLanguage(lang)
	if err != nil {
		t.Errorf("TestCSGetFunction failed to set language: %v", err)
	}
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "public int Subtract(int x, int y)") {
		t.Errorf(
			"TestCSGetMultipleFunction failed want %q got %v",
			"public int Subtract(int x, int y)",
			slice,
		)
	}
	if !slices.Contains(slice, "public string Add(string x, int y)") {
		t.Errorf(
			"TestCSGetMultipleFunction failed want %q got %v",
			"public string Add(string x, int y)",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestCSGetMultipleFunction failed length wrong")
	}
}

func TestCSGetNoFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := ``
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	err := parser.SetLanguage(lang)
	if err != nil {
		t.Errorf("TestCSGetFunction failed to set language: %v", err)
	}
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if len(slice) != 0 {
		t.Errorf("TestCSGetNoFunction failed length wrong")
	}
}

func TestCSGetEmptyFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `public void Test() {}`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if len(slice) != 1 {
		t.Errorf("TestCSGetEmptyFunction failed length wrong")
	}
	if !slices.Contains(slice, "public void Test()") {
		t.Errorf(
			"TestCSGetEmptyFunction failed want %q got %v",
			"public void Test()",
			slice,
		)
	}
}

func TestCSEmptyFile(t *testing.T) {
	InitsrcToText(Test())
	testcode := ``
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if len(slice) != 0 {
		t.Errorf("TestCSEmptyFile failed length wrong")
	}

}

func TestCSGetNestedFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `public void Nested(int x) {
                    void Nested2(int y) {
                        Console.WriteLine(y);
                    }
                    Nested2(x);
                }`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "public void Nested(int x)") {
		t.Errorf(
			"TestCSGetNestedFunction failed want %q got %v",
			"public void Nested(int x)",
			slice,
		)
	}
	if !slices.Contains(slice, "void Nested2(int y)") {
		t.Errorf(
			"TestCSGetNestedFunction failed want %q got %v",
			"void Nested2(int y)",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestCSGetNestedFunction failed length wrong")
	}
}

func TestCSGetFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `public int Subtract(int x, int y) {
                    return x - y;
                }`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if string(slice[0]) != "int int int" {
		t.Errorf(
			"TestCSGetFunctionSignature failed want %q got %v",
			"int int int",
			slice[0],
		)
	}
}

func TestCSGetMultipleFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `public int Subtract(int x, int y) {
                    return x - y;
                }
                
                public string Add(string x, int y) {
                    return x + y.ToString();
                }`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "int int int") {
		t.Errorf(
			"TestCSGetMultipleFunctionSignature failed want %q got %v",
			"int int int",
			slice,
		)
	}
	if !slices.Contains(slice, "string int string") {
		t.Errorf(
			"TestCSGetMultipleFunctionSignature failed want %q got %v",
			"string int string",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestCSGetMultipleFunctionSignature failed length wrong")
	}
}

func TestCSGetNoFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := ``
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if len(slice) != 0 {
		t.Errorf("TestCSGetNoFunctionSignature failed length wrong")
	}
}

func TestCSGetEmptyFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `public void Test() {}`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if len(slice) != 1 {
		t.Errorf("TestCSGetEmptyFunctionSignature failed length wrong")
	}
	if !slices.Contains(slice, "void") {
		t.Errorf(
			"TestCSGetEmptyFunctionSignature failed want %q got %v",
			"void",
			slice,
		)
	}
}

func TestCSGetDocs(t *testing.T) {
	InitsrcToText(Test())
	testcode := `
/* <summary>
Adds two numbers.
</summary>
<param name="x">An integer.</param>
<param name="y">Another integer.</param>
<returns>The sum of x and y.</returns>*/
public int Add(int x, int y) {
    return x + y;
}`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindDocs([]byte(testcode), parser, conf)
	if !slices.Contains(slice, `/* <summary>
Adds two numbers.
</summary>
<param name="x">An integer.</param>
<param name="y">Another integer.</param>
<returns>The sum of x and y.</returns>*/`) {
		t.Errorf(
			"TestCSGetDocs failed want %q got %v",
			"/// <summary> Adds two numbers...",
			slice,
		)
	}
	if len(slice) != 1 {
		t.Errorf("TestCSGetDocs failed length wrong")
	}
}

func TestCSGetClass(t *testing.T) {
	InitsrcToText(Test())
	testcode := `public class Test {
		public int Add(int x, int y) {
			return x + y;
		}
		public int Subtract(int x, int y) {
			return x - y;
		}		
	}`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "Test int int int") {
		t.Errorf(
			"TestCSGetClass failed want %q got %v",
			"Test int int int",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestCSGetClass failed length wrong")
	}
}

func TestCSGetLineDocs(t *testing.T) {
	InitsrcToText(Test())
	testcode := `/// Adds two numbers.
/// <param name="x">An integer.</param>
/// <param name="y">Another integer.</param>
public int Add(int x, int y) {
    return x + y;
}
`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindDocs([]byte(testcode), parser, conf)
	if slices.Contains(slice, `/// Adds two numbers.
/// <param name="x">An integer.</param>
/// <param name="y">Another integer.</param>`) {
		t.Errorf(
			"TestCSGetLineDocs failed want %q got %v",
			"/// Adds two numbers...",
			slice,
		)
	}
	if len(slice) != 1 {
		t.Errorf("TestCSGetLineDocs failed length wrong")
	}
}

func TestCSGetMultipleDocs(t *testing.T) {
	InitsrcToText(Test())
	testcode := `/// Adds two numbers.
/// <param name="x">An integer.</param>
/// <param name="y">Another integer.</param>
public int Add(int x, int y) {
    return x + y;
}

/// Subtracts two numbers.
/// <param name="x">An integer.</param>
/// <param name="y">Another integer.</param>
public int Subtract(int x, int y) {
	return x - y;
}
`
	parser := tree_sitter.NewParser()
	conf := config[".cs"]
	lang, _ := getLanguageFileExt(".cs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindDocs([]byte(testcode), parser, conf)
	if slices.Contains(slice, `/// Adds two numbers.
/// <param name="x">An integer.</param>
/// <param name="y">Another integer.</param>`) {
		t.Errorf(
			"TestCSGetMultipleDocs failed want %q got %v",
			"/// Adds two numbers...",
			slice,
		)
	}
	if !slices.Contains(slice, `/// Subtracts two numbers.
// / <param name="x">An integer.</param>/// <param name="y">Another
// integer.</param>`) {
		t.Errorf(
			"TestCSGetMultipleDocs failed want %q got %v",
			"/// Subtracts two numbers...",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestCSGetMultipleDocs failed length wrong")
	}
}
