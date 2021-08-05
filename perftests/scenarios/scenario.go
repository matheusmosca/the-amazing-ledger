package scenarios

type Scenario interface {
	GetMethod() string
	GetJSON() string
}
