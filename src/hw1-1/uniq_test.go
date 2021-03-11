package main

import (
	"reflect"
	"testing"
)

func TestDefaultSuccess(t *testing.T) {
	data := []string{
		"I love music.",
		"I love music.",
		"I love music.",
		"",
		"I love music of Kartik.",
		"I love music of Kartik.",
		"Thanks.",
	}
	options := &Options{
		C: false,
		D: true,
		U: false,
		F: 0,
		S: 0,
		I: false,
	}
	expected := []string{
		"I love music.",
		"I love music of Kartik.",
	}
	result := execute(options, data)

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Check default behaviour failed")
	}
}

func TestUniqueSuccess(t *testing.T) {
	data := []string{
		"I love music.",
		"I love music.",
		"I love music.",
		"",
		"I love music of Kartik.",
		"I love music of Kartik.",
		"Thanks.",
	}
	options := &Options{
		C: false,
		D: false,
		U: true,
		F: 0,
		S: 0,
		I: false,
	}
	expected := []string{
		"",
		"Thanks.",
	}
	result := execute(options, data)

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Check default behaviour failed")
	}
}

func TestCountSuccess(t *testing.T) {
	data := []string{
		"I love music.",
		"I love music.",
		"I love music.",
		"",
		"I love music of Kartik.",
		"I love music of Kartik.",
		"Thanks.",
	}
	options := &Options{
		C: true,
		D: false,
		U: false,
		F: 0,
		S: 0,
		I: false,
	}
	expected := []string{
		"3 I love music.",
		"1",
		"2 I love music of Kartik.",
		"1 Thanks.",
	}
	result := execute(options, data)

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Check default behaviour failed")
	}
}

func TestSkipFieldSuccess(t *testing.T) {
	data := []string{
		"I love music.",
		"A love music.",
		"C love music.",
		"",
		"I love music of Kartik.",
		"We love music of Kartik.",
		"Thanks.",
	}
	options := &Options{
		C: false,
		D: false,
		U: false,
		F: 1,
		S: 0,
		I: false,
	}
	expected := []string{
		"I love music.",
		"",
		"I love music of Kartik.",
		"Thanks.",
	}
	result := execute(options, data)

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Check default behaviour failed")
	}
}

func TestSkipCharsSuccess(t *testing.T) {
	data := []string{
		"I love music.",
		"A love music.",
		"C love music.",
		"",
		"I love music of Kartik.",
		"We love music of Kartik.",
		"Thanks.",
	}
	options := &Options{
		C: false,
		D: false,
		U: false,
		F: 1,
		S: 0,
		I: false,
	}
	expected := []string{
		"I love music.",
		"",
		"I love music of Kartik.",
		"Thanks.",
	}
	result := execute(options, data)

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Check default behaviour failed")
	}
}

func TestRegSuccess(t *testing.T) {
	data := []string{
		"I LOVE MUSIC.",
		"I love music.",
		"I LoVe MuSiC.",
		"",
		"I love MuSIC of Kartik.",
		"I love music of kartik.",
		"Thanks.",
	}
	options := &Options{
		C: false,
		D: false,
		U: false,
		F: 0,
		S: 0,
		I: true,
	}
	expected := []string{
		"I LOVE MUSIC.",
		"",
		"I love MuSIC of Kartik.",
		"Thanks.",
	}
	result := execute(options, data)

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Check default behaviour failed")
	}
}

func TestCountFail(t *testing.T) {
	data := []string{
		"I love music.",
		"I love music.",
		"I love music.",
		"",
		"I love music of Kartik.",
		"I love music of Kartik.",
		"Thanks.",
	}
	options := &Options{
		C: true,
		D: false,
		U: false,
		F: 0,
		S: 0,
		I: false,
	}
	expected := []string{
		"2 I love music.",
		"1",
		"2 I love music of Kartik.",
		"1 Thanks.",
	}
	result := execute(options, data)

	if reflect.DeepEqual(expected, result) {
		t.Fatalf("Check default behaviour failed")
	}
}
