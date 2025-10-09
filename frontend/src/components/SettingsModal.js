import React from "react";
import { getAppSettings, storeAppSettings } from '../services/appSettings';
import {useMutation, useQuery} from '@tanstack/react-query';
import {addBaseSettings} from "../services/settingservice";

const { useState, useEffect, useRef } = React;
const useLoadSettings = (setFormData) => {
    return useQuery({
        queryKey: ['settings'],
        queryFn: getAppSettings,
        select: (data) => data.data,
        refetchOnMount: 'always',
    });
}

const SettingsModal = ({onClose}) => {
    const [formData, setFormData] = useState(null);
    const [handleSubmitError, setError] = useState('');
    const addSettings = useMutation({
        mutationFn: storeAppSettings
    });
    const {data: settings, isLoading, error } = useLoadSettings();
    if (error) {
        console.log(error)
    }

    useEffect(() => {
        if (settings) {
            setFormData(settings);
        }
    }, [settings]);


    const handleSubmit = (e) => {
        e.preventDefault();
        console.log('Сохранение глобальных настроек:');

        addSettings.mutate( { id: formData.ID, data: formData }, {
            onError: (err) => setError('Error: '+err.message),
            onSuccess: () => { setError('');  onClose()}
        });
    };

    return (
        <div className="modal modal-medium">
            <div className="modal-content">
                <div className="modal-header">
                    <h2 className="modal-title">Глобальные настройки</h2>
                    <button className="close-btn" onClick={onClose}>×</button>
                </div>
                <form onSubmit={handleSubmit}>
                    <div className="modal-body">
                        <div className="settings-form">
                            <div className="settings-section">
                                <h3>Путь к платформе 1С</h3>
                                <div className="form-group">
                                    <label className="form-label">Исполняемый файл платформы</label>
                                    <input
                                        type="text"
                                        className="form-control"
                                        value={formData?.PlatformPath}
                                        onChange={(e) => setFormData({...formData, PlatformPath: e.target.value})}
                                    />
                                </div>
                            </div>

                            {/*<div className="settings-section">*/}
                            {/*    <h3>Настройки подключения</h3>*/}
                            {/*    <div className="form-row">*/}
                            {/*        <div className="form-group">*/}
                            {/*            <label className="form-label">Таймаут подключения (сек)</label>*/}
                            {/*            <input*/}
                            {/*                type="number"*/}
                            {/*                className="form-control"*/}
                            {/*                value={formData.connectionTimeout}*/}
                            {/*                min="10"*/}
                            {/*                max="120"*/}
                            {/*            />*/}
                            {/*        </div>*/}
                            {/*        <div className="form-group">*/}
                            {/*            <label className="form-label">Таймаут анализа (сек)</label>*/}
                            {/*            <input*/}
                            {/*                type="number"*/}
                            {/*                className="form-control"*/}
                            {/*                value={formData.analysisTimeout}*/}
                            {/*                min="60"*/}
                            {/*                max="1800"*/}
                            {/*            />*/}
                            {/*        </div>*/}
                            {/*    </div>*/}
                            {/*</div>*/}
                        </div>
                    </div>
                    <div className="modal-footer">
                        <button type="button" className="btn btn--secondary" onClick={onClose}>
                            Отмена
                        </button>
                        <button type="submit" className="btn btn--primary">
                            Сохранить
                        </button>
                    </div>
                </form>
                <br/>
                {handleSubmitError && <div style={{
                    color: 'var(--color-error)',
                    marginBottom: '16px',
                    fontSize: '14px'
                }}>{handleSubmitError}</div>}
            </div>
        </div>
    )
}

export default SettingsModal;