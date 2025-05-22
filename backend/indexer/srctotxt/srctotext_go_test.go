package srctotxt

import "testing"
import "github.com/tree-sitter/go-tree-sitter"
import "slices"

func TestGetFunction(t *testing.T) {
	initsrcToText(Test())
	testcode := `func sub(x int, y int) int{
					return x-y
				}`
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncs([]byte(testcode), parser, conf)
	if(string(slice[0]) != "func sub(x int, y int) int"){
		t.Errorf("TestGetFunction failed want %q got %v","func sub(x int, y int) int", slice[0])
	}
}

func TestGetMultipleFunction(t *testing.T){
	initsrcToText(Test())
	testcode := `func sub(x int, y int) int{
					return x-y
				}
				
				func add(x string, y int) (uint, int){
					return uint(y), y
				}`
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncs([]byte(testcode), parser, conf)
	if(!slices.Contains(slice, "func sub(x int, y int) int")){
		t.Errorf("TestGetMultipleFunction failed want %q got %v","func sub(x int, y int) int", slice)
	}
	if(!slices.Contains(slice, "func add(x string, y int) (uint, int)")){
		t.Errorf("TestGetMultipleFunction failed want %q got %v","func add(x string, y int) (uint, int)", slice)
	}
	if(len(slice) != 2){
		t.Errorf("TestGetMultipleFunction failed length wrong")
	}
}

func TestGetNoFunction(t *testing.T){
	initsrcToText(Test())
	testcode := ``
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncs([]byte(testcode), parser, conf)
	if(len(slice) != 0){
		t.Errorf("TestGetNoFunction failed length wrong")
	}
}

func TestGetEmptyFunction(t *testing.T){
	initsrcToText(Test())
	testcode := `func test(){}`
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncs([]byte(testcode), parser, conf)
	if(len(slice) != 1){
		t.Errorf("TestGeEmptyFunction failed length wrong")
	}
	if(!slices.Contains(slice, "func test()")){
		t.Errorf("TestGetEmptyFunction failed want %q got %v","func test()", slice)
	}
}

func TestGetNestedFunction(t *testing.T){
	initsrcToText(Test())
	testcode := `func nested(x int){
					var nested2 func(y int)
					nested2 = func(y int){
						return
					}
					nested2(x)
				}`
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncs([]byte(testcode), parser, conf)
	if(!slices.Contains(slice, "func nested(x int)")){
		t.Errorf("TestGetMultipleFunction failed want %q got %v","func nested(x int)", slice)
	}
	if(!slices.Contains(slice, "func(y int)")){
		t.Errorf("TestGetMultipleFunction failed want %q got %v","func(y int)", slice)
	}
	if(len(slice) != 2){
		t.Errorf("TestGetMultipleFunction failed length wrong")
	}
}

func TestGetFunctionSignature(t *testing.T){
	initsrcToText(Test())
	testcode := `func sub(x int, y int) int{
					return x-y
				}`
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncSignature([]byte(testcode), parser, conf)
	if(string(slice[0]) != "int int int"){
		t.Errorf("TestGetFunction failed want %q got %v","int int int", slice[0])
	}
}

func TestGetMultipleFunctionSignature(t *testing.T){
	initsrcToText(Test())
	testcode := `func sub(x int, y int) int{
					return x-y
				}
				
				func add(x string, y int) (uint, int){
					return uint(y), y
				}`
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncSignature([]byte(testcode), parser, conf)
	if(!slices.Contains(slice, "int int int")){
		t.Errorf("TestGetMultipleFunction failed want %q got %v","int int int", slice)
	}
	if(!slices.Contains(slice, "string int (uint, int)")){
		t.Errorf("TestGetMultipleFunction failed want %q got %v","string int (uint, int)", slice)
	}
	if(len(slice) != 2){
		t.Errorf("TestGetMultipleFunction failed length wrong")
	}
}

func TestGetNoFunctionSignature(t *testing.T){
	initsrcToText(Test())
	testcode := ``
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncSignature([]byte(testcode), parser, conf)
	if(len(slice) != 0){
		t.Errorf("TestGetNoFunctionSignature failed length wrong")
	}
}

func TestGetEmptyFunctionSignature(t *testing.T){
	initsrcToText(Test())
	testcode := `func test(){}`
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncSignature([]byte(testcode), parser, conf)
	if(len(slice) != 1){
		t.Errorf("TestGeEmptyFunction failed length wrong")
	}
	if(!slices.Contains(slice, "void")){
		t.Errorf("TestGetEmptyFunction failed want %q got %v","void", slice)
	}
}

func TestGetNestedFunctionSignature(t *testing.T){
	initsrcToText(Test())
	testcode := `func nested(x int){
					var nested2 func(y int)
					nested2 = func(y int){
						var nested3 func(z int)string
						nested3 = func(z int)string{
							return "test"
						}
						nested3(y)
					}
					nested2(x)
				}`
	parser := tree_sitter.NewParser()
	conf := config[".go"]
	lang, _ := getLanguageFileExt(".go", conf)
    parser.SetLanguage(lang)
    slice,_ := findFuncSignature([]byte(testcode), parser, conf)
	if(!slices.Contains(slice, "int void")){
		t.Errorf("TestGetMultipleFunction failed want %q got %v","int void", slice)
	}
	if(!slices.Contains(slice, "int string")){
		t.Errorf("TestGetMultipleFunction failed want %q got %v","int string", slice)
	}
	if(len(slice) != 3){
		t.Errorf("TestGetMultipleFunction failed length wrong")
	}
}

func TestGetFunctionReceiver(t *testing.T){
	//cooked
}

func TestGetDocs(t *testing.T){

}

func TestGetMultipleDocs(t *testing.T){

}

func TestGetDocsComments(t *testing.T){

}




