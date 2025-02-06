import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export const fetchPings = async () => {
  const response = await axios.get(`${API_BASE_URL}/api/pings`);
  return response.data;
};
