package service

import (
	"encoding/json"
	"os"
)

type StateService struct{}

func NewStateService() *StateService { //создание
	return &StateService{}
}

type AnswerError struct{ Msg string }

type Answer struct {
	Err       *AnswerError
	Value     interface{} //ответ
	Ready     bool        //флаг
	Expresion string      //пример
}

func (s *StateService) SaveState(m map[int]*Answer) error {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	jsonFile, err := os.Create("state.json")
	defer jsonFile.Close()
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (s *StateService) RestoreState() (map[int]*Answer, int, error) {
	f, err := os.Open("state.json")
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	result, err := os.ReadFile("state.json")
	if err != nil {
		return nil, 0, err
	}
	var ma map[int]*Answer
	err = json.Unmarshal([]byte(result), &ma)
	if err != nil {
		return nil, 0, err
	}

	var requestID int
	for id, answer := range ma {
		if requestID < id {
			requestID = id
		}
		if !answer.Ready {
			answer.Ready = true
			answer.Err = &AnswerError{Msg: "отменено"}
		}
	}
	return ma, requestID, nil
}
