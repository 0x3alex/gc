package main

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	/*
		Create a new cache with a cleaing interval of 5 minutes and an buffer
		of 100 actions
	*/
	gc := New(5*time.Minute, 100)
	gc.Run()
	if gc == nil {
		t.Fatal("GC is nil")
	}
	t.Log("GC init ok!")
	/*
		Store the pair ("hello", "world") with a time to live of the cache cleaning interval
	*/
	gc.Set("hello", "world", gc.DefaultTTL)
	/*
		Get a value.
		@param ok - bool, if the value was found
		@param value - the value
	*/
	if ok, value := gc.Get("hello"); !ok || value.(string) != "world" {
		t.Fatal("Could not retrieve value or value was not correct")
	}
	t.Log("GC set ok!")
	/*
		Overwrite the pair ("hello", "world") with the pair ("hello", "mars"), which
		never expires and won't be cleaned
	*/
	gc.SetNoTTL("hello", "mars")
	if ok, value := gc.Get("hello"); !ok || value.(string) != "mars" {
		t.Fatal("Could not retrieve value or value was not correct")
	}
	t.Log("GC set ok!")
	gc.Set("deleteMe", "lol", 2*time.Minute)
	/*
		Delete an key-value-pair manually
	*/
	gc.Delete("deleteMe")
	if ok, _ := gc.Get("deleteMe"); ok {
		t.Fatal("Deleted value is still there")
	}
	t.Log("GC delete ok!")
	time.Sleep(21 * time.Second)
	if ok, _ := gc.Get("hello"); ok {
		t.Fatal("Value should be deleted")
	}
	t.Log("GC clean ok!")
}
