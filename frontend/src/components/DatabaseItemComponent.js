import { getBaseSettings, deleteBaseSettings } from '../services/settingservice';
import { useQuery } from '@tanstack/react-query';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useDatabase } from "./ContexDatabaset";
import React from "react";
import ScheduleEditorModal from "./CRONSchedule";
import {BDStates} from "./States";
import {useStartAnalysis, useLogs} from "./MainContent";


const { useState, useEffect, useRef } = React;

const getStatusText = (status) => {
    switch (status) {
        case BDStates.NEW: return 'Новая';
        case BDStates.ANALYZING: return 'Анализ';
        case BDStates.DONE: return 'Готово';
        case BDStates.ERROR: return 'Ошибка';
        default: return 'Неизвестно';
    }
};

function useDeleteBaseSettings() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (id) => deleteBaseSettings(id),
        onSuccess: () => {
            // после успешного удаления — обновим кэш
            queryClient.invalidateQueries({
                queryKey: ['databases']
            }).then(r => {});
        },
        onError: (e) => {
            console.log(e)
        }
    });
}

function useDatabaseList() {
    return useQuery({
        queryKey: ['databases'],
        staleTime: 30_000, // кеш 30 секунд
        queryFn: getBaseSettings,
        select: (data) => data.data,
    });
}

const DatabaseList = () => {
    const { selectedDb } = useDatabase();
    const { data: databases, isLoading, error } = useDatabaseList()
    if(error) {
        console.log(error)
    }

   const [activeDatabaseId, setActiveDatabaseId] = useState(null);

    useEffect(() => {
        if (databases && databases.length > 0) {
            setActiveDatabaseId(selectedDb?.ID || databases[0].ID);
        }
    }, [databases]); // срабатывает когда databases загружаются

    return(
        <div className="database-list">
            { databases && databases.map((db) => (
                    <DatabaseItem
                        key={db.ID}
                        database={db}
                        isActive={db.ID === activeDatabaseId}
                        setActiveDatabaseId={(dbID) => setActiveDatabaseId(dbID)}
                     />
                ))
            }
        </div>
    )
}

const DatabaseItem = ({database, isActive, setActiveDatabaseId}) => {
    const [scheduleModalIsShow, setShowScheduleModal] = useState(false);
    const { selectedDb, setSelectedDb } = useDatabase();
    const deleteMutation = useDeleteBaseSettings();
    const {progress: progress, setProgress: setProgress, setLogs} = useLogs(database.ID)
    const onClick = (database) => {
        setActiveDatabaseId(database.ID);
        setSelectedDb(database);
    }
    const {start: startAnalysis, error: errorAnalysis} = useStartAnalysis(selectedDb?.ID);

    return (
        <div
            className={`database-item ${isActive ? 'active' : ''}`}
            onClick={() => onClick(database)}
        >
            <div className="database-actions">
                <button
                    className="close-btn"
                    onClick={(e) => {
                        const confirmDelete = window.confirm("Вы уверены, что хотите удалить запись?");
                        if (!confirmDelete) {
                            return
                        }

                        e.stopPropagation();
                        deleteMutation.mutate(database.ID);
                    }}
                    title="Удалить базу"
                >
                    ×
                </button>
            </div>

            <div className="database-item-header">
                <div className="database-name">{database.Name}</div>
                <div className="database-status">
                    <div className={`status-indicator ${database.Status}`}></div>
                    <span>{getStatusText(database.Status)}</span>
                </div>
            </div>

            {!database.Cron && (
                <div className="schedule-status">
                <div className="schedule-icon">⏸️</div>
                <span>Только ручной запуск</span></div>)}

            <div className="database-details">
                {database.ConnectionString}
                {database.LastCheckAsString && (
                    <div>Последний анализ: {database.LastCheckAsString}</div>
                )}
            </div>

            {database.Status === BDStates.ANALYZING && (
                <div className="database-progress">
                    <div
                        className="progress-fill"
                        style={{width: `${progress||1}%`}}
                    ></div>
                </div>
            )}

            {database.Status !== BDStates.ANALYZING && (
                <div style={{display: 'flex', gap: '8px'}}>
                    <button
                        className="btn btn--primary btn--sm"
                        style={{marginTop: '8px', width: '100%'}}
                        onClick={() => setShowScheduleModal(true)}
                    >
                        Расписание
                    </button>
                    <button className="btn btn--primary btn--sm" db={database.ID}
                            onClick={(e) => {
                                setProgress(0)
                                if(e.target.attributes.db.value == selectedDb?.ID)  // защита что б можно было запустить только на активной БД
                                    startAnalysis(e);
                            } }
                            style={{marginTop: '8px', width: '100%'}}>
                        Начать анализ
                    </button>
                </div>
            )
            }

            {scheduleModalIsShow && (<ScheduleEditorModal
                database={database}
                onClose={() => setShowScheduleModal(false)}
            />)}
        </div>
    );
};

export default DatabaseList;
