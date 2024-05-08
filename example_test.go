package osrelease_test

import (
	"fmt"

	"github.com/RaphaelPour/osrelease"
)

func ExampleNew() {
	version := osrelease.New(6, 5, 4)
	fmt.Printf("Major: %d\n", version.Major())
	fmt.Printf("Minor: %d\n", version.Minor())
	fmt.Printf("Patch: %d\n", version.Patch())
	fmt.Println(version)

	// Output:
	// Major: 6
	// Minor: 5
	// Patch: 4
	// 6.5.4
}

func ExampleNew_withSuffix() {
	version := osrelease.New(5, 4, 0, osrelease.WithSuffix("-200.fc39.x86_64"))
	fmt.Printf("Major: %d\n", version.Major())
	fmt.Printf("Minor: %d\n", version.Minor())
	fmt.Printf("Patch: %d\n", version.Patch())
	fmt.Printf("Suffix: %s\n", version.Suffix())
	fmt.Println(version)

	// Output:
	// Major: 5
	// Minor: 4
	// Patch: 0
	// Suffix: -200.fc39.x86_64
	// 5.4.0-200.fc39.x86_64
}

func ExampleParseString() {
	version, err := osrelease.ParseString("6.8.8-200.fc39.x86_64")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Major: %d\n", version.Major())
	fmt.Printf("Minor: %d\n", version.Minor())
	fmt.Printf("Patch: %d\n", version.Patch())
	fmt.Printf("Suffix: %s\n", version.Suffix())
	fmt.Println(version)
	// Output:
	// Major: 6
	// Minor: 8
	// Patch: 8
	// Suffix: -200.fc39.x86_64
	// 6.8.8-200.fc39.x86_64
}

func ExampleVersion_NewerThan() {
	oldVersion, err := osrelease.ParseString("1.2.3")
	if err != nil {
		fmt.Println(err)
		return
	}

	newVersion, err := osrelease.ParseString("2.0.0")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"%s newer than %s: %t\n",
		newVersion,
		oldVersion,
		newVersion.NewerThan(oldVersion),
	)

	// Output:
	// 2.0.0 newer than 1.2.3: true
}

func ExampleVersion_OlderThan() {
	oldVersion, err := osrelease.ParseString("1.1.5")
	if err != nil {
		fmt.Println(err)
		return
	}

	newVersion, err := osrelease.ParseString("1.2.1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"%s older than %s: %t\n",
		oldVersion,
		newVersion,
		oldVersion.OlderThan(newVersion),
	)

	// Output:
	// 1.1.5 older than 1.2.1: true
}
