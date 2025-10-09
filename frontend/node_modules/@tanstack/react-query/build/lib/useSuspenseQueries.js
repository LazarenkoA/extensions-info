'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var useQueries = require('./useQueries.js');

function useSuspenseQueries({
  queries,
  context
}) {
  return useQueries.useQueries({
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

exports.useSuspenseQueries = useSuspenseQueries;
//# sourceMappingURL=useSuspenseQueries.js.map
