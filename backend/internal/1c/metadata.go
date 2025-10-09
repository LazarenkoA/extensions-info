package onec

import (
	"encoding/xml"
	"errors"
	"log"
	"os"
	"path/filepath"
	"your-app/internal/models"
)

func (a *Analyzer1C) metadataAnalyzing(extDir string, confID int32) error {
	dirs, err := os.ReadDir(extDir)
	if err != nil {
		return err
	}

	var metadataInfo []*models.MetadataInfo

	for _, dir := range dirs {
		confStr, err := parseConfigurationFile(filepath.Join(extDir, dir.Name(), "Configuration.xml"))
		if err != nil {
			log.Printf("error analyzing the %q extension\n", dir)
			continue
		}

		metadataInfo = convToMetadataInfo(confStr)
	}

	_ = metadataInfo
	return nil
}

func parseConfigurationFile(path string) (*ConfigurationStruct, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

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

func convToMetadataInfo(cfg *ConfigurationStruct) []*models.MetadataInfo {
	var metadataInfo []*models.MetadataInfo

	typeHandlers := []struct {
		names []string
		t     models.ObjectType
		dir   string
	}{
		{cfg.Subsystems, models.ObjectTypeSubsystems, "Subsystems"},
		{cfg.Roles, models.ObjectTypeRoles, "ObjectTypeRoles"},
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
			for _, item := range h.names {
				metadataInfo = append(metadataInfo, &models.MetadataInfo{
					ObjectName: item,
					Type:       h.t,
					Path:       h.dir,
				})
			}
		}
	}

	return metadataInfo
}
