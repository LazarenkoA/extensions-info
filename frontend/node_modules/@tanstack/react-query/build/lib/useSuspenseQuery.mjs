import { QueryObserver } from '@tanstack/query-core';
import { useBaseQuery } from './useBaseQuery.mjs';

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
//# sourceMappingURL=useSuspenseQuery.mjs.map
