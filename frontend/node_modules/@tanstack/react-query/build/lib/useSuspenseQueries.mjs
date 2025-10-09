import { useQueries } from './useQueries.mjs';

function useSuspenseQueries({
  queries,
  context
}) {
  return useQueries({
    queries: queries.map(query => ({ ...query,
      enabled: true,
      useErrorBoundary: true,
      suspense: true,
      placeholderData: undefined,
      networkMode: 'always'
    })),
    context
  });
}

export { useSuspenseQueries };
//# sourceMappingURL=useSuspenseQueries.mjs.map
