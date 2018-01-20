package examples

// PanicOnErr panics when err is not nil
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
