# Naming convention

This is how to name constants, variables, functions, types and strucs

## Private constants

All caps with underscores of both side of the name.   
Example: `_TESTCONST_`

## Anything else

Pascal case if it is exported, camel case otherwise.   
Example: `ExampleExportFunction`, `examplePrivateVariable`

# Documentation

Documentation follows https://tip.golang.org/doc/comment   
Function parameter and return types should not need be mentioned in function documentation.
The reason for this is that all documentation for the types should be at the type definition.   
You can use either `/* */` or `//`, but it is preferred to use `//`
when the documentation is three lines or less.   
Example:   
```go
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

## Types

As previously mentioned, types should be documented, including any invariant, range, etc.
Types should usually not be aliases.   
Types should be added if the intended use of the currently used type
is not obvious or there is some not obvious restriction on the type of values.    
For example, if you have a string that represents a date,
any string could be seen as a valid date.
If you create a new type "Date", the differnce is clearer.

## Testing

Test function style: https://pkg.go.dev/testing   
The assert library should be used.

## Packages

Right now it is not necessary to write documentation for packages.
Feel free to add it if you want.

# Formatting

For backend code, use golangci-lint.   
For frontend code, use Prettier.
