import React from "react";

const { useState, useEffect, useRef } = React;

export const DetailsPanel = ({ selectedItem, extensions, onFunctionClick }) => {
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


    const renderMetadataDetails = (objectData) => {
        return (
            <div className="details-panel">
                <div className="details-header">
                    <h3 className="details-title">Переопределения: {objectData.ObjectName}</h3>
                </div>
                <div className="details-content">
                    {extensions && extensions.map((ext, index) => (
                        <div key={index} className="extension-section">
                            <div className="extension-header">
                                <h4 className="extension-name">{ext.Name} (v{ext.Version})</h4>
                            </div>
                            <div className="extension-body">
                            <ul className="changes-list">
                                    {/*{ext.changes.map((change, changeIndex) => (*/}
                                    {/*    <li key={changeIndex}>{change}</li>*/}
                                    {/*))}*/}
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
            switch (mode) {
                case 'Вместо': return 'vmestore';
                case 'ИзменениеИКонтроль': return 'izmenenie-i-kontrol';
                case 'Перед': return 'pered';
                case 'После': return 'posle';
                default: return 'posle';
            }
        };

        return (
            <div className="details-panel">
                <div className="details-header">
                    <h3 className="details-title">Переопределения функции: {functionData.Name}</h3>
                </div>
                <div className="details-content">
                    {extensions && extensions.map((ext, index) => (
                        <div key={index} className="extension-section">
                            <div className="extension-header">
                                <h4 className="extension-name">{ext.Name} (v{ext.Version})</h4>
                            </div>
                            <div className="extension-body">
                                <div className="function-override">
                                    <div className="override-header">
                    <span className={`override-type-badge ${getOverrideTypeClass(functionData.RedefinitionMethod)}`}>
                      {functionData.RedefinitionMethod}
                    </span>
                                        <div
                                            className="function-name-link"
                                            // onClick={() => onFunctionClick(ext)}
                                        >
                                            /api/v1/getConfiguration77777777777777777777Info?id=1
                                            {/*{ext.functionName}*/}
                                        </div>
                                    </div>

                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        );
    };

    // Determine if it's a metadata object or function
    // const objectData = database.objectOverrides[selectedItem.id];
    // const functionData = database.functionOverrides[selectedItem.id];

    if (selectedItem.Type === 'function') {
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
    const handleToggle = () => {
        setExpanded(!expanded);
    };

    if(!node) {
        return null
    }


    const hasChildren = node.Children && node.Children.length > 0;
    const hasFunctions = node.Funcs && node.Funcs.length > 0;
    //const hasModules = node.modules && node.modules.length > 0;

    const icon = {
        'function': '📐',
        'folder': '📁',
        'commonModule': '⚙️',
        'document': '📑',
        'catalog': '📝',
        'configuration': '🧩'
    }

    expanded = expanded || level == 0 // первый уровент всегда открыт

    return (
        <div>
            <div className="tree-node">
                <div
                    className={`tree-item ${expanded ? 'expanded' : ''} ${selectedItem ? 'selected' : ''}`}
                    onClick={() => setSelectedItem(node)}
                >
                    <div className="tree-toggle" onClick={handleToggle}>
                        {(hasChildren || hasFunctions) ? (expanded ? '▼' : '▶') : ''}
                    </div>
                    <div className={`tree-icon`}>{icon[node.Type]}</div>
                    <div className="tree-label">
                        {node.ObjectName || node.Name || node.Type}
                    </div>
                    {node.status && (
                        <div className={`tree-status ${node.status}`}>
                            {node.status === 'modified' ? 'Изменен' :
                                node.status === 'added' ? 'Добавлен' : 'Удален'}
                        </div>
                    )}
                </div>

                {expanded  && (
                    <div className="tree-children">
                        {hasChildren && node.Children.map((child, index) => (
                            <TreeNode
                                key={index}
                                node={child}
                                level={level + 1}
                                selectedItem={selectedItem}
                                setSelectedItem={setSelectedItem}
                            />
                        ))}
                        {/*{hasModules && node.modules.map((module, index) => (*/}
                        {/*    <TreeNode*/}
                        {/*        key={index}*/}
                        {/*        node={module}*/}
                        {/*        level={level + 1}*/}
                        {/*        type="module"*/}
                        {/*    />*/}
                        {/*))}*/}
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
        </div>
    );

};

export default TreeNode;