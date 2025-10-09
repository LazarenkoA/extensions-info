
import {apiRequest} from './request'

// Получить настройки по базам
export const getBaseSettings = () => {
  return apiRequest('/getBaseSettings');
};

// Добавить новую базу
export const addBaseSettings = (data) => {
  return apiRequest('/addBaseSettings', { method: 'POST', body: JSON.stringify(data)});
};

// Удалить базу
export const deleteBaseSettings = (id) => {
  return apiRequest(`/deleteBaseSettings/${id}`, { method: 'DELETE'});
};

