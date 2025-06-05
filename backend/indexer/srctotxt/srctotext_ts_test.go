//nolint:all
package srctotxt

import "testing"
import "github.com/tree-sitter/go-tree-sitter"
import "slices"

//import "fmt"

func TestTSGetFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `function add(x: int, y:int):int{
					return x + y;
				}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if string(slice[0]) != "function add(x: int, y:int):int" {
		t.Errorf(
			"TestTSGetFunction failed want %q got %v",
			"function add(x: int, y:int):int",
			slice[0],
		)
	}
}

func TestTSGetMultipleFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `function sub(x: int, y: int): int{
					return x - y;
				}
				
				function add(x: string, y: int): any{
					return y;
				}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "function sub(x: int, y: int): int") {
		t.Errorf(
			"TestTSGetMultipleFunction failed want %q got %v",
			"function sub(x: int, y: int): int",
			slice,
		)
	}
	if !slices.Contains(slice, "function add(x: string, y: int): any") {
		t.Errorf(
			"TestTSGetMultipleFunction failed want %q got %v",
			"function add(x: string, y: int): any",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestTSGetMultipleFunction failed length wrong")
	}
}

func TestTSGetNoFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := ``
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if len(slice) != 0 {
		t.Errorf("TestTSGetNoFunction failed length wrong")
	}
}

func TestTSGetEmptyFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `function test(){}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if len(slice) != 1 {
		t.Errorf("TestTSGetEmptyFunction failed length wrong")
	}
	if !slices.Contains(slice, "function test()") {
		t.Errorf(
			"TestTSGetEmptyFunction failed want %q got %v",
			"function test()",
			slice,
		)
	}
}

func TestTSGetNestedFunction(t *testing.T) {
	InitsrcToText(Test())
	testcode := `function nested(x: int){
					function nested2(y: int):int{
						return y;
					}
					nested2(x);
				}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncs([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "function nested(x: int)") {
		t.Errorf(
			"TestTSGetNestedFunction failed want %q got %v",
			"function nested(x: int)",
			slice,
		)
	}
	if !slices.Contains(slice, "function nested2(y: int):int") {
		t.Errorf(
			"TestTSGetNestedFunction failed want %q got %v",
			"function nested2(y: int):int",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestTSGetNestedFunction failed length wrong")
	}
}

func TestTSGetFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `function sub(x: int, y: int):int{
					return x-y;
				}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if string(slice[0]) != ": int : int :int" {
		t.Errorf(
			"TestTSGetFunctionSignature failed want %q got %v",
			"int int int",
			slice[0],
		)
	}
}

func TestTSGetMultipleFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `function sub(x: int, y: int): int{
					return x-y;
				}
				
				function add(x: string, y: int): any{
					return y
				}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if !slices.Contains(slice, ": int : int : int") {
		t.Errorf(
			"TestTSGetMultipleFunctionSignature failed want %q got %v",
			"int int int",
			slice,
		)
	}
	if !slices.Contains(slice, ": string : int : any") {
		t.Errorf(
			"TestTSGetMultipleFunctionSignature failed want %q got %v",
			"string int any",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestTSGetMultipleFunctionSignature failed length wrong")
	}
}

func TestTSGetNoFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := ``
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if len(slice) != 0 {
		t.Errorf("TestTSGetNoFunctionSignature failed length wrong")
	}
}

func TestTSGetEmptyFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `function test(){}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if len(slice) != 1 {
		t.Errorf("TestTSGetEmptyFunctionSignature failed length wrong")
	}
	if !slices.Contains(slice, "void") {
		t.Errorf(
			"TestTSGetEmptyFunctionSignature failed want %q got %v",
			"void",
			slice,
		)
	}
}

func TestTSGetNestedFunctionSignature(t *testing.T) {
	InitsrcToText(Test())
	testcode := `function nested(x: int){
					function nested2(y: int){
						function nested3(z: int):string{
							return "test";
						}
						nested3(y);
					}
					nested2(x);
				}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if !slices.Contains(slice, ": int void") {
		t.Errorf(
			"TestTSGetNestedFunctionSignature failed want %q got %v",
			": int void",
			slice,
		)
	}
	if !slices.Contains(slice, ": int :string") {
		t.Errorf(
			"TestTSGetNestedFunctionSignature failed want %q got %v",
			": int :string",
			slice,
		)
	}
	if len(slice) != 3 {
		t.Errorf("TestTSGetNestedFunctionSignature failed length wrong")
	}
}

func TestTSGetClass(t *testing.T) {
	InitsrcToText(Test())
	testcode := `class Test{
					function test():string{
						return "hello";
					}
				}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindFuncSignature([]byte(testcode), parser, conf)
	if !slices.Contains(slice, "Test :string") {
		t.Errorf("TestTSGetClass failed want %q got %v", "Test :string", slice)
	}
	if len(slice) != 1 {
		t.Errorf("TestTSGetClass failed length wrong")
	}
}

func TestTSGetDocs(t *testing.T) {
	InitsrcToText(Test())
	testcode := `/*add
@param x an int
@returns an int
*/
				function add(x: int):int{
					return x;
				}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindDocs([]byte(testcode), parser, conf)
	if !slices.Contains(slice, `/*add
@param x an int
@returns an int
*/`) {
		t.Errorf(
			"TestTSGetDocs failed want %q got %v",
			"add @param x an int @returns an int",
			slice,
		)
	}
	if len(slice) != 1 {
		t.Errorf("TestTSGetDocs failed length wrong")
	}
}

func TestTSGetLineDocs(t *testing.T) {
	InitsrcToText(Test())
	testcode := `
//add
//@param x an int
//@returns an int
//
				 function add(x: int):int{
				 	 return x;
				 }`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindDocs([]byte(testcode), parser, conf)
	if slice[0] !=
		`//add
//@param x an int//@returns an int//` {
		t.Errorf(
			"TestTSGetLineDocs failed want %q got %v",
			"//add//@param x an int//@returns an int//",
			slice,
		)
	}
	if len(slice) != 1 {
		t.Errorf("TestTSGetLineDocs failed length wrong")
	}
}

func TestTSGetMultipleDocs(t *testing.T) {
	InitsrcToText(Test())
	testcode := `
//add
//@param x an int
//@returns an int
//
				function add(x: int):int{
					return x;
				}
					
/*sub
@param x an int
@returns an int
*/
				
				function sub(x: int):int{
					return x;
				}`
	parser := tree_sitter.NewParser()
	conf := config[".ts"]
	lang, _ := getLanguageFileExt(".ts", conf)
	parser.SetLanguage(lang)
	slice, _ := FindDocs([]byte(testcode), parser, conf)
	if !slices.Contains(slice,
		`//add
//@param x an int//@returns an int//`) {
		t.Errorf(
			"TestTSGetMultipleDocs failed want %q got %v",
			"add @param x an int @returns an int",
			slice,
		)
	}
	if !slices.Contains(slice,
		`/*sub
@param x an int
@returns an int
*/`) {
		t.Errorf(
			"TestTSGetMultipleDocs failed want %q got %v",
			"sub @param x an int @returns an int",
			slice,
		)
	}
	if len(slice) != 2 {
		t.Errorf("TestTSGetMultipleDocs failed length wrong")
	}
}
