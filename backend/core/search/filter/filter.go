package filter

type Filter interface{}

type Includes string

type Excludes string

func ParseFilter(filter string) Filter {
	return ""
}
