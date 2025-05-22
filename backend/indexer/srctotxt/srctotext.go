package srctotxt

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"seekourney/utils"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/tree-sitter/go-tree-sitter"
)

type FileExtension string
type Function string

type TreeSitterConf struct {
	grammarPath utils.Path
	libFunc string
	parameters string
	functionDeclaration []string
	classDeclaration []string
	rightSideTypeParam bool
	rightSideReturn bool
	receiver string
	returnType []string
	blockComment string
	lineComment string
}

type ExtensionMap map[FileExtension]TreeSitterConf

var config ExtensionMap

func initsrcToText(newConfig ExtensionMap){
	config = newConfig
}

func getLanguage(path utils.Path, conf TreeSitterConf)(*tree_sitter.Language, error){
	fileExtension := filepath.Ext(string(path))
	pathSO := config[FileExtension(fileExtension)].grammarPath

	lib, err := purego.Dlopen(string(pathSO), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
        return nil, err
    }

	var language func() uintptr
	purego.RegisterLibFunc(&language, lib, conf.libFunc)
	sitterLanguage := tree_sitter.NewLanguage(unsafe.Pointer(language()))
	if(sitterLanguage == nil){
		return nil, errors.New("tree sitter language not found")
	}
	return sitterLanguage, nil
}

func getLanguageFileExt(fileExtension FileExtension, conf TreeSitterConf)(*tree_sitter.Language, error){
	pathSO := config[fileExtension].grammarPath
	lib, err := purego.Dlopen(string(pathSO), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
        return nil, err
    }
	var language func() uintptr
	purego.RegisterLibFunc(&language, lib, conf.libFunc)
	sitterLanguage := tree_sitter.NewLanguage(unsafe.Pointer(language()))
	if(sitterLanguage == nil){
		return nil, errors.New("tree sitter language not found")
	}
	return sitterLanguage, nil
}

func contains(str string, lst []string) bool{
	for i := 0; i < len(lst); i++{
		if(str == lst[i]){
			return true
		}
	}
	return false
}


func findFuncs(sourceCode[]byte, parser *tree_sitter.Parser, conf TreeSitterConf) ([]Function, error){
    defer parser.Close()
    tree := parser.Parse(sourceCode, nil)
    defer tree.Close()

	rootNode := tree.RootNode()
	var funcs []Function

	var findFuncsHelper func(node tree_sitter.Node) error
	findFuncsHelper = func(node tree_sitter.Node) error{
		if(node.NamedChildCount() == 0){
			return nil
		}
		if (contains(node.GrammarName(), conf.functionDeclaration)) { //check if we have a function declaration
			if(node.ChildCount() < 2){
				return errors.New("function declaration invalid")
			}
			functionSignature := sourceCode[node.StartByte():node.Child(node.ChildCount() - 2).EndByte()] //take out body by checking all the children - 2 (since last child node is usually function body)
			funcs = append(funcs, Function(functionSignature))
		}
		for i := uint(0); i < node.NamedChildCount(); i++ {
			findFuncsHelper(*node.NamedChild(i))
		}
		return nil
	}
	if(rootNode == nil || rootNode == &tree_sitter.Node{}){
		return funcs, nil
	}
	err := findFuncsHelper(*rootNode)
	if(err != nil){
		return nil, err
	}
	return funcs, nil
}

//findClass
//Helper function that finds the class of a given function
func findClass(node *tree_sitter.Node, sourceCode []byte, conf TreeSitterConf) string{
	for node.Parent() != nil {
		node = node.Parent()
		if(contains(node.GrammarName(), conf.classDeclaration)){
			for i := uint(0); node.NamedChildCount() > i; i++ {
				if(node.NamedChild(i).GrammarName() == "identifier"){
					node = node.NamedChild(i)
					break
				}
			} 
			return string(sourceCode[node.StartByte():node.EndByte()]) + " "
		}
	}
	return ""
}

//getSrcCode
//gets the sourcecode of a file from a given path
func getSrcCode(path utils.Path) ([]byte, error){
	sourceCode, err := os.ReadFile(string(path))
	if err != nil{
		return nil, err
	}
	return sourceCode, nil
}

//toTree
//Gets a parser and config for a given file on a file path
//errors if the file is not found or the language is not supported
func toTree(path utils.Path)(*tree_sitter.Parser, TreeSitterConf, error){
	currentLang := FileExtension(filepath.Ext(string(path)))
	conf, exists := config[currentLang]
	if(!exists){
		return nil, conf, errors.New("language not found")
	}
	parser := tree_sitter.NewParser()
	language, err := getLanguage(path, conf)
	if err != nil {
		return nil, conf, err
	}
    err = parser.SetLanguage(language)
	if( err != nil) {
		return nil, conf, err
	}
	return parser, conf, nil
}

//findFuncSignature
//Extracts all function signatures from a given sourcecode
func findFuncSignature(sourceCode[]byte, parser *tree_sitter.Parser, conf TreeSitterConf)([]Function, error){
    defer parser.Close()
    tree := parser.Parse(sourceCode, nil)
    defer tree.Close()

	rootNode := tree.RootNode()
	var funcs []Function

	var findFuncsHelper func(node tree_sitter.Node) error
	findFuncsHelper = func(node tree_sitter.Node) error{
		var findParameters func(node tree_sitter.Node, paramIndex uint) (string, uint, error)
		findParameters = func(node tree_sitter.Node, paramIndex uint) (string, uint, error){
			if node.NextNamedSibling() == nil{
				return "", paramIndex, nil
			}
			if node.GrammarName() == conf.parameters{
				var parameters string
				if(conf.rightSideTypeParam){
					for i := uint(0); i < node.NamedChildCount(); i++ {
						if(node.NamedChild(i).ChildCount() < 2){
							continue
						}
						paramType := node.NamedChild(i).NamedChild(1)
						if(paramType == nil){
							paramType = node.NamedChild(i).NamedChild(0)
						}
						if(node.NamedChild(i).GrammarName() == conf.receiver){
							if(len(sourceCode) < int(node.NamedChild(i).Child(0).EndByte())){
								return "", paramIndex, errors.New("source code smaller than node byte")
							}
							parameters += string(sourceCode[node.NamedChild(i).Child(0).StartByte():paramType.EndByte()]) + " "
						}else {
							if(len(sourceCode) < int(paramType.EndByte())){
								return "", paramIndex, errors.New("source code smaller than node byte")
							}
							parameters += string(sourceCode[paramType.StartByte():paramType.EndByte()]) + " "
						}
					}
				}else{
					for i := uint(0); i < node.NamedChildCount(); i++ {
						if(node.NamedChild(i).ChildCount() < 1){
							continue
						}
						typeNode := node.NamedChild(i).NamedChild(0)
						if(len(sourceCode) < int(typeNode.EndByte())){
							return "", paramIndex, errors.New("source code smaller than node byte")
						}
						parameters +=  string(sourceCode[typeNode.StartByte():typeNode.EndByte()]) + " "
					}
				}
				return parameters, paramIndex, nil
			}
			return findParameters(*node.NextNamedSibling(), paramIndex + 1)
		}
		if(node.NamedChildCount() == 0){
			return nil
		}
		if (contains(node.GrammarName(), conf.functionDeclaration)) { //check if we have a function declaration
			currentSignature, paramIndex, err := findParameters(*node.Child(0), 0) //params + where the parameters begin
			if(err != nil){
				return err
			}
			currentSignature = findClass(&node, sourceCode, conf) + currentSignature
			if(conf.rightSideReturn){
				nodes := node.NamedChildCount() - 2
				if(paramIndex > nodes){ //compare to see that return and parameters arent the same
					currentSignature += "void"
				} else if(contains(node.NamedChild(nodes).GrammarName(), conf.returnType)){
					if(len(sourceCode) < int(node.NamedChild(nodes).EndByte())){
						return errors.New("source code smaller than node byte")
					}
					currentSignature += string(sourceCode[node.NamedChild(nodes).StartByte():node.NamedChild(nodes).EndByte()])
				} else {
					currentSignature += "void"
				}
			} else {
				for i:= uint(0);i<node.NamedChildCount();i++{
					if(contains(node.NamedChild(i).GrammarName(), conf.returnType)){
						currentSignature += string(sourceCode[node.NamedChild(i).StartByte():node.NamedChild(i).EndByte()])
						break
					}
				}
			}
			funcs = append(funcs, Function(currentSignature))
		}
		for i := uint(0); i < node.NamedChildCount(); i++ {
			findFuncsHelper(*node.NamedChild(i))
		}
		return nil
	}
	if(rootNode == nil || rootNode == &tree_sitter.Node{}){
		return funcs, nil
	}
	err := findFuncsHelper(*rootNode)
	if(err != nil){
		return nil, err
	}

	return funcs, nil
}

func isNested(node *tree_sitter.Node, searchFor []string, searchLimit int, childLimit int)bool{
	if(searchLimit < 1){
		return false
	}
	childrenAmount := node.NamedChildCount()
	if(childrenAmount>uint(childLimit)){
		return false
	}
	for i := uint(0); i<childrenAmount;i++{
		if contains(node.GrammarName(), searchFor){
			return true
		}
		if isNested(node.NamedChild(i), searchFor, searchLimit-1, childLimit) {
            return true
        }
	}
	return false
}

func findDocs(sourceCode[]byte, parser *tree_sitter.Parser, conf TreeSitterConf) ([]string, error){
    defer parser.Close()
    tree := parser.Parse(sourceCode, nil)
    defer tree.Close()

	rootNode := tree.RootNode()
	var docs []string

	var findDocsHelper func(node tree_sitter.Node, acc string) (string, error)
	findDocsHelper = func(node tree_sitter.Node, acc string) (string, error){
		if(node.GrammarName() == conf.blockComment || node.GrammarName() == conf.lineComment){
			if(node.NextNamedSibling() != nil){
				if(acc != ""){
					if(isNested(node.NextNamedSibling(), conf.functionDeclaration, 3, 5)){
						docs = append(docs, acc)
						return "", nil
					} else if node.NextSibling().GrammarName() != conf.lineComment{
						return "", nil
					}
					if(len(sourceCode) < int(node.NextNamedSibling().EndByte())){ //this has never happened in my tests, but just in case, this is an issue with tree-sitter if it happens.
						return "", errors.New("source code smaller than node byte")
					}
					return acc + string(sourceCode[node.NextNamedSibling().StartByte():node.NextNamedSibling().EndByte()]), nil
				}
				if(node.NextNamedSibling().GrammarName() == conf.lineComment){
					if(len(sourceCode) < int(node.NextNamedSibling().EndByte())){
						return "", errors.New("source code smaller than node byte")
					}
					return acc + string(sourceCode[node.StartByte():node.NextNamedSibling().EndByte()]), nil
				} else if isNested(node.NextNamedSibling(), conf.functionDeclaration, 5, 10){ //five is how many layers deep our search is allowed to go, 10 is how many children is max before we consider it not a functon declaration
					if(len(sourceCode) < int(node.EndByte())){
						return "", errors.New("source code smaller than node byte")
					}
					docs = append(docs, string(sourceCode[node.StartByte():node.EndByte()]))
				}
			}
		}
		for i := uint(0); i < node.ChildCount(); i++ {
			acc,_ = findDocsHelper(*node.Child(i), acc)
		}
		return acc, nil
	}
	var acc string
	findDocsHelper(*rootNode, acc)
	return docs, nil
}




func Run(){
	initsrcToText(Default())
	parser, conf,_ := toTree("./srctotxt/testfiles/test1.go")
	sourceCode, _ := getSrcCode("./srctotxt/testfiles/test1.go")
	fmt.Println(findFuncSignature(sourceCode, parser, conf))
}