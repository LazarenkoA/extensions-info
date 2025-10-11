import React from "react";
import { useDatabase } from "./ContexDatabaset";
import {getConfigurationInfo } from '../services/configuration';
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query';
import AnalysisView from "./AnalysisViewComponent";
import ResultsView from "./ResultsViewComponent";
import {BDStates} from "./States";
import {useWebSocket} from "../services/contexWebSocket";
import {startManualAnalysis} from "../services/cron";

const { useState, useEffect } = React;

export const useStartAnalysis = (databaseID) => {
    const [error, setError] = useState(null);
    const manualAnalysis = useMutation({
        mutationFn: (id) => startManualAnalysis(id),
        onSuccess: () => {
            setError(null)
        },
        onError: (e) => {
            setError(e);
            console.log(e);
        }
    });

    return {
        start: (e) => {
            e.stopPropagation();
            manualAnalysis.mutate(databaseID);
        },
        error: error,
    }
};

// отлов событий смены статуса по ws
function useChangeState() {
    const { subscribe } = useWebSocket();
    const [state, setState] = useState({});

    useEffect(() => {
        try {
            const unsubscribe = subscribe((msg) => {
                const obj = JSON.parse(msg);

                // фильтруем по типу
                if (obj.type === "new_state") {
                    setState({db: obj.db_id, state: obj.msg});
                }
            });

            return unsubscribe
        } catch (e) {
            console.error("Invalid WS message");
        }
    }, [subscribe]);

    return state;
}

// отлов логов по ws
export function useLogs(dbId) {
    const { subscribe } = useWebSocket();
    const [logs, setLogs] = useState([]);
    const [progress, setProgress] = useState(0);

    useEffect(() => {
        try {
            const unsubscribe = subscribe((msg) => {
                const obj = JSON.parse(msg);

                // фильтруем по типу и ID
                if (obj.type === "log" && obj.db_id === dbId) {
                    setLogs((prev) => [...prev, `${obj.time}: ${obj.msg}`]);
                }
                if (obj.type === "progress" && obj.db_id === dbId) {
                    setProgress(obj.msg);
                }
            });

            return unsubscribe
        } catch (e) {
            console.error("Invalid WS message");
        }
    }, [dbId, subscribe]);

    return {logs, progress, setLogs, setProgress};
}

const MainContent = () => {
    const { selectedDb, setSelectedDb } = useDatabase();
    const queryClient = useQueryClient();
    const state = useChangeState();
    const {logs,progress,setLogs,setProgress} = useLogs(selectedDb?.ID)
    const {start: startAnalysis, error: errorAnalysis} = useStartAnalysis(selectedDb?.ID);

    // при смене статуса обновляем список БД
    useEffect(() => {
        if (state.state === BDStates.DONE) {
            setLogs([])
        }

        setProgress(0)
        queryClient.invalidateQueries({
            queryKey: ['databases']
        }).then(r => {
            if(selectedDb?.ID === state.db)
                setSelectedDb({...selectedDb, Status: state.state})
        });
    }, [state, queryClient]);


    const query = useQuery({
        queryKey: ['Configuration', selectedDb?.ID],
        //staleTime: 10_000, // кеш 10 секунд
        queryFn: ({queryKey }) => {
            const [, id] = queryKey;
            return getConfigurationInfo(id)
        } ,
        select: (data) => data.data,
        enabled: !!selectedDb?.ID, // запрос выполняется только если selectedDb есть
    });
    useEffect(() => {
        if (selectedDb?.ID) {
            query.refetch().then(()=>{}); // рефрешем при смене БД
        }
    }, [selectedDb]);

    if(query.error) {
        console.log(query.error)
        return null
    }


    let analyseError = null
    if(selectedDb?.Status === BDStates.ERROR)
        analyseError = logs[logs.length-1];

    if (selectedDb?.Status === BDStates.ANALYZING || analyseError) {
        return (
            <AnalysisView
                database={selectedDb}
                progress={progress}
                logs={logs}
                error={analyseError}
            />
        );
    }

    if (selectedDb?.Status === BDStates.DONE) {
        return <ResultsView conf={query.data} database={selectedDb} />;
    }

    return (
        <div className="empty-state">
            <div className="empty-state-icon">⚡</div>

            {!selectedDb && (<p>Выберите базу данных из списка слева</p>)}
            {selectedDb && (
                <div>
                    <h2>{selectedDb?.Name}</h2>
                    <p>База данных готова к анализу</p>
                </div>
            )}

            {selectedDb?.Status === BDStates.NEW && (
                <button
                    className="btn btn--primary btn--lg"
                    onClick={(e) => {
                        setProgress(0)
                        startAnalysis(e)
                    }}
                >
                    Начать анализ
                </button>
            )}
        </div>
    );
}

export default MainContent;