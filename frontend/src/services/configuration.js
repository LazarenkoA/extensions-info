import {apiRequest} from './request'

// Получить информацию по конфигурации
export const getConfigurationInfo = (id) => {
    return apiRequest(`/getConfigurationInfo?id=${id}`);
};

// Получить информацию по коду
export const getSourceCode = (extid, modulekey) => {
    return apiRequest(`/getSourceCode?extid=${extid}&modulekey=${modulekey}`);
};