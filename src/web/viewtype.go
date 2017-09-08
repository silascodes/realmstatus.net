package main

type ViewType interface {
	GetValue() string
	GetDisplay() string
}