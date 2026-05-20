package must

func Do[OUT any](o OUT, err error) OUT {
	if err != nil {
		panic(err)
	}
	return o
}
