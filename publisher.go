package main

type Publisher interface {
	Publish(msg string) error
}
