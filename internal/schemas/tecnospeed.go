package schemas

type IssueGNREResponse struct {
	Sucess  *IssueGNREResponseSucess
	Failure *IssueGNREResponseFailure
}

type IssueGNREResponseSucess struct {
	NumRecibo    string
	Situacao     string
	Motivo       string
	UFFavorecida string
	Receita      string
}

type IssueGNREResponseFailure struct {
	Exception string
	Class     string
	Message   string
}
