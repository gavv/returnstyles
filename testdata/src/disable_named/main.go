package main

func empty() {
	return
}

func unnamed() (int, int) {
	return 11, 22
}

func namedPartial() (int, b int) { // want `functions with named return variables not allowed`
	return 11, 22
}

func namedNormal() (a int, b int) { // want `functions with named return variables not allowed`
	return 11, 22
}

func outer() {
	// empty()
	func() {
		return
	}()

	// unnamed()
	func() (int, int) {
		return 11, 22
	}()

	// namedPartial()
	func() (int, b int) { // want `functions with named return variables not allowed`
		return 11, 22
	}()

	// namedNormal()
	func() (a int, b int) { // want `functions with named return variables not allowed`
		return 11, 22
	}()
}
