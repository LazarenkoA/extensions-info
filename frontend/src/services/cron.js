import {apiRequest} from './request'

// Получить настройки CRON
export const getCronSettings = () => {
    return apiRequest(`/getCronSettings`);
};

// Сохранить настройки CRON
export const setCronSettings = ({db_id, data}) => {
    return apiRequest(`/setCronSettings/${db_id}`, { method: 'POST', body: JSON.stringify(data)});
};

// Удаление настройки CRON
export const deleteCronSettings = (db_id) => {
    return apiRequest(`/deleteCronSettings/${db_id}`, { method: 'DELETE'});
};

// Ручной запуск анализа
export const startManualAnalysis = (db_id) => {
    return apiRequest(`/startManualAnalysis/${db_id}`, { method: 'POST'});
};



