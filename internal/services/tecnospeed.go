package services

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

type Tecnospeed interface {
	IssueGNRE(xmlContent string, group string, document string) (*schemas.IssueGNREResponse, error)
	DownloadGNRE(group string, document string, chaveNota string, numRecibo string) ([]byte, error)
}

type tecnospeedService struct {
	client   *resty.Client
	username string
	password string
	baseUrl  string
}

func NewTecnospeedService(client *resty.Client, username string, password string, baseUrl string) Tecnospeed {
	return &tecnospeedService{
		client:   client,
		username: username,
		password: password,
		baseUrl:  baseUrl,
	}
}

func (ts *tecnospeedService) IssueGNRE(xmlContent string, group string, document string) (*schemas.IssueGNREResponse, error) {
	endpoint := fmt.Sprintf("%s/gnre/envia", ts.baseUrl)
	resp, err := ts.client.R().
		SetBasicAuth(ts.username, ts.password).
		SetFormData(map[string]string{
			"Grupo":   group,
			"CNPJ":    document,
			"Arquivo": xmlContent,
		}).
		SetResult(&schemas.IssueGNREResponseSucess{}).
		Post(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error issuing gnre %s", err)
	}

	responseBody := resp.String()
	responseBody = strings.TrimSpace(responseBody)

	issueResponse := &schemas.IssueGNREResponse{}

	if strings.HasPrefix(responseBody, "EXCEPTION") {
		parts := strings.Split(responseBody, ",")
		if len(parts) < 3 {
			return nil, fmt.Errorf("formato de erro inesperado")
		}
		issueResponse.Failure = &schemas.IssueGNREResponseFailure{
			Exception: strings.Trim(parts[0], `"`),
			Class:     strings.Trim(parts[1], `"`),
			Message:   strings.Trim(parts[2], `"`),
		}
	} else {
		parts := strings.Split(responseBody, ",")
		if len(parts) < 5 {
			return nil, fmt.Errorf("formato de sucesso inesperado")
		}
		issueResponse.Sucess = &schemas.IssueGNREResponseSucess{
			NumRecibo:    strings.TrimSpace(parts[0]),
			Situacao:     strings.TrimSpace(parts[1]),
			Motivo:       strings.TrimSpace(parts[2]),
			UFFavorecida: strings.TrimSpace(parts[3]),
			Receita:      strings.TrimSpace(parts[4]),
		}
	}

	return issueResponse, nil
}

func (ts *tecnospeedService) DownloadGNRE(group string, document string, chaveNota string, numRecibo string) ([]byte, error) {
	var fileBytes []byte
	endpoint := fmt.Sprintf("%s/gnre/imprime", ts.baseUrl)
	_, err := ts.client.R().
		SetBasicAuth(ts.username, ts.password).
		SetQueryParams(map[string]string{
			"Grupo":        group,
			"CNPJ":         document,
			"ChaveNota":    chaveNota,
			"Url":          "0",
			"NumeroRecibo": numRecibo,
		}).
		SetResult(fileBytes).
		Post(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error issuing gnre %s", err)
	}

	return fileBytes, nil
}
