package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"

	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

type XML interface {
	ValidateAndProcess([]byte) (schemas.ValidateAndProcess, error)
}

type xmlService struct{}

func NewXmlService() XML {
	return &xmlService{}
}

func (xs *xmlService) ValidateAndProcess(xmlBytes []byte) (schemas.ValidateAndProcess, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlBytes); err != nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("XML inválido: %w", err)
	}

	nfeProc := doc.FindElement("./nfeProc")
	if nfeProc == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("nfeProc não encontrado")
	}

	protNFe := nfeProc.FindElement("protNFe")
	if protNFe == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("protNFe não encontrado")
	}

	infProt := protNFe.FindElement("infProt")
	if infProt == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("infProt não encontrado")
	}

	chNFe := infProt.SelectElement("chNFe")
	if chNFe == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("chNFe não encontrado")
	}
	chaveNota := chNFe.Text()

	nfe := nfeProc.FindElement("NFe")
	if nfe == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("NFe não encontrado")
	}

	infNFe := nfe.FindElement("infNFe")
	if infNFe == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("infNFe não encontrado")
	}

	ide := infNFe.FindElement("ide")
	if ide == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("ide não encontrado")
	}

	nNF := ide.SelectElement("nNF")
	if nNF == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("nNF não encontrado")
	}
	numNota := nNF.Text()

	dest := infNFe.FindElement("dest")
	if dest == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("dest não encontrado")
	}

	var cpfCnpj, destinatario, uf string

	cpfElem := dest.SelectElement("CPF")
	if cpfElem != nil {
		cpfCnpj = cpfElem.Text()
	} else {
		cnpjElem := dest.SelectElement("CNPJ")
		if cnpjElem != nil {
			cpfCnpj = cnpjElem.Text()
		}
	}

	xNome := dest.SelectElement("xNome")
	if xNome == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("xNome não encontrado")
	}
	destinatario = xNome.Text()

	enderDest := dest.SelectElement("enderDest")
	if enderDest == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("enderDest não encontrado")
	}

	UF := enderDest.SelectElement("UF")
	if UF == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("UF não encontrado")
	}
	uf = UF.Text()

	total := infNFe.SelectElement("total")
	if total == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("total não encontrado")
	}

	ICMSTot := total.SelectElement("ICMSTot")
	if ICMSTot == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("ICMSTot não encontrado")
	}

	vICMSUFDestElem := ICMSTot.SelectElement("vICMSUFDest")
	if vICMSUFDestElem == nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("vICMSUFDest não encontrado")
	}

	vICMSUFDestStr := strings.TrimSpace(vICMSUFDestElem.Text())
	vICMSUFDest, err := strconv.ParseFloat(strings.Replace(vICMSUFDestStr, ",", ".", 1), 64)
	if err != nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("valor inválido para vICMSUFDest: %w", err)
	}

	var icmsValue float64

	if vICMSUFDest != 0 {
		icmsValue = vICMSUFDest
	} else {
		infAdic := infNFe.SelectElement("infAdic")
		if infAdic == nil {
			return schemas.ValidateAndProcess{}, fmt.Errorf("infAdic não encontrado")
		}

		infCpl := infAdic.SelectElement("infCpl")
		if infCpl == nil {
			return schemas.ValidateAndProcess{}, fmt.Errorf("infCpl não encontrado")
		}

		infCplText := infCpl.Text()

		estaduaisValue, err := extractEstaduaisValue(infCplText)
		if err != nil {
			return schemas.ValidateAndProcess{}, fmt.Errorf("erro ao extrair valor de Estaduais: %w", err)
		}

		if estaduaisValue == 0 {
			return schemas.ValidateAndProcess{}, fmt.Errorf("valor de Estaduais é zero ou inexistente")
		}

		icmsValue = estaduaisValue

		vICMSUFDestElem.SetText(fmt.Sprintf("%.2f", icmsValue))
	}

	processedXML, err := doc.WriteToString()
	if err != nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("erro ao re-encodar XML: %w", err)
	}

	return schemas.ValidateAndProcess{
		IcmsValue:    icmsValue,
		ChaveNota:    chaveNota,
		NumNota:      numNota,
		Destinatario: destinatario,
		CpfCnpj:      cpfCnpj,
		UF:           uf,
		ProcessedXML: processedXML,
	}, nil
}

func extractEstaduaisValue(text string) (float64, error) {
	prefix := "Estaduais R$ "
	index := strings.Index(text, prefix)
	if index == -1 {
		return 0, fmt.Errorf("valor de Estaduais não encontrado")
	}

	substr := text[index+len(prefix):]

	endIndex := strings.IndexAny(substr, " (")
	if endIndex == -1 {
		endIndex = len(substr)
	}

	valueStr := strings.TrimSpace(substr[:endIndex])
	valueStr = strings.Replace(valueStr, ",", ".", 1)

	estaduaisValue, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("valor inválido para Estaduais: %w", err)
	}

	return estaduaisValue, nil
}
