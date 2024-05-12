package main

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	gc := New(20*time.Second, 100)
	gc.Run()
	if gc == nil {
		t.Fatal("GC is nil")
	}
	t.Log("GC init ok!")
	gc.Set("hello", "world", 1*time.Minute)
	if ok, value := gc.Get("hello"); !ok || value.(string) != "world" {
		t.Fatal("Could not retrieve value or value was not correct")
	}
	t.Log("GC set ok!")
	gc.Set("hello", "mars", 10*time.Second)
	if ok, value := gc.Get("hello"); !ok || value.(string) != "mars" {
		t.Fatal("Could not retrieve value or value was not correct")
	}
	t.Log("GC set ok!")
	gc.Set("deleteMe", "lol", 2*time.Minute)
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
