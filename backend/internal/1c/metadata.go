package onec

import (
	"context"
	"encoding/xml"
	"github.com/beevik/etree"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"your-app/internal/models"
	"your-app/internal/utils"
)

func (a *Analyzer1C) metadataAnalyzing(extDir string, confID int32) error {
	defer os.RemoveAll(extDir)

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

	doc := etree.NewDocument()
	_, err = doc.ReadFrom(f)
	if err != nil {
		log.Println(errors.Wrap(err, "read xml"))
		return new(models.MetadataInfo)
	}

	elemBorrowed := doc.FindElement("/MetaDataObject/*/Properties/ExtendedConfigurationObject")
	return &models.MetadataInfo{
		ObjectName: objectName,
		Type:       objectType,
		Borrowed:   elemBorrowed != nil,
	}
}

func merge(metadataInfoPrev, metadataInfoNew []*models.MetadataInfo, extID int32) []*models.MetadataInfo {
	for _, md := range metadataInfoNew {
		exist := false
		for j, mdPrev := range metadataInfoPrev {
			if md.Type == mdPrev.Type && md.ObjectName == mdPrev.ObjectName {
				metadataInfoPrev[j].ExtensionIDs = append(metadataInfoPrev[j].ExtensionIDs, extID)
				exist = true
				break
			}
		}

		if !exist {
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
