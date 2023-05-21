package main

func empty() {
	return
}

func unnamed() (int, int) {
	return 11, 22
}

func partial() (int, b int) { // want `functions with partially named return variables not allowed`
	return 11, 22
}

func normal() (a int, b int) {
	return 11, 22
}

func naked() (a int, b int) {
	a = 11
	b = 22
	return // want `naked returns not allowed`
}

func mixing() (a int, b int) {
	if true {
		return 11, 22
	} else {
		a = 11
		b = 22
		return // want `naked returns not allowed`
	}
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

	// partial()
	func() (int, b int) { // want `functions with partially named return variables not allowed`
		return 11, 22
	}()

	// normal()
	func() (a int, b int) {
		return 11, 22
	}()

	// naked()
	func() (a int, b int) {
		a = 11
		b = 22
		return // want `naked returns not allowed`
	}()

	// mixing()
	func() (a int, b int) {
		if true {
			return 11, 22
		} else {
			a = 11
			b = 22
			return // want `naked returns not allowed`
		}
	}()
}
