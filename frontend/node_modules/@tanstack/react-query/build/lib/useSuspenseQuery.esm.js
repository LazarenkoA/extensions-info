import { QueryObserver } from '@tanstack/query-core';
import { useBaseQuery } from './useBaseQuery.esm.js';

function useSuspenseQuery(options) {
  return useBaseQuery({ ...options,
    enabled: true,
    useErrorBoundary: true,
    suspense: true,
    placeholderData: undefined,
    networkMode: 'always',
    onSuccess: undefined,
    onError: undefined,
    onSettled: undefined
  }, QueryObserver);
}

export { useSuspenseQuery };
//# sourceMappingURL=useSuspenseQuery.esm.js.map
