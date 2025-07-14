package component

func SplitURI(uri string) []string {

	// "direct:foo?a=1" -> ["direct", "foo"]
	for i := 0; i < len(uri); i++ {
		if uri[i] == ':' {
			return []string{uri[:i], uri[i+1:]}
		}
	}

	return []string{uri}
}
