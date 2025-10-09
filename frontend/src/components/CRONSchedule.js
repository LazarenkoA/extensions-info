import React from "react";
import { setCronSettings, deleteCronSettings } from '../services/cron';
import {useMutation, useQueryClient} from '@tanstack/react-query';

const cronPresets = [
    {"name": "Ежедневно в 2:00", "expression": "0 2 * * *"},
    {"name": "Еженедельно (понедельник)", "expression": "0 0 * * 1"},
    {"name": "Ежемесячно (1 число)", "expression": "0 0 1 * *"},
    {"name": "Каждые 6 часов", "expression": "0 */6 * * *"},
    {"name": "Рабочие дни в 9:00", "expression": "0 9 * * 1-5"}
]

function validateCron(cron) {
    const parts = cron.trim().split(/\s+/);
    if (parts.length !== 5) return false;

    const ranges = [
        [0, 59],  // minute
        [0, 23],  // hour
        [1, 31],  // day of month
        [1, 12],  // month
        [0, 6]    // day of week
    ];

    const validatePart = (part, min, max) => {
        const elements = part.split(',');
        for (let el of elements) {
            el = el.trim();
            if (el === '*') continue;

            // Шаги вида */n
            if (/^\*\/\d+$/.test(el)) {
                continue;
            }

            // Диапазон или одиночное число
            const rangeMatch = el.match(/^(\d+)(-(\d+))?$/);
            if (!rangeMatch) return false;

            const start = parseInt(rangeMatch[1], 10);
            const end = rangeMatch[3] ? parseInt(rangeMatch[3], 10) : start;

            if (start < min || end > max || start > end) return false;
        }
        return true;
    };

    return parts.every((part, i) => validatePart(part, ranges[i][0], ranges[i][1]));
}

function useAddMutation() {
    const queryClient = useQueryClient();

   return useMutation({
        mutationFn: setCronSettings,
        onError: (error) => {
            console.log(error)
        },
        onSuccess: () => {
            // обновим кэш
            queryClient.invalidateQueries({
                queryKey: ['databases']
            }).then(r => {});
        },
    });
}

function useDeleteMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (id) => deleteCronSettings(id),
        onError: (error) => {
            console.log(error)
        },
        onSuccess: () => {
            // после успешного удаления — обновим кэш
            queryClient.invalidateQueries({
                queryKey: ['databases']
            }).then(r => {});
        },
    });
}

const ScheduleEditorModal = ({ onClose, database }) => {
    const {useState, useEffect, useRef} = React;
    const [autoAnalysisEnabled, setAutoAnalysisEnabled] = useState(database?.Cron != undefined); // включает чекбокс "Включить автоматический анализ"
    const [cronExpression, setCronExpression] = useState(database?.Cron?.Schedule || '');
    const [selectedPreset, setSelectedPreset] = useState('');
    const [validationMessage, setValidationMessage] = useState('');
    const setSchedule = useAddMutation()
    const delSchedule = useDeleteMutation()

    useEffect(() => {
        if (database) {
            setAutoAnalysisEnabled(database?.Cron != undefined);
            setCronExpression(database?.Cron?.Schedule || '');
        }
    }, [database]);


    const validateCronExpression = (expression) => {
        if (!expression.trim()) {
            setValidationMessage('');
            return false;
        }

        const parts = expression.trim().split(/\s+/);
        if (!validateCron(expression)) {
            setValidationMessage('CRON выражение не вальдно');
            return false;
        }

        setValidationMessage('Валидное CRON выражение');
        return true;
    };

    const handleCronChange = (value) => {
        setCronExpression(value);
        validateCronExpression(value);
        setSelectedPreset('');
    };

    const handlePresetSelect = (preset) => {
        setCronExpression(preset.expression);
        setSelectedPreset(preset.expression);
        validateCronExpression(preset.expression);
    };

    const handleSave = () => {
        if (!validateCronExpression(cronExpression)) {
            return;
        }

        if(!autoAnalysisEnabled){
            console.log('Удаляем расписания для базы:', database.ID, {
                cronExpression
            });

            delSchedule.mutate(database.ID)
        } else {
            console.log('Сохранение расписания для базы:', database.ID, {
                cronExpression
            });

            setSchedule.mutate( { db_id: database.ID, data: {Schedule: cronExpression} });
        }
        onClose();
    };

     if (!database) return null;

    return (
        <div className="modal">
            <div className="modal-content">
                <div className="modal-header">
                    <h2 className="modal-title">Настройка расписания анализа</h2>
                    <button className="close-btn" onClick={onClose}>×</button>
                </div>

                <div className="modal-body">
                    <div className="database-info">
                        <strong>База данных:</strong> {database.Name}
                    </div>

                    <div className="checkbox-group">
                        <input
                            type="checkbox"
                            id="autoAnalysis"
                            checked={autoAnalysisEnabled}
                            onChange={(e) => setAutoAnalysisEnabled(e.target.checked)}
                        />
                        <label htmlFor="autoAnalysis">Включить автоматический анализ</label>
                    </div>

                    {autoAnalysisEnabled && (
                        <div className="cron-section">
                            <div className="form-group">
                                <label className="form-label">CRON выражение</label>
                                <input
                                    type="text"
                                    className="form-control"
                                    value={cronExpression}
                                    onChange={(e) => handleCronChange(e.target.value)}
                                    placeholder="0 2 * * * (ежедневно в 2:00)"
                                />
                                {validationMessage && (
                                    <div className={`cron-validation ${validationMessage.includes('Валидное') ? 'valid' : 'invalid'}`}>
                                        {validationMessage}
                                    </div>
                                )}
                            </div>

                            <div className="form-group">
                                <label className="form-label">Готовые шаблоны</label>
                                <div className="cron-presets">
                                    {cronPresets.map((preset, index) => (
                                        <button
                                            key={index}
                                            className={`preset-btn ${selectedPreset === preset.expression ? 'selected' : ''}`}
                                            onClick={() => handlePresetSelect(preset)}
                                        >
                                            <div>{preset.name}</div>
                                            <div style={{ fontSize: '10px' }}>
                                                {preset.expression}
                                            </div>
                                        </button>
                                    ))}
                                </div>
                            </div>
                        </div>
                    )}
                </div>

                <div className="modal-footer">
                    <button className="btn btn--secondary" onClick={onClose}>
                        Отмена
                    </button>
                    <button className="btn btn--primary" onClick={handleSave}>
                        Сохранить
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ScheduleEditorModal;