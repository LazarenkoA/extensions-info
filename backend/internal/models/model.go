package models

import "time"

type ObjectType string
type Redefinition string

const (
	ObjectTypeUndefined              ObjectType = ""
	ObjectTypeConf                   ObjectType = "configuration"
	ObjectTypeLanguage               ObjectType = "language"
	ObjectTypeSubsystems             ObjectType = "subsystem"
	ObjectTypeRoles                  ObjectType = "role"
	ObjectTypeCommonModules          ObjectType = "commonModule"
	ObjectTypeExchangePlans          ObjectType = "exchangePlan"
	ObjectTypeHTTPServices           ObjectType = "httpService"
	ObjectTypeEventSubscriptions     ObjectType = "eventSubscription"
	ObjectTypeScheduledJobs          ObjectType = "scheduledJob"
	ObjectTypeDefinedTypes           ObjectType = "definedType"
	ObjectTypeConstants              ObjectType = "constant"
	ObjectTypeCatalogs               ObjectType = "catalog"
	ObjectTypeDocuments              ObjectType = "document"
	ObjectTypeDocumentJournals       ObjectType = "documentJournal"
	ObjectTypeEnums                  ObjectType = "enum"
	ObjectTypeReports                ObjectType = "report"
	ObjectTypeDataProcessors         ObjectType = "dataProcessor"
	ObjectTypeInformationRegisters   ObjectType = "informationRegister"
	ObjectTypeAccumulationRegisters  ObjectType = "accumulationRegister"
	ObjectTypeChartsOfCharacteristic ObjectType = "chartOfCharacteristicTypes"
	ObjectTypeChartsOfAccounts       ObjectType = "chartOfAccounts"
	ObjectTypeAccountingRegisters    ObjectType = "accountingRegister"
	ObjectTypeChartsOfCalculation    ObjectType = "chartOfCalculationTypes"
	ObjectTypeCalculationRegisters   ObjectType = "calculationRegister"
	ObjectTypeBusinessProcesses      ObjectType = "businessProcess"
	ObjectTypeTasks                  ObjectType = "task"
	ObjectTypeFunction               ObjectType = "function"
	// и остальное
)

const (
	RedefinitionUndefined     Redefinition = ""
	RedefinitionBefore        Redefinition = "Перед"
	RedefinitionAfter         Redefinition = "После"
	RedefinitionChangeControl Redefinition = "ИзменениеИКонтроль"
	RedefinitionInstead       Redefinition = "Вместо"
)

type CRONInfo struct {
	Schedule          string
	NextCheck         time.Time
	NextCheckAsString string
	Status            string
	DatabaseID        int32
}

type DatabaseSettings struct {
	Status            string
	Name              string
	ConnectionString  string
	Username          *string
	Password          *string
	ID                int32
	LastCheck         *time.Time
	LastCheckAsString string
	Cron              *CRONInfo
}

type ConfigurationInfo struct {
	DatabaseID   int32
	ID           int32
	Description  *string
	Version      string
	Name         string
	LastAnalysis time.Time
	Extensions   []ExtensionsInfo
	MetadataTree *MetadataInfo
}

type ExtensionsInfo struct {
	ConfigurationInfo
	Purpose string
}

type MetadataInfo struct {
	ObjectName   string
	Path         string
	Type         ObjectType
	Funcs        []FuncInfo
	Children     []*MetadataInfo
	ExtensionIDs []int32
	Borrowed     bool
}

type FuncInfo struct {
	RedefinitionMethod Redefinition
	Type               ObjectType
	Name               string
	Code               string
	ExtensionIDs       []int32
}

type AppSettings struct {
	ID           int32
	PlatformPath string
}
