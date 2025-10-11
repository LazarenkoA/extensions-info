package onec

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	astpkg "github.com/LazarenkoA/1c-language-parser/ast"
	"github.com/LazarenkoA/extensions-info/internal/models"
	"github.com/LazarenkoA/extensions-info/internal/utils"
	"github.com/antchfx/xmlquery"
	"github.com/beevik/etree"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (a *Analyzer1C) metadataAnalyzing(extRootDir string, confID int32) error {
	dirs, err := os.ReadDir(extRootDir)
	if err != nil {
		return err
	}

	extensionsRep, _ := a.repo.GetExtensionsInfo(context.Background(), confID)
	extensions := make(map[string]int32, len(extensionsRep))
	for _, e := range extensionsRep {
		extensions[e.Name] = e.ID
	}

	var metadataInfo []*models.MetadataInfo
	cfgStr, err := a.repo.GetChildObjectsConf(context.Background(), confID)
	if err != nil {
		log.Println("GetChildObjectsConf error", err)
	}

	mapCfg := structs.Map(cfgStr)

	for _, extDir := range dirs {
		metadataInfo = merge(metadataInfo, a.convToMetadataInfo(extensions[extDir.Name()], filepath.Join(extRootDir, extDir.Name()), mapCfg))
	}

	metadataInfo = groupMetadata(metadataInfo)
	data, _ := json.Marshal(metadataInfo)
	return a.repo.SetMetadata(context.Background(), confID, data)
}

func groupMetadata(metadataInfo []*models.MetadataInfo) []*models.MetadataInfo {
	tmp := lo.GroupBy(metadataInfo, func(item *models.MetadataInfo) string {
		return item.Type
	})

	result := make([]*models.MetadataInfo, 0, len(tmp))
	for key, val := range tmp {
		result = append(result, &models.MetadataInfo{
			ObjectName: key,
			Type:       key,
			Children:   val,
			ID:         uuid.NewString(),
		})
	}

	return result
}

func (a *Analyzer1C) codeAnalyzing(confID int32) error {
	// получаем из БД структуру метаданных ранее проанализированную
	// нас интересуют только заимствованные объекты, будем проверять их

	data, err := a.repo.GetMetadata(context.Background(), confID)
	if err != nil {
		return errors.Wrap(err, "couldn't get the configuration structure")
	}

	var metadata []*models.MetadataInfo
	_ = json.Unmarshal(data, &metadata)

	err = recursionRead(metadata, func(md *models.MetadataInfo) error {
		if utils.PtrToVal(md.Borrowed) {
			for i, ext := range md.Extension {
				err := filepath.WalkDir(filepath.Join(ext.PathObject, md.ObjectName), a.walk(ext.ExtensionRootPath, md, md.Extension[i].ID))
				if err != nil {
					log.Printf("parse error %s: %v\n", md.ObjectName, err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	data, _ = json.Marshal(metadata)
	return a.repo.SetMetadata(context.Background(), confID, data)
}

func recursionRead(metadata []*models.MetadataInfo, f func(md *models.MetadataInfo) error) error {
	if len(metadata) == 0 {
		return nil
	}

	for _, md := range metadata {
		if err := recursionRead(md.Children, f); err != nil {
			return err
		}

		if err := f(md); err != nil {
			return err
		}
	}

	return nil
}

func (a *Analyzer1C) walk(extensionRootPath string, md *models.MetadataInfo, extID int32) func(s string, d fs.DirEntry, err error) error {
	return func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		if !d.IsDir() && filepath.Ext(s) == ".bsl" {
			ast, err := parse(s)
			if err != nil {
				return err
			}

			fExist := make(map[string]*models.ExtChanges)

			for _, stm := range ast.ModuleStatement.Body {
				f := astpkg.Cast[*astpkg.FunctionOrProcedure](stm)
				if len(utils.Opt(f).Directives) > 0 {
					for _, d := range f.Directives {
						if d.Src == "" {
							continue
						}

						info := models.FuncInfo{
							Directive: d.Name,
							Type:      "Functions",
							Name:      f.Name,
							ModuleKey: utils.Hash([]byte(strings.TrimPrefix(s, extensionRootPath) + f.Name)),
						}

						ext := models.ExtChanges{
							ID:           extID,
							FuncsChanges: []models.FuncInfo{info},
						}

						objectName := lo.If(f.Type == astpkg.PFTypeProcedure, "Процедура: ").Else("Функция: ") + d.Src
						if _, ok := fExist[d.Src]; !ok {
							newItem := md.Children.Find(objectName, "Functions")
							if newItem == nil {
								newItem = &models.MetadataInfo{
									ObjectName: objectName,
									Type:       "Functions",
									Borrowed:   utils.Ptr(true),
									ID:         uuid.NewString(),
									Extension:  []*models.ExtChanges{&ext},
								}

								md.Children = append(md.Children, newItem)
							} else {
								newItem.Extension = append(newItem.Extension, &ext)
							}

							_ = a.repo.SetCode(context.Background(), extID, info.ModuleKey, ast.PrintStatement(f)) // сохраняем в БД
							fExist[d.Src] = &ext
							continue
						}

						fExist[d.Src].FuncsChanges = append(fExist[d.Src].FuncsChanges, info)
					}
				}
			}
		}

		return nil
	}
}

func parse(fPath string) (*astpkg.AstNode, error) {
	f, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := transform.NewReader(f, unicode.BOMOverride(unicode.UTF8BOM.NewDecoder()))
	data, _ := io.ReadAll(reader)

	a := astpkg.NewAST(string(data))
	err = a.Parse()
	return a, err
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

func (a *Analyzer1C) convToMetadataInfo(extID int32, rootDir string, cfgStr map[string]interface{}) []*models.MetadataInfo {
	var metadataInfo []*models.MetadataInfo

	dirs, err := os.ReadDir(rootDir)
	if err != nil {
		log.Println(errors.Wrap(err, "ReadDir error"))
	}

	for _, file := range dirs {
		if file.IsDir() {
			files, _ := os.ReadDir(filepath.Join(rootDir, file.Name()))
			for _, fileConf := range files {
				if !fileConf.IsDir() && filepath.Ext(fileConf.Name()) == ".xml" {
					path := filepath.Join(rootDir, file.Name(), fileConf.Name())
					metadataInfo = append(metadataInfo, a.readExtensionObject(extID, rootDir, file.Name(), path, cfgStr))
				}
			}
		}
	}

	return metadataInfo
}

func (a *Analyzer1C) readExtensionObject(extID int32, rootDir, groupName, path string, cfgStr map[string]interface{}) *models.MetadataInfo {
	f, err := os.Open(path)
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

	// ExtendedConfigurationObject может не быть, скорец всего когда расширение создали на одной конфе, потом подключили к другой БД с такой же конфой
	// по этому сравниваем по имени
	elemName := doc.SelectElement("/MetaDataObject/*/Properties/Name/text()")

	objectName := utils.Ptr(utils.Opt[xmlquery.Node](elemName)).Data
	return &models.MetadataInfo{
		ObjectName: objectName,
		Type:       groupName,
		Borrowed:   utils.Ptr(a.isBorrowed(cfgStr, groupName, objectName)),
		ID:         uuid.NewString(),
		Extension:  []*models.ExtChanges{changes(extID, rootDir, filepath.Join(rootDir, groupName), doc)},
	}
}

func (a *Analyzer1C) isBorrowed(cfgStr map[string]interface{}, groupName, objectName string) bool {
	if v, ok := cfgStr[groupName]; ok {
		for _, item := range utils.Cast[[]string](v) {
			if item == objectName {
				return true
			}
		}
	}

	return false
}

// changes информация что изменено
func changes(extID int32, rootPath, pathObject string, doc *xmlquery.Node) *models.ExtChanges {
	result := models.ExtChanges{ID: extID, PathObject: pathObject, ExtensionRootPath: rootPath}

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

			result.MetadataChanges = append(result.MetadataChanges, txt.String())
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

			result.MetadataChanges = append(result.MetadataChanges, txt.String())
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
			f, err := os.Open(filepath.Join(pathObject, fmt.Sprintf("%s\\Forms\\%s.xml", utils.Opt(elemName).Data, form.Data)))

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
			result.MetadataChanges = append(result.MetadataChanges, txt.String())
		}
	}

	return &result
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

func merge(metadataInfoPrev, metadataInfoNew []*models.MetadataInfo) []*models.MetadataInfo {
	for _, md := range metadataInfoNew {
		exist := false
		for j, mdPrev := range metadataInfoPrev {
			if md.Type == mdPrev.Type && md.ObjectName == mdPrev.ObjectName {
				metadataInfoPrev[j].Extension = append(metadataInfoPrev[j].Extension, md.Extension...)
				if md.Borrowed != nil {
					metadataInfoPrev[j].Borrowed = utils.Ptr(utils.PtrToVal(metadataInfoPrev[j].Borrowed) || utils.PtrToVal(md.Borrowed))
				}
				exist = true
				break
			}
		}

		if !exist {
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

	f.Close()
	childObjects, _ := parseConfigurationFile(f.Name())
	return &ConfigurationInfo{
		Name:         utils.Ptr(utils.Opt[etree.Element](elemName)).Text(),
		Synonym:      utils.Ptr(utils.Opt[etree.Element](elemSynonym)).Text(),
		Version:      utils.Ptr(utils.Opt[etree.Element](elemVersion)).Text(),
		Vendor:       utils.Ptr(utils.Opt[etree.Element](elemVendor)).Text(),
		Purpose:      utils.Ptr(utils.Opt[etree.Element](elemPurpose)).Text(),
		ChildObjects: childObjects,
	}, nil
}
