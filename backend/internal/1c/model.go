package onec

type ConfigurationInfo struct {
	ID           int32
	Name         string
	Synonym      string
	Version      string
	Vendor       string
	Purpose      string
	ChildObjects *ConfigurationStruct
}

type MetadataObjectStruct struct {
}

type ConfigurationStruct struct {
	Subsystems             []string `xml:"Subsystem"`
	Roles                  []string `xml:"Role"`
	CommonModules          []string `xml:"CommonModule"`
	ExchangePlans          []string `xml:"ExchangePlan"`
	HTTPServices           []string `xml:"HTTPService"`
	EventSubscriptions     []string `xml:"EventSubscription"`
	ScheduledJobs          []string `xml:"ScheduledJob"`
	DefinedTypes           []string `xml:"DefinedType"`
	Constants              []string `xml:"Constant"`
	Catalogs               []string `xml:"Catalog"`
	Documents              []string `xml:"Document"`
	DocumentJournals       []string `xml:"DocumentJournal"`
	Enums                  []string `xml:"Enum"`
	Reports                []string `xml:"Report"`
	DataProcessors         []string `xml:"DataProcessor"`
	InformationRegisters   []string `xml:"InformationRegister"`
	AccumulationRegisters  []string `xml:"AccumulationRegister"`
	ChartsOfCharacteristic []string `xml:"ChartOfCharacteristicTypes"`
	ChartsOfAccounts       []string `xml:"ChartOfAccounts"`
	AccountingRegisters    []string `xml:"AccountingRegister"`
	ChartsOfCalculation    []string `xml:"ChartOfCalculationTypes"`
	CalculationRegisters   []string `xml:"CalculationRegister"`
	BusinessProcesses      []string `xml:"BusinessProcess"`
	Tasks                  []string `xml:"Task"`
}

var ruName = map[string]string{}
