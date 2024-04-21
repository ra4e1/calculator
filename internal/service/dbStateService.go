package service

type DbStateService struct {
	db *DbService
}

func NewDbStateService(db *DbService) *DbStateService { //создание
	return &DbStateService{
		db: db,
	}
}

type ExpressionError struct{ Msg string }

type Expression struct {
	ID         int64
	Err        *ExpressionError
	UserId     int64
	Value      interface{} //ответ
	Ready      bool        //флаг
	Expression string      //пример
}

func bool2Int(val bool) int {
	var result int
	if val {
		result = 1
	} else {
		result = 0
	}
	return result
}

func mapToExpression(dbExp *DbExpression) *Expression {
	var expErr *ExpressionError
	if dbExp.Err != "" {
		expErr = &ExpressionError{Msg: dbExp.Err}
	}

	exp := &Expression{
		ID:         dbExp.ID,
		UserId:     dbExp.UserId,
		Expression: dbExp.Expresion,
		Ready:      dbExp.Ready != 0,
		Err:        expErr,
		Value:      dbExp.Answer,
	}
	return exp
}

func mapToDbExpression(exp *Expression) *DbExpression {
	var err string
	if exp.Err != nil {
		err = exp.Err.Msg
	}

	var value float64
	if exp.Value != nil {
		value = exp.Value.(float64)
	}

	dbexp := &DbExpression{
		ID:        exp.ID,
		UserId:    exp.UserId,
		Ready:     bool2Int(exp.Ready),
		Answer:    value,
		Err:       err,
		Expresion: exp.Expression,
	}
	return dbexp
}

func (s *DbStateService) AddExpression(expression *Expression) (int64, error) {
	dbexp := mapToDbExpression(expression)
	id, err := s.db.InsertExpression(dbexp)
	return id, err
}

func (s *DbStateService) UpdateExpression(id int64, expression *Expression) error {
	dbexp := mapToDbExpression(expression)
	s.db.UpdateExpression(id, dbexp)
	return nil
}

func (s *DbStateService) FindUserExpression(userId, id int64) (*Expression, error) {
	dbExp, err := s.db.FindExpressionById(userId, id)
	if err != nil {
		return nil, err
	}
	return mapToExpression(dbExp), nil
}

func (s *DbStateService) FindUserAllExpressions(userId int64) ([]*Expression, error) {
	dbExps, err := s.db.FindExpressions(userId)
	if err != nil {
		return nil, err
	}

	var expressions []*Expression
	for _, v := range dbExps {
		exp := mapToExpression(v)
		expressions = append(expressions, exp)
	}

	return expressions, nil
}
