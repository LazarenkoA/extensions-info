import {apiRequest} from './request'

// Получить информацию по конфигурации
export const getConfigurationInfo = (id) => {
    return apiRequest(`/getConfigurationInfo?id=${id}`);
};