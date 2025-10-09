import {apiRequest} from './request'

// Получить настройки
export const getAppSettings = () => {
    return apiRequest(`/getAppSettings`);
};

// Сохранить настройки
export const storeAppSettings = ({id, data}) => {
    return apiRequest(`/appSettings/${id}`, { method: 'POST', body: JSON.stringify(data)});
};