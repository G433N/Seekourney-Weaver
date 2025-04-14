# Naming convention

This is how to name constants, variables, functions, types and strucs

## Private constants

All caps with an underscore at the start.   
Example: `_TESTCONST_`

## Anything else

Snake case with first letter upper case if it exported, Snake case with first letter lower case otherwise.   
Example: `TestExportFunction`, `testPrivateVariable`

# Documentation

Documentation follows https://tip.golang.org/doc/comment   
Function parameter and return types should not need be mentioned in function documentation. The reason for this is that all documentation for the types should be at the type definition.   
You can use either `/* */` or `//`, but it is preferred to use `//` when the documentation is three lines or less.   
Example:   
```
/*
Copy copies from src to dst until either EOF is reached
on src or an error occurs. It returns the total number of bytes
written and the first error encountered while copying, if any.

A successful Copy returns err == nil, not err == EOF.
Because Copy is defined to read from src until EOF, it does
not treat an EOF from Read as an error to be reported.
*/
func Copy(dst Writer, src Reader) (n int64, err error) {
    ...
}

// Sort sorts data in ascending order as determined by the Less method.
// It makes one call to data.Len to determine n and O(n*log(n)) calls to
// data.Less and data.Swap. The sort is not guaranteed to be stable.
func Sort(data Interface) {
    ...
}
```

## Packages

Right now it is not necessary to write documentation for packages. Feel free to add it if you want.
