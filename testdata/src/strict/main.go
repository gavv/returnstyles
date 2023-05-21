package main

func empty() {
	return
}

func unnamed() (int, int) {
	return 11, 22
}

func namedPartial() (int, b int) { // want `functions with partially named return variables not allowed`
	return 11, 22 // want `normal \(non-naked\) returns not allowed in functions with named return variables`
}

func namedNormal() (a int, b int) {
	return 11, 22 // want `normal \(non-naked\) returns not allowed in functions with named return variables`
}

func namedNaked() (a int, b int) {
	a = 11
	b = 22
	return
}

func namedMixing1() (a int, b int) {
	if true {
		return 11, 22 // want `normal \(non-naked\) returns not allowed in functions with named return variables`
	} else {
		a = 11
		b = 22
		return
	}
}

func namedMixing2() (a int, b int) {
	if true {
		a = 11
		b = 22
		return
	} else {
		return 11, 22 // want `normal \(non-naked\) returns not allowed in functions with named return variables`
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

	// namedPartial()
	func() (int, b int) { // want `functions with partially named return variables not allowed`
		return 11, 22 // want `normal \(non-naked\) returns not allowed in functions with named return variables`
	}()

	// namedNormal()
	func() (a int, b int) {
		return 11, 22 // want `normal \(non-naked\) returns not allowed in functions with named return variables`
	}()

	// namedNaked()
	func() (a int, b int) {
		a = 11
		b = 22
		return
	}()

	// namedMixing1()
	func() (a int, b int) {
		if true {
			return 11, 22 // want `normal \(non-naked\) returns not allowed in functions with named return variables`
		} else {
			a = 11
			b = 22
			return
		}
	}()

	// namedMixing2()
	func() (a int, b int) {
		if true {
			a = 11
			b = 22
			return
		} else {
			return 11, 22 // want `normal \(non-naked\) returns not allowed in functions with named return variables`
		}
	}()
}
