package assertions

func AssertOk(err error) {
	if err != nil {
		panic(err)
	}
}
