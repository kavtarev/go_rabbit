package main

import (
	"fmt"
	"testing"
)

func TestQueryOrder(t *testing.T) {
	inputs := []string{"desc", "DESC", "Desc"}
	for _, name := range inputs {
		testName := fmt.Sprintf("test_%v", name)
		t.Run(testName, func(t *testing.T) {
			val := make(map[string][]string)
			val["order"] = []string{name}

			par := QueryParamsParser{
				values: val,
			}
			par.parseOrder()

			if par.Order != defaultOrder {
				t.Errorf("on value %v got error parsing order", name)
			}
		})
	}
}

func TestQueryPage(t *testing.T) {
	val := make(map[string][]string)
	val["page"] = []string{"0tt"}

	par := QueryParamsParser{
		values: val,
	}

	isCorrect := par.CheckCorrectness()
	if isCorrect {
		t.Error("not correct when page is string")
	}
}
