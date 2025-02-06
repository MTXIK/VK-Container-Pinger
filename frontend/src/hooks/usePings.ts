import { useEffect, useState } from 'react';
import { fetchPings } from '../services/api';

export interface Ping {
  id?: number;
  ip_address: string;
  container_name: string;
  ping_time: number;
  last_success: string;
}

export const usePings = (interval: number = 10000) => {
  const [pings, setPings] = useState<Ping[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const loadPings = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchPings();
      setPings(data);
    } catch (err: any) {
      setError(err.message || 'Ошибка при загрузке данных');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPings();
    const timer = setInterval(loadPings, interval);
    return () => clearInterval(timer);
  }, [interval]);

  return { pings, loading, error, reload: loadPings };
};
