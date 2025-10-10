package onec

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/LazarenkoA/extensions-info/internal/models"
	"github.com/LazarenkoA/extensions-info/internal/utils"
	"github.com/antchfx/xmlquery"
	"github.com/beevik/etree"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (a *Analyzer1C) metadataAnalyzing(extDir string, confID int32) error {
	dirs, err := os.ReadDir(extDir)
	if err != nil {
		return err
	}

	extensionsRep, _ := a.repo.GetExtensionsInfo(context.Background(), confID)
	extensions := make(map[string]int32, len(extensionsRep))
	for _, e := range extensionsRep {
		extensions[e.Name] = e.ID
	}

	var metadataInfo []*models.MetadataInfo

	for _, dir := range dirs {
		confStr, err := parseConfigurationFile(filepath.Join(extDir, dir.Name(), "Configuration.xml"))
		if err != nil {
			log.Printf("error analyzing the %q extension\n", dir)
			continue
		}

		metadataInfo = merge(metadataInfo, convToMetadataInfo(confStr, filepath.Join(extDir, dir.Name())), extensions[dir.Name()])
	}

	data, _ := json.Marshal(metadataInfo)
	return a.repo.SetMetadata(context.Background(), confID, data)
}

func (a *Analyzer1C) codeAnalyzing(extDir string, confID int32) error {
	// получаем из БД структуру метаданных ранее проанализированную
	// нас интересуют только заимствованные объекты, будем проверять их

	data, err := a.repo.GetMetadata(context.Background(), confID)
	if err != nil {
		return errors.Wrap(err, "couldn't get the configuration structure")
	}

	var metadata []*models.MetadataInfo
	_ = json.Unmarshal(data, &metadata)

	for _, md := range metadata {
		if utils.PtrToVal(md.Borrowed) {
			fmt.Println(1)
		}
	}

	return nil
}

func parseConfigurationFile(path string) (*ConfigurationStruct, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	decoder := xml.NewDecoder(f)
	for {
		tok, _ := decoder.Token()
		if tok == nil {
			break
		}

		if se, ok := tok.(xml.StartElement); ok && se.Name.Local == "ChildObjects" {
			var child ConfigurationStruct
			if err := decoder.DecodeElement(&child, &se); err != nil {
				return nil, err
			}

			return &child, nil
		}
	}

	return nil, errors.New("no ChildObjects found in xml")
}

func convToMetadataInfo(cfg *ConfigurationStruct, rootDir string) []*models.MetadataInfo {
	var metadataInfo []*models.MetadataInfo

	typeHandlers := []struct {
		names []string
		t     models.ObjectType
		dir   string
	}{
		{cfg.Subsystems, models.ObjectTypeSubsystems, "Subsystems"},
		{cfg.Roles, models.ObjectTypeRoles, "Roles"},
		{cfg.CommonModules, models.ObjectTypeCommonModules, "CommonModules"},
		{cfg.ExchangePlans, models.ObjectTypeExchangePlans, "ExchangePlans"},
		{cfg.HTTPServices, models.ObjectTypeHTTPServices, "HTTPServices"},
		{cfg.EventSubscriptions, models.ObjectTypeEventSubscriptions, "EventSubscriptions"},
		{cfg.ScheduledJobs, models.ObjectTypeScheduledJobs, "ScheduledJobs"},
		{cfg.DefinedTypes, models.ObjectTypeDefinedTypes, "DefinedTypes"},
		{cfg.Constants, models.ObjectTypeConstants, "Constants"},
		{cfg.Catalogs, models.ObjectTypeCatalogs, "Catalogs"},
		{cfg.Documents, models.ObjectTypeDocuments, "Documents"},
		{cfg.DocumentJournals, models.ObjectTypeDocumentJournals, "DocumentJournals"},
		{cfg.Enums, models.ObjectTypeEnums, "Enums"},
		{cfg.Reports, models.ObjectTypeReports, "Reports"},
		{cfg.DataProcessors, models.ObjectTypeDataProcessors, "DataProcessors"},
		{cfg.InformationRegisters, models.ObjectTypeInformationRegisters, "InformationRegisters"},
		{cfg.AccumulationRegisters, models.ObjectTypeAccumulationRegisters, "AccumulationRegisters"},
		{cfg.ChartsOfCharacteristic, models.ObjectTypeChartsOfCharacteristic, "ChartsOfCharacteristic"},
		{cfg.ChartsOfAccounts, models.ObjectTypeChartsOfAccounts, "ChartsOfAccounts"},
		{cfg.AccountingRegisters, models.ObjectTypeAccountingRegisters, "AccountingRegisters"},
		{cfg.ChartsOfCalculation, models.ObjectTypeChartsOfCalculation, "ChartsOfCalculation"},
		{cfg.CalculationRegisters, models.ObjectTypeCalculationRegisters, "CalculationRegisters"},
		{cfg.BusinessProcesses, models.ObjectTypeBusinessProcesses, "BusinessProcesses"},
		{cfg.Tasks, models.ObjectTypeTasks, "Tasks"},
	}

	for _, h := range typeHandlers {
		if len(h.names) > 0 {
			for _, name := range h.names {
				metadataInfo = append(metadataInfo, readExtensionObject(name, h.t, filepath.Join(rootDir, h.dir)))
			}
		}
	}

	return metadataInfo
}

func readExtensionObject(objectName string, objectType models.ObjectType, dir string) *models.MetadataInfo {
	f, err := os.Open(filepath.Join(dir, objectName+".xml"))
	if err != nil {
		log.Println(errors.Wrap(err, "open object file"))
		return new(models.MetadataInfo)
	}
	defer f.Close()

	doc, err := xmlquery.Parse(f)
	if err != nil {
		log.Println(errors.Wrap(err, "read xml"))
		return new(models.MetadataInfo)
	}

	elemBorrowed := doc.SelectElement("/MetaDataObject/*/Properties/ExtendedConfigurationObject")
	return &models.MetadataInfo{
		ObjectName: objectName,
		Type:       objectType,
		Borrowed:   utils.Ptr(elemBorrowed != nil),
		ID:         uuid.NewString(),
		Changes_:   changes(dir, doc),
		Changes:    map[int32][]string{},
	}
}

func changes(rootDir string, doc *xmlquery.Node) []string {
	var result []string

	elemName := doc.SelectElement("/MetaDataObject/*/Properties/Name/text()")
	// реквизиты
	{
		addAttr, extAttr := changesElement(doc, "/MetaDataObject/*/ChildObjects/Attribute/Properties/Name")
		if len(addAttr) > 0 || len(extAttr) > 0 {
			var txt strings.Builder
			if len(addAttr) > 0 || len(extAttr) > 0 {
				txt.WriteString("Реквизиты:\n")
			}
			if len(addAttr) > 0 {
				txt.WriteString("- Добавленные: " + strings.Join(addAttr, ", ") + "\n")
			}
			if len(extAttr) > 0 {
				txt.WriteString("- Заимствованные: " + strings.Join(extAttr, ", ") + "\n")
			}

			result = append(result, txt.String())
		}
	}

	// таб части
	{
		addAttr, extAttr := changesElement(doc, "/MetaDataObject/*/ChildObjects/TabularSection/Properties/Name")
		if len(addAttr) > 0 || len(extAttr) > 0 {
			var txt strings.Builder
			if len(addAttr) > 0 || len(extAttr) > 0 {
				txt.WriteString("Таб. часть:\n")
			}
			if len(addAttr) > 0 {
				txt.WriteString("- Добавленные: " + strings.Join(addAttr, ", ") + "\n")
			}
			if len(extAttr) > 0 {
				txt.WriteString("- Заимствованные: " + strings.Join(extAttr, ", ") + "\n")
			}

			result = append(result, txt.String())
		}
	}

	// формы
	{
		var txt strings.Builder
		elemForms := doc.SelectElements("/MetaDataObject/*/ChildObjects/Form/text()")
		if len(elemForms) > 0 {
			txt.WriteString("Формы:\n")
		}

		for _, form := range elemForms {
			f, err := os.Open(filepath.Join(rootDir, fmt.Sprintf("%s\\Forms\\%s.xml", utils.Opt(elemName).Data, form.Data)))

			if err == nil {
				doc, err := xmlquery.Parse(f)
				if err == nil {
					addAttr, extAttr := changesElement(doc, "/MetaDataObject/Form/Properties/Name")
					if len(addAttr) > 0 || len(extAttr) > 0 {
						if len(addAttr) > 0 {
							txt.WriteString("- Добавленные: " + strings.Join(addAttr, ", ") + "\n")
						}
						if len(extAttr) > 0 {
							txt.WriteString("- Заимствованные: " + strings.Join(extAttr, ", ") + "\n")
						}
					}
				}

				f.Close()
			}

		}
		if txt.Len() > 0 {
			result = append(result, txt.String())
		}
	}

	return result
}

func changesElement(doc *xmlquery.Node, xpath string) ([]string, []string) {
	elemAttributeName := doc.SelectElements(xpath)
	addAttr, extAttr := []string{}, []string{}
	for _, nameElem := range elemAttributeName {
		txt := nameElem.SelectElement("text()")
		ext := nameElem.Parent.SelectElement("ExtendedConfigurationObject")
		if ext != nil {
			extAttr = append(extAttr, txt.Data)
		} else {
			addAttr = append(addAttr, txt.Data)
		}
	}

	return addAttr, extAttr
}

func merge(metadataInfoPrev, metadataInfoNew []*models.MetadataInfo, extID int32) []*models.MetadataInfo {
	for _, md := range metadataInfoNew {
		exist := false
		for j, mdPrev := range metadataInfoPrev {
			if md.Type == mdPrev.Type && md.ObjectName == mdPrev.ObjectName {
				metadataInfoPrev[j].ExtensionIDs = append(metadataInfoPrev[j].ExtensionIDs, extID)
				metadataInfoPrev[j].Changes[extID] = append(metadataInfoPrev[j].Changes[extID], md.Changes_...)
				exist = true
				break
			}
		}

		if !exist {
			md.Changes[extID] = append(md.Changes[extID], md.Changes_...)
			md.ExtensionIDs = append(md.ExtensionIDs, extID)
			metadataInfoPrev = append(metadataInfoPrev, md)
		}
	}

	return metadataInfoPrev
}

func readConfigurationFile(dir string) (*ConfigurationInfo, error) {
	f, err := os.Open(filepath.Join(dir, "Configuration.xml"))
	if err != nil {
		return nil, errors.Wrap(err, "open Configuration.xml error")
	}
	defer f.Close()

	doc := etree.NewDocument()
	_, err = doc.ReadFrom(f)
	if err != nil {
		return nil, errors.Wrap(err, "read Configuration.xml error")
	}

	elemName := doc.FindElement("/MetaDataObject/Configuration/Properties/Name")
	elemSynonym := doc.FindElement("/MetaDataObject/Configuration/Properties/Synonym/v8:item/v8:content")
	elemVersion := doc.FindElement("/MetaDataObject/Configuration/Properties/Version")
	elemVendor := doc.FindElement("/MetaDataObject/Configuration/Properties/Vendor")
	elemPurpose := doc.FindElement("/MetaDataObject/Configuration/Properties/ConfigurationExtensionPurpose")
	return &ConfigurationInfo{
		Name:    utils.Ptr(utils.Opt[etree.Element](elemName)).Text(),
		Synonym: utils.Ptr(utils.Opt[etree.Element](elemSynonym)).Text(),
		Version: utils.Ptr(utils.Opt[etree.Element](elemVersion)).Text(),
		Vendor:  utils.Ptr(utils.Opt[etree.Element](elemVendor)).Text(),
		Purpose: utils.Ptr(utils.Opt[etree.Element](elemPurpose)).Text(),
	}, nil
}
