package main

func empty() {
	return
}

func unnamed() (int, int) { // want `functions with unnamed return variables not allowed`
	return 11, 22
}

func namedPartial() (int, b int) {
	return 11, 22
}

func namedNormal() (a int, b int) {
	return 11, 22
}

func outer() {
	// empty()
	func() {
		return
	}()

	// unnamed()
	func() (int, int) { // want `functions with unnamed return variables not allowed`
		return 11, 22
	}()

	// namedPartial()
	func() (int, b int) {
		return 11, 22
	}()

	// namedNormal()
	func() (a int, b int) {
		return 11, 22
	}()
}
