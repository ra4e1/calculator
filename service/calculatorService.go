package service

import (
	"time"

	"github.com/Knetic/govaluate"
)

type CalculatorService struct {
	Delay int //задержка в секундах для притормаживания вычеслений
}

func NewCalculatorService(delay int) *CalculatorService { //создание
	return &CalculatorService{
		Delay: delay,
	}
}

func (c *CalculatorService) Calculate(line string) (answer interface{}, err error) { //решение примера
	expression, err := govaluate.NewEvaluableExpression(line)
	if err != nil {
		return nil, err
	}
	result, err := expression.Evaluate(nil)
	time.Sleep(time.Duration(c.Delay) * time.Second)
	return result, err
}
