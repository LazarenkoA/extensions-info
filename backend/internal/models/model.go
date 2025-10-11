package models

import "time"

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

type ExtChanges struct {
	ID                int32
	MetadataChanges   []string
	FuncsChanges      []FuncInfo
	PathObject        string
	ExtensionRootPath string
}

type ChildrenObject []*MetadataInfo

type MetadataInfo struct {
	ObjectName string
	Type       string
	Children   ChildrenObject
	Extension  []*ExtChanges
	Borrowed   *bool
	ID         string
}

type FuncInfo struct {
	Directive string
	Type      string
	Name      string
	Code      string
	ModuleKey string
}

type AppSettings struct {
	ID           int32
	PlatformPath string
}

func (c ChildrenObject) Find(name string, otype string) *MetadataInfo {
	for _, item := range c {
		if item.ObjectName == name && item.Type == otype {
			return item
		}
	}

	return nil
}
