//nolint:all
package srctotxt

func Default() ExtensionMap {
	return ExtensionMap{
		".go": {
			grammarPath: "./srctotxt/grammarlibs/libtree-sitter-go.so",
			libFunc:     "tree_sitter_go",
			parameters:  "parameter_list",
			functionDeclaration: []string{
				"function_declaration",
				"func_literal",
				"method_declaration",
			},
			classDeclaration:   []string{""},
			rightSideTypeParam: true,
			rightSideReturn:    true,
			receiver:           "receiver",
			returnType:         []string{"identifier", "parameter_list"},
			blockComment:       "comment",
			lineComment:        "comment",
		},
		".ts": {
			grammarPath: "./srctotxt/grammarlibs/libtree-sitter-typescript.so",
			libFunc:     "tree_sitter_typescript",
			parameters:  "formal_parameters",
			functionDeclaration: []string{
				"function_declaration",
				"method_definition",
			},
			classDeclaration:   []string{"class_declaration"},
			rightSideTypeParam: true,
			rightSideReturn:    true,
			receiver:           "",
			returnType: []string{
				"formal_parameters",
				"type_annotation",
			},
			blockComment: "comment",
			lineComment:  "comment",
		},
		".cs": {
			grammarPath: "./srctotxt/grammarlibs/libtree-sitter-c-sharp.so",
			libFunc:     "tree_sitter_c_sharp",
			parameters:  "parameter_list",
			functionDeclaration: []string{"method_declaration",
				"constructor_declaration",
				"operator_declaration",
				"conversion_operator_declaration",
				"destructor_declaration",
				"local_function_statement",
			},
			classDeclaration:   []string{"class_declaration"},
			rightSideTypeParam: false,
			rightSideReturn:    false,
			receiver:           "",
			returnType: []string{
				"identifier",
				"generic_name",
				"predefined_type",
				"variable_declaration",
			},
			blockComment: "comment",
			lineComment:  "comment",
		},
		".rs": {
			grammarPath:         "./srctotxt/grammarlibs/libtree-sitter-rust.so",
			libFunc:             "tree_sitter_rust",
			parameters:          "parameters",
			functionDeclaration: []string{"function_item"},
			classDeclaration:    []string{"impl_item"},
			rightSideTypeParam:  true,
			rightSideReturn:     true,
			receiver:            "self_parameter",
			returnType: []string{
				"generic_type",
				"primitive_type",
				"identifier",
				"reference_type",
				"where_clause",
			},
			blockComment: "block_comment",
			lineComment:  "line_comment",
		},
	}
}

func Test() ExtensionMap {
	return ExtensionMap{
		".go": {
			grammarPath: "./grammarlibs/libtree-sitter-go.so",
			libFunc:     "tree_sitter_go",
			parameters:  "parameter_list",
			functionDeclaration: []string{
				"function_declaration",
				"func_literal",
				"method_declaration",
			},
			classDeclaration:   []string{""},
			rightSideTypeParam: true,
			rightSideReturn:    true,
			receiver:           "receiver",
			returnType:         []string{"identifier", "parameter_list"},
			blockComment:       "comment",
			lineComment:        "comment",
		},
		".ts": {
			grammarPath: "./grammarlibs/libtree-sitter-typescript.so",
			libFunc:     "tree_sitter_typescript",
			parameters:  "formal_parameters",
			functionDeclaration: []string{
				"function_declaration",
				"method_definition",
			},
			classDeclaration:   []string{"class_declaration"},
			rightSideTypeParam: true,
			rightSideReturn:    true,
			receiver:           "",
			returnType: []string{
				"formal_parameters",
				"type_annotation",
			},
			blockComment: "comment",
			lineComment:  "comment",
		},
		".cs": {
			grammarPath: "./grammarlibs/libtree-sitter-c-sharp.so",
			libFunc:     "tree_sitter_c_sharp",
			parameters:  "parameter_list",
			functionDeclaration: []string{"method_declaration",
				"constructor_declaration",
				"operator_declaration",
				"conversion_operator_declaration",
				"destructor_declaration",
				"local_function_statement",
			},
			classDeclaration:   []string{"class_declaration"},
			rightSideTypeParam: false,
			rightSideReturn:    false,
			receiver:           "",
			returnType: []string{
				"identifier",
				"generic_name",
				"predefined_type",
				"variable_declaration",
			},
			blockComment: "comment",
			lineComment:  "comment",
		},
		".rs": {
			grammarPath:         "./grammarlibs/libtree-sitter-rust.so",
			libFunc:             "tree_sitter_rust",
			parameters:          "parameters",
			functionDeclaration: []string{"function_item"},
			classDeclaration:    []string{"impl_item"},
			rightSideTypeParam:  true,
			rightSideReturn:     true,
			receiver:            "self_parameter",
			returnType: []string{
				"generic_type",
				"primitive_type",
				"identifier",
				"reference_type",
				"where_clause",
			},
			blockComment: "block_comment",
			lineComment:  "line_comment",
		},
	}
}

//config
/*
func Default() Config {
	return Config{
		".go": LanguageConfig{
			typeConf: TypeConfig {str: "string"},
			functionConf: FunctionConfig {
				functionDeclaration: "func",
				continueAfterDeclare: false,
				parameterDelimiter: ",",
				parameterEnclosure: "(",
			},
			scopeDelimStart: "{",
			scopeDelimEnd: "}",
			commentConf: CommentConfig{
				comment: "//",
				beginComment: "/*",
				endComment: "*/ /*",
	},
	terminator: "\n",
},
".ts": LanguageConfig{
	typeConf: TypeConfig {str: "string"},
	functionConf: FunctionConfig {
		functionDeclaration: "function",
		continueAfterDeclare: false,
		parameterDelimiter: ",",
		parameterEnclosure: "(",
	},
	scopeDelimStart: "{",
	scopeDelimEnd: "}",
	commentConf: CommentConfig{
		comment: "//",
		beginComment: "/*",
		endComment: "*/ /*",
			},
			terminator: "\n",
		},
		"ex": LanguageConfig{
			typeConf: TypeConfig{str: "string"},
			functionConf: FunctionConfig{
				functionDeclaration: "def",
				continueAfterDeclare: false,
				parameterDelimiter: ",",
				parameterEnclosure: "(",
			},
			scopeDelimStart: "do",
			scopeDelimEnd: "end",
			commentConf: CommentConfig{
				comment: "#",
				beginComment: "@doc\"\"\"",
				endComment: "\"\"\"",
			},
			terminator: "\n",
		},
	}
}
*/
