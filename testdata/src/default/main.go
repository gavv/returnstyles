package main

/*
void add(){
	int a = 11 + 22;
	return;
}
*/
import "C"

func empty() {
	return
}

func unnamed() (int, int) {
	return 11, 22
}

func namedPartial() (int, b int) { // want `functions with partially named return variables not allowed`
	return 11, 22
}

func namedNormal() (a int, b int) {
	return 11, 22
}

func namedNaked() (a int, b int) {
	a = 11
	b = 22
	return
}

func namedMixing1() (a int, b int) {
	if true {
		return 11, 22
	} else {
		a = 11
		b = 22
		return // want `mixing normal and naked returns not allowed`
	}
}

func namedMixing2() (a int, b int) {
	if true {
		a = 11
		b = 22
		return
	} else {
		return 11, 22 // want `mixing normal and naked returns not allowed`
	}
}

func main() {
	C.add()
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
		return 11, 22
	}()

	// namedNormal()
	func() (a int, b int) {
		return 11, 22
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
			return 11, 22
		} else {
			a = 11
			b = 22
			return // want `mixing normal and naked returns not allowed`
		}
	}()

	// namedMixing2()
	func() (a int, b int) {
		if true {
			a = 11
			b = 22
			return
		} else {
			return 11, 22 // want `mixing normal and naked returns not allowed`
		}
	}()
}
