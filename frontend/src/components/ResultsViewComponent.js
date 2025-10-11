import TreeNode from './TreeNodeComponent';
import {DetailsPanel} from './TreeNodeComponent'
import React from "react";

const { useState, useEffect, useRef } = React;

const ResultsView = ({ conf, database }) => {
    const [selectedItem, setSelectedItem] = useState(null);

    if (!conf || !conf.Extensions) {
        return (
            <div className="empty-state">
                <div className="empty-state-icon">📊</div>
                <h2>Результаты анализа недоступны</h2>
                <p>Данные анализа для этой базы еще не готовы.</p>
            </div>
        );
    }

    return (
        <div className="results-view">
            <div className="results-header">
                <h1>Результаты анализа: {database.Name}</h1>
                <div className="completion-time">
                    Анализ завершен: {database.LastCheckAsString}
                </div>
            </div>
            <div className="results-grid">
                <div className="result-card config-info">
                    <h3>Информация о конфигурации</h3>
                    <div className="config-details">
                        <div className="config-detail">
                            <div className="config-detail-label">Название</div>
                            <div className="config-detail-value">{conf.Name}</div>
                        </div>
                        <div className="config-detail">
                            <div className="config-detail-label">Версия</div>
                            <div className="config-detail-value">{conf.Version}</div>
                        </div>
                        <div className="config-detail">
                            <div className="config-detail-label">Расширений</div>
                            <div className="config-detail-value">{conf?.Extensions?.length}</div>
                        </div>
                    </div>
                    <div>
                        <div className="config-detail-label">Установленные расширения:</div>
                        <div className="extensions-list">
                            {conf?.Extensions?.map((ext, index) => (
                                <span key={index} className="extension-tag">{ext.Name} {ext.Version != '' ? `(v: ${ext.Version})` : ''}</span>
                            ))}
                        </div>
                    </div>
                </div>
                <div className="result-card config-info">
                    <div className="metadata-container">
                        <div className="tree-container">
                            <h3>Переопределенные метаданные</h3>
                            {conf.MetadataTree && (
                                <TreeNode
                                    node={conf.MetadataTree}
                                    setSelectedItem={setSelectedItem}
                                    selectedItem={selectedItem}
                                />)}
                        </div>

                        {selectedItem && (
                            <DetailsPanel
                                selectedItem={selectedItem}
                                extensionInfo={conf.Extensions}
                            />
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ResultsView;