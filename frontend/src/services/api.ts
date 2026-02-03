import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'https://github-resume-builder-eihv.onrender.com';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authService = {
  login: () => {
    window.location.href = `${API_URL}/auth/login`;
  },
  
  setToken: (token: string) => {
    localStorage.setItem('token', token);
  },
  
  getToken: () => {
    return localStorage.getItem('token');
  },
  
  logout: () => {
    localStorage.removeItem('token');
  },
};

export const resumeService = {
  generate: async (targetRole: string) => {
    const response = await api.post('/resumes/generate', { target_role: targetRole });
    return response.data;
  },
  
  list: async () => {
    const response = await api.get('/resumes');
    return response.data;
  },
  
  get: async (id: number) => {
    const response = await api.get(`/resumes/${id}`);
    return response.data;
  },
  
  update: async (id: number, data: any) => {
    const response = await api.put(`/resumes/${id}`, data);
    return response.data;
  },
  
  delete: async (id: number) => {
    await api.delete(`/resumes/${id}`);
  },
};

export default api;
