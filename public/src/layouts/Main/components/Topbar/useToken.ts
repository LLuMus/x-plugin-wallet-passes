import { useState } from 'react';

export interface Customer {
  id: string;
  email: string;
  details: string;
  credit_tokens: number;
  auth_code: string;
  updated_at: number;
  created_at: number;
}

export default function useToken() {
  const getToken = () => {
    return sessionStorage.getItem('token');
  };

  const [token, setToken] = useState(getToken());

  const saveToken = userToken => {
    sessionStorage.setItem('token', userToken);
    setToken(userToken);
  };

  return {
    setToken: saveToken,
    token
  };
}
