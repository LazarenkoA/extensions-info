package models

import "time"

type ObjectType string
type Redefinition string

const (
	ObjectTypeUndefined    ObjectType = ""
	ObjectTypeConf         ObjectType = "configuration"
	ObjectTypeDocument     ObjectType = "document"
	ObjectTypeCatalog      ObjectType = "catalog"
	ObjectTypeCommonModule ObjectType = "commonModule"
	ObjectTypeFunction     ObjectType = "function"

	// ....
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
