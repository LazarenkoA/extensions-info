import React from "react";
import {useDatabase} from "./ContexDatabaset";
import { useQuery } from '@tanstack/react-query';
import hljs from 'highlight.js/lib/core';
import 'highlight.js/styles/github-dark.css';
import onec from 'highlight.js/lib/languages/1c';
import {getConfigurationInfo, getSourceCode} from "../services/configuration";

hljs.registerLanguage('1c', onec);

const ruNames = {
    'Configuration': {name:'Конфигурация', icon: 'root'},
    'Functions': {name:'Метод', icon: 'func'},
    'CommonModules': {name:'Общий модуль',icon: 'commonModule'},
    'Documents': {name:'Документ',icon: 'document'},
    'Catalogs': {name:'Справочник',icon: 'catalog'},
    'Constants': {name:'Константы',icon: 'сonstant'},
    'Roles': {name:'Роль',icon: 'role'},
    'Enums': {name:'Перечисление',icon: 'enum'},
    'Languages': {name:'Язык',icon: 'language'},
    'InformationRegisters': {name:'Регистр сведений',icon: 'informationRegister'},
    'HTTPServices': {name:'HTTP сервис',icon: 'httpService'},
    'WebServices': {name:'Web сервис',icon: 'webService'},
    'WSReferences': {name:'WS ссылка',icon: 'webService'},
    'XDTOPackages': {name:'XDTO пакет',icon: 'xdtoPackage'},
    'Reports': {name:'Отчет',icon: 'report'},
    'Subsystems': {name:'Подсистема',icon: 'subsystem'},
    'Styles': {name:'Стриль',icon: 'style'},
    'StyleItems': {name:'Элемент стиля',icon: 'style'},
    'CommonForms': {name:'Общая форма',icon: 'form'},
    'CommonCommands': {name:'Общая команда',icon: 'command'},
    'CommandGroups': {name:'Группа команд',icon: 'command'},
    'AccountingRegisters': {name:'Регистр бухгалтерии',icon: 'accountingRegister'},
    'AccumulationRegisters': {name:'Регистр накопления',icon: 'accumulationRegister'},
    'CalculationRegisters': {name:'Регистр расчета',icon: 'calculationRegister'},
    'CommonTemplates': {name:'Общий макет',icon: 'template'},
    'CommonPictures': {name:'Общая картинка',icon: 'commonPictures'},
    'DataProcessors': {name:'Обработка',icon: 'dataProcessors'},
    'SessionParameters': {name:'Параметр сеанса',icon: 'sessionParameter'},
    'SettingsStorages': {name:'Хранилище настроек',icon: 'settingsStorage'},
}

const { useState,  useEffect, useRef } = React;

const useHighlightTheme = () => {

        const updateTheme = () => {
            const colorScheme = document.documentElement.getAttribute("data-color-scheme");
            const existingLink = document.querySelector('link[href*="highlight.js"]');

            if (existingLink) {
                existingLink.remove();
            }

            const link = document.createElement('link');
            link.rel = 'stylesheet';

            const theme = colorScheme === 'dark' ? 'a11y-dark' : '1c-light';
            link.href = `https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/styles/${theme}.css`;

            link.onload = () => {
                // Убираем старые классы highlight.js
                document.querySelectorAll('pre code').forEach(block => {
                    block.removeAttribute('data-highlighted');
                    block.className = block.className.replace(/hljs[\w-]*/g, '').trim();
                });

                // Повторно применяем подсветку
                hljs.highlightAll();
            };

            document.head.appendChild(link);
        };

        // Обновляем тему при монтировании
        updateTheme();

        // Наблюдаем за изменениями атрибута data-color-scheme
        const observer = new MutationObserver(() => {
            updateTheme();
        });

        observer.observe(document.documentElement, {
            attributes: true,
            attributeFilter: ['data-color-scheme']
        });

        return () => observer.disconnect();

};

function useSourceCode(extID, moduleKey) {
    return useQuery({
        queryKey: ['sourceCode', extID, moduleKey],
        queryFn: ({queryKey} ) => {
            const [, extid, modulekey] = queryKey;
            return getSourceCode(extid, modulekey)
        } ,
        select: (data) => data.data,
        enabled: !!extID && !!moduleKey,
    });
}

const CodeModal = ({ onClose, functionData }) => {
    useHighlightTheme()

    const { data: code, isLoading, error } = useSourceCode(functionData?.extID, functionData?.moduleKey)
    if(error) {
        console.log(error)
    }

    if (!functionData) return null;

    return (
        <div className="code-modal">
            <div className="code-modal-content">
                <div className="code-modal-header">
                    <h3 className="code-modal-title">
                        {functionData.functionName}
                    </h3>
                    <button className="close-btn" onClick={onClose}>×</button>
                </div>
                <pre><code className="language-1c">{code}</code></pre>
            </div>
        </div>
    );
};

export const DetailsPanel = ({selectedItem, extensionInfo }) => {
    const [codeModal, setCodeModal] = useState(null);

    if (!selectedItem) {
        return (
            <div className="details-panel hidden">
                <div className="details-header">
                    <h3 className="details-title">Детали</h3>
                </div>
                <div className="details-content">
                    <div className="empty-state">
                        <div className="empty-state-icon">📋</div>
                        <h4>Выберите элемент</h4>
                        <p>Выберите объект метаданных или функцию для просмотра переопределений</p>
                    </div>
                </div>
            </div>
        );
    }

    const extPreview = (id) => {
        const ext = extensionInfo?.filter(ext => ext.ID === id)
        if(ext?.length === 1)
            return `${ext[0].Name} ${ext[0].Version !== '' ? `(v: ${ext[0].Version})` : ''}`
    }

    const renderMetadataDetails = (objectData) => {
        return (
            <div className="details-panel">
                <div className="details-header">
                    <h3 className="details-title">{ruNames[objectData.Type]?.name || objectData.Type}: {objectData.ObjectName}</h3>
                </div>
                <div className="details-content">
                    {objectData.Extension && objectData.Extension.map((ext, index) => (
                        <div key={index} className="extension-section">
                            <div className="extension-header">
                                <h4 className="extension-name">{extPreview(ext.ID)}</h4>
                            </div>
                            <div className="extension-body">
                                <ul className="changes-list">
                                    {ext.MetadataChanges && ext.MetadataChanges.map((change, changeIndex) => (
                                        <li key={changeIndex}>{change}</li>
                                    ))}
                                </ul>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        );
    };

    const renderFunctionDetails = (functionData) => {
        const getOverrideTypeClass = (mode) => {
            switch (mode.toLowerCase()) {
                case '&вместо': return 'vmestore';
                case '&изменениеиконтроль': return 'izmenenie-i-kontrol';
                case '&перед': return 'pered';
                case '&после': return 'posle';
                default: return '-';
            }
        };

        return (
            <div className="details-panel">
                <div className="details-header">
                    <h3 className="details-title">Переопределения функции: {functionData.Name}</h3>
                </div>
                <div className="details-content">
                    {functionData.Extension && functionData.Extension.map((ext, index) => (
                        <div key={index} className="extension-section">
                            <div className="extension-header">
                                <h4 className="extension-name">{extPreview(ext.ID)}</h4>
                            </div>
                            <div className="extension-body">
                                {functionData && ext.FuncsChanges && ext.FuncsChanges.map((func, id) => (
                                    <div key={id} className="function-override">
                                        <div className="override-header">
                    <span
                        className={`override-type-badge ${getOverrideTypeClass(func.Directive)}`}>
                      {func.Directive}
                    </span>
                                            <div
                                                className="function-name-link" data-tooltip={func.Name}
                                                onClick={() => setCodeModal({ isOpen: true, functionName: func.Name, extID: ext.ID, moduleKey: func.ModuleKey})}
                                            >
                                                {func.Name}
                                            </div>
                                        </div>

                                    </div>
                                ))}

                            </div>
                        </div>
                    ))}
                </div>
                {codeModal && codeModal.isOpen && <CodeModal
                    onClose={() => setCodeModal(null)}
                    functionData={codeModal}
                />}
            </div>
        );
    };

    if (selectedItem.Type === 'Functions') {
        return renderFunctionDetails(selectedItem);
    }

    return renderMetadataDetails(selectedItem);

    return (
        <div className="details-panel">
            <div className="details-header">
                <h3 className="details-title">Детали недоступны</h3>
            </div>
            <div className="details-content">
                <div className="empty-state">
                    <div className="empty-state-icon">❓</div>
                    <h4>Данные не найдены</h4>
                    <p>Информация о переопределениях для выбранного элемента недоступна</p>
                </div>
            </div>
        </div>
    );
};

const TreeNode = ({ node, level = 0, selectedItem, setSelectedItem }) => {
    let [expanded, setExpanded] = useState(false);
    const { selectedDb } = useDatabase();
    const handleToggle = () => {
        setExpanded(!expanded);
    };
    useEffect(() => {
        setSelectedItem(null)
    }, [selectedDb])

    if(!node) {
        return null
    }


    const hasChildren = node.Children && node.Children.length > 0;
    const hasFunctions = node.Funcs && node.Funcs.length > 0;

    expanded = expanded || level == 0 // первый уровент всегда открыт

    return (
            <div className="tree-node">
                <div
                    className={`tree-item ${expanded ? 'expanded' : ''} ${selectedItem && selectedItem.ID == node.ID ? 'selected' : ''}`}
                    onClick={() => setSelectedItem(node)}
                >
                    <div className="tree-toggle" onClick={handleToggle}>
                        {(hasChildren || hasFunctions) ? (expanded ? '▼' : '▶') : ''}
                    </div>
                    <div className={`icon-${ruNames[node.Type]?.icon || 'default'}`}></div>
                    <div className="tree-label">
                        {node.ObjectName || node.Name || node.Type}
                    </div>

                    {node.Borrowed != undefined && (<div className={`tree-status ${node.Borrowed ? 'modified':'added'}`}>
                        {node.Borrowed ? 'Изменен' : 'Добавлен' }
                    </div>)}

                </div>

                {expanded  && (
                    <div className="tree-children">
                        {hasChildren && node.Children?.sort((a, b) => b.Type.localeCompare(a.Type)).map((child, index) => (
                            <TreeNode
                                key={index}
                                node={child}
                                level={level + 1}
                                selectedItem={selectedItem}
                                setSelectedItem={setSelectedItem}
                            />
                        ))}
                        {hasFunctions && node.Funcs.map((func, index) => (
                            <TreeNode
                                key={index}
                                node={func}
                                level={level + 1}
                                selectedItem={selectedItem}
                                setSelectedItem={setSelectedItem}
                            />
                        ))}
                    </div>
                )}
            </div>
    );

};

export default TreeNode;