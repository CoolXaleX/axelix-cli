package models

// ConditionsFeed is the response of the conditions actuator endpoint.
type ConditionsFeed struct {
	PositiveMatches []PositiveCondition `json:"positiveMatches"`
	NegativeMatches []NegativeCondition `json:"negativeMatches"`
}

// PositiveCondition represents a configuration class where all conditions matched.
type PositiveCondition struct {
	ClassName  string           `json:"className"`
	MethodName *string          `json:"methodName"`
	Matched    []ConditionMatch `json:"matched"`
}

// NegativeCondition represents a configuration class where some conditions did not match.
type NegativeCondition struct {
	ClassName  string           `json:"className"`
	MethodName *string          `json:"methodName"`
	NotMatched []ConditionMatch `json:"notMatched"`
	Matched    []ConditionMatch `json:"matched"`
}

// ConditionMatch represents the result of evaluating a single condition.
type ConditionMatch struct {
	Condition string `json:"condition"`
	Message   string `json:"message"`
}
