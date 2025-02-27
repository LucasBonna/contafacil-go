package services

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/shopspring/decimal"

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

	if cpfElem := dest.SelectElement("CPF"); cpfElem != nil {
		cpfCnpj = cpfElem.Text()
	} else if cnpjElem := dest.SelectElement("CNPJ"); cnpjElem != nil {
		cpfCnpj = cnpjElem.Text()
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
	vICMSUFDestDec, err := decimalFromString(vICMSUFDestElem.Text())
	if err != nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("valor inválido para vICMSUFDest: %w", err)
	}

	var icmsDec decimal.Decimal

	if vICMSUFDestDec.GreaterThan(decimal.Zero) {
		icmsDec = vICMSUFDestDec
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

		aliquotaDec, errAliquota := extractAliquotaIcmsDecimal(infCplText)
		if errAliquota == nil && aliquotaDec.GreaterThan(decimal.Zero) {

			vBCElem := ICMSTot.SelectElement("vBC")
			if vBCElem == nil {
				return schemas.ValidateAndProcess{}, fmt.Errorf("vBC não encontrado em ICMSTot")
			}

			vBCDec, errParseBC := decimalFromString(vBCElem.Text())
			if errParseBC != nil {
				return schemas.ValidateAndProcess{}, fmt.Errorf("valor inválido para vBC: %w", errParseBC)
			}

			icmsDec = vBCDec.Mul(aliquotaDec.Div(decimal.NewFromInt(100)))
		} else {

			estaduaisDec, errEstaduais := extractEstaduaisValueDecimal(infCplText)
			if errEstaduais != nil {
				return schemas.ValidateAndProcess{}, fmt.Errorf(
					"erro ao extrair valor de Estaduais: %w. Erro alíquota: %s",
					errEstaduais, errAliquota,
				)
			}

			if estaduaisDec.LessThanOrEqual(decimal.Zero) {
				return schemas.ValidateAndProcess{}, fmt.Errorf(
					"valor de Estaduais é zero ou inexistente",
				)
			}

			icmsDec = estaduaisDec
		}

		vICMSUFDestElem.SetText(icmsDec.Round(2).StringFixed(2))
	}

	vICMSUFDestElem.SetText(icmsDec.Round(2).StringFixed(2))

	processedXML, err := doc.WriteToString()
	if err != nil {
		return schemas.ValidateAndProcess{}, fmt.Errorf("erro ao re-encodar XML: %w", err)
	}

	return schemas.ValidateAndProcess{
		IcmsValue:    icmsDec.Round(2).InexactFloat64(),
		ChaveNota:    chaveNota,
		NumNota:      numNota,
		Destinatario: destinatario,
		CpfCnpj:      cpfCnpj,
		UF:           uf,
		ProcessedXML: processedXML,
	}, nil
}

// decimalFromString converte uma string para decimal.Decimal,
// trocando vírgula por ponto na primeira ocorrência.
func decimalFromString(strVal string) (decimal.Decimal, error) {
	strVal = strings.TrimSpace(strVal)
	// Se existir ',', substituir por '.' (apenas a primeira).
	strVal = strings.Replace(strVal, ",", ".", 1)
	return decimal.NewFromString(strVal)
}

// extractAliquotaIcmsDecimal extrai a alíquota do ICMS do estado de destino como decimal.
// Exemplo esperado no infCplText:
//
//	"Aliquota do ICMS do estado de destino 20,00"
func extractAliquotaIcmsDecimal(text string) (decimal.Decimal, error) {
	prefix := "Aliquota do ICMS do estado de destino "
	index := strings.Index(text, prefix)
	if index == -1 {
		return decimal.Zero, fmt.Errorf("não encontrou string de aliquota (prefixo '%s')", prefix)
	}

	substr := text[index+len(prefix):]

	// Geralmente vem "20,00" ou "20.00" antes de espaço, <br>, etc.
	endIndex := strings.IndexAny(substr, " <,(&\n\t\r")
	if endIndex == -1 {
		// Se não achou separador, pega tudo até o fim
		endIndex = len(substr)
	}
	aliqStr := strings.TrimSpace(substr[:endIndex])

	aliquotaDec, err := decimalFromString(aliqStr)
	if err != nil {
		return decimal.Zero, fmt.Errorf("valor inválido para aliquota: %w", err)
	}
	return aliquotaDec, nil
}

// extractEstaduaisValueDecimal extrai o valor numérico após "Estaduais R$ "
// Ex: "Estaduais R$ 45,67" → 45.67
func extractEstaduaisValueDecimal(text string) (decimal.Decimal, error) {
	prefix := "Estaduais R$ "
	index := strings.Index(text, prefix)
	if index == -1 {
		return decimal.Zero, fmt.Errorf("valor de Estaduais não encontrado")
	}

	substr := text[index+len(prefix):]

	endIndex := strings.IndexAny(substr, " (")
	if endIndex == -1 {
		endIndex = len(substr)
	}

	valueStr := strings.TrimSpace(substr[:endIndex])
	return decimalFromString(valueStr)
}
