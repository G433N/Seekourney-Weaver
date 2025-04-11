package main

func Bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

func Italic(text string) string {
	return "\033[3m" + text + "\033[0m"
}

func LightBlue(text string) string {
	return "\033[94m" + text + "\033[0m"
}

func Green(text string) string {
	return "\033[92m" + text + "\033[0m"
}
