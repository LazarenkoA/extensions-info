import React from "react";
import {useWebSocket} from '../services/contexWebSocket'

const { useState, useEffect, useRef } = React;


const AnalysisView = ({ database, logs, progress, error }) => {
    const [showLogs, setShowLogs] = useState(false);
    const logEndRef = useRef(null);

    useEffect(() => {
        if (logEndRef.current) {
            logEndRef.current.scrollIntoView({ behavior: 'smooth' });
        }
    }, [logs]);

    if (error) {
        return (
            <div className="empty-state">
                <div className="empty-state-icon">⛔</div>
                <h2>Произошла ошибка</h2>
                <p>{error}</p>
            </div>
        );
    }

    return (
        <div className="analysis-view">
            <div className="analysis-status">
                <div className="spinner"></div>
                <div className="status-text">Анализ базы: {database.Name}</div>
                <div className="substatus-text">{logs[logs.length-1]}</div>
                <div className="substatus-text">Прогресс: {progress}%</div>
            </div>

            <div className="log-section">
                <div className="card">
                    <button
                        className="btn btn--secondary log-toggle"
                        onClick={() => setShowLogs(!showLogs)}
                    >
                        {showLogs ? 'Скрыть лог анализа' : 'Показать лог анализа'}
                    </button>

                    {showLogs && (
                        <div className="log-viewer">
                            {logs?.map((log, index) => (
                                <div key={index} className="log-entry">
                                    {log}
                                </div>
                            ))}
                            <div ref={logEndRef} />
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default AnalysisView;