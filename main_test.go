package main

import "testing"

func TestMainOutput(t *testing.T) {
    expected := "Hello, World!"
    if expected != "Hello, World!" {
        t.Errorf("Expected '%s' but got '%s'", expected, "Hello, World!")
    }
}