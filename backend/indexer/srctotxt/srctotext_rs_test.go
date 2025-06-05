//nolint:all
package srctotxt

import "testing"
import "github.com/tree-sitter/go-tree-sitter"
import "slices"

func TestRSGetFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `fn add(x: i32, y: i32) -> i32 {
                    x + y
                }`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if string(slice[0]) != "fn add(x: i32, y: i32) -> i32" {
		t.Errorf(
			"TestRSGetFunction failed want %q got %v",
			"fn add(x: i32, y: i32) -> i32",
			slice[0],
		)
	}
}

func TestRSGetMultipleFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `fn sub(x: i32, y: i32) -> i32 {
                    x - y
                }
                
                fn add(x: String, y: i32) -> String {
                    x + &y.to_string()
                }`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "fn sub(x: i32, y: i32) -> i32") {
		t.Errorf(
			"TestRSGetMultipleFunction failed want %q got %v",
			"fn sub(x: i32, y: i32) -> i32",
			slice,
		)
	}
	if !slices.Contains(slice, "fn add(x: String, y: i32) -> String") {
		t.Errorf(
			"TestRSGetMultipleFunction failed want %q got %v",
			"fn add(x: String, y: i32) -> String",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestRSGetMultipleFunction failed length wrong")
	}
}

func TestRSGetNoFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := ``
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if len(slice) != 0 {
		t.Errorf("TestRSGetNoFunction failed length wrong")
	}
}

func TestRSGetEmptyFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `fn test() {}`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if len(slice) != 1 {
		t.Errorf("TestRSGetEmptyFunction failed length wrong")
	}
	if !slices.Contains(slice, "fn test()") {
		t.Errorf(
			"TestRSGetEmptyFunction failed want %q got %v",
			"fn test()",
			slice,
		)
	}
}

func TestRSGetNestedFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `fn nested(x: i32) {
                    fn nested2(y: i32) -> i32 {
                        y
                    }
                    nested2(x);
                }`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "fn nested(x: i32)") {
		t.Errorf(
			"TestRSGetNestedFunction failed want %q got %v",
			"fn nested(x: i32)",
			slice,
		)
	}
	if !slices.Contains(slice, "fn nested2(y: i32) -> i32") {
		t.Errorf(
			"TestRSGetNestedFunction failed want %q got %v",
			"fn nested2(y: i32) -> i32",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestRSGetNestedFunction failed length wrong")
	}
}

func TestRSGetFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `fn sub(x: i32, y: i32) -> i32 {
                    x - y
                }`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if string(slice[0]) != "i32 i32 i32" {
		t.Errorf(
			"TestRSGetFunctionSignature failed want %q got %v",
			"i32 i32 i32",
			slice[0],
		)
	}
}

func TestRSGetMultipleFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `fn sub(x: i32, y: i32) -> i32 {
                    x - y
                }
                
                fn add(x: String, y: i32) -> String {
                    x + &y.to_string()
                }`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "i32 i32 i32") {
		t.Errorf(
			"TestRSGetMultipleFunctionSignature failed want %q got %v",
			"i32 i32 i32",
			slice,
		)
	}
	if !slices.Contains(slice, "String i32 String") {
		t.Errorf(
			"TestRSGetMultipleFunctionSignature failed want %q got %v",
			"String i32 String",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestRSGetMultipleFunctionSignature failed length wrong")
	}
}

func TestRSGetNoFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := ``
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if len(slice) != 0 {
		t.Errorf("TestRSGetNoFunctionSignature failed length wrong")
	}
}

func TestRSGetEmptyFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `fn test() {}`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if len(slice) != 1 {
		t.Errorf("TestRSGetEmptyFunctionSignature failed length wrong")
	}
	if !slices.Contains(slice, "void") {
		t.Errorf(
			"TestRSGetEmptyFunctionSignature failed want %q got %v",
			"void",
			slice,
		)
	}
}

func TestRSGetNestedFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `fn nested(x: i32) {
                    fn nested2(y: i32) {
                        fn nested3(z: i32) -> String {
                            "test".to_string()
                        }
                        nested3(y);
                    }
                    nested2(x);
                }`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "i32 void") {
		t.Errorf(
			"TestRSGetNestedFunctionSignature failed want %q got %v",
			"i32 void",
			slice,
		)
	}
	if !slices.Contains(slice, "i32 String") {
		t.Errorf(
			"TestRSGetNestedFunctionSignature failed want %q got %v",
			"i32 String",
			slice,
		)
	}
	if len(slice) != 3 {
		t.Errorf("TestRSGetNestedFunctionSignature failed length wrong")
	}
}

func TestRSGetClass(t *testing.T) {
	InitsrcToText(Test())
	testcode := `struct Person {
					name: String,
					age: i32,
				}
				
				impl Person {
					fn new(name: String, age: i32) -> Person {
						Person { name, age }
					}
				}`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "Person String i32 Person") {
		t.Errorf(
			"TestRSGetClass failed want %q got %v",
			"Person String i32 Person",
			slice,
		)
	}
}

func TestRSGetReceiver(t *testing.T) {
	InitsrcToText(Test())
	testcode := `struct Person {
					name: String,
					age: i32,
				}
				
				impl Person {
					fn new(&mut self, age: i32) -> Person {
						Person { name, age }
					}
				}`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "Person &mut self i32 Person") {
		t.Errorf(
			"TestRSGetReceiver failed want %q got %v",
			"Person &mut self i32 Person",
			slice,
		)
	}
}

func TestRSGetDocs(t *testing.T) {
	InitsrcToText(Test())
	testcode := `/* Adds two numbers.

# Arguments
 
* x - An integer.
* y - Another integer.
 
# Returns
 
The sum of x and y.
*/
fn add(x: i32, y: i32) -> i32 {
    x + y
}`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindDocs([]byte(testcode), parser, conf)
	if !slices.Contains(slice, `/* Adds two numbers.

# Arguments
 
* x - An integer.
* y - Another integer.
 
# Returns
 
The sum of x and y.
*/`) {
		t.Errorf(
			"TestRSGetDocs failed want %q got %v",
			"/* Adds two numbers...",
			slice,
		)
	}
	if len(slice) != 1 {
		t.Errorf("TestRSGetDocs failed length wrong")
	}
}

func TestRSGetLineDocs(t *testing.T) {
	InitsrcToText(Test())
	testcode := `/// Adds two numbers.
///
/// * x - An integer.
/// * y - Another integer.
fn add(x: i32, y: i32) -> i32 {
    x + y
}`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindDocs([]byte(testcode), parser, conf)
	if slice[0] !=
		`/// Adds two numbers.
///
/// * x - An integer.
/// * y - Another integer.
` {
		t.Errorf(
			"TestRSGetLineDocs failed want %q got %v",
			"/// Adds two numbers...",
			slice,
		)
	}
	if len(slice) != 1 {
		t.Errorf("TestRSGetLineDocs failed length wrong")
	}
}

func TestRSGetMultipleDocs(t *testing.T) {
	InitsrcToText(Test())
	testcode := `/// Adds two numbers.
/// 
/// * x - An integer.
/// * y - Another integer.
fn add(x: i32, y: i32) -> i32 {
    x + y
}

/// Subtracts two numbers.
/// 
/// * x - An integer.
/// * y - Another integer.
fn sub(x: i32, y: i32) -> i32 {
    x - y
}`
	parser := tree_sitter.NewParser()
	conf := config[".rs"]
	lang, _ := getLanguageFileExt(".rs", conf)
	parser.SetLanguage(lang)
	slice, _ := FindDocs([]byte(testcode), parser, conf)
	if !slices.Contains(slice,
		`/// Adds two numbers.
/// 
/// * x - An integer.
/// * y - Another integer.
`) {
		t.Errorf(
			"TestRSGetMultipleDocs failed want %q got %v",
			"/// Adds two numbers...",
			slice,
		)
	}
	if !slices.Contains(slice, `/// Subtracts two numbers.
/// 
/// * x - An integer.
/// * y - Another integer.
`) {
		t.Errorf(
			"TestRSGetMultipleDocs failed want %q got %v",
			"/// Subtracts two numbers...",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestRSGetMultipleDocs failed length wrong")
	}
}
