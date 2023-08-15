import { useLocation } from 'react-router-dom';
import React from 'react';

const localHardCoded = undefined;

export function getBaseUrl() {
  return (
    localHardCoded || window.location.protocol + '//' + window.location.host
  );
}

export function useQuery() {
  const { search } = useLocation();
  return React.useMemo(() => new URLSearchParams(search), [search]);
}
