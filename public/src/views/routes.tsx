import React from 'react';

import {
  LegalPage as LegalPageView,
  DesktopApp as DesktopAppView,
  PassPage as PassPageView,
  SigninSimple as SigninSimpleView,
  Credit as CreditView,
  History as HistoryView,
  NotFoundCover as NotFoundCoverView,
  ErrorCover as ErrorCoverView
} from 'views';

const routes = [
  {
    path: '/',
    renderer: (params = {}): JSX.Element => <DesktopAppView {...params} />
  },
  {
    path: '/legal',
    renderer: (params = {}): JSX.Element => <LegalPageView {...params} />
  },
  {
    path: '/history',
    renderer: (params = {}): JSX.Element => <HistoryView {...params} />
  },
  {
    path: '/credit',
    renderer: (params = {}): JSX.Element => <CreditView {...params} />
  },
  {
    path: '/pass/:id',
    renderer: (params = {}): JSX.Element => <PassPageView {...params} />
  },
  {
    path: '/welcome',
    renderer: (params = {}): JSX.Element => <SigninSimpleView {...params} />
  },
  {
    path: '/not-found',
    renderer: (params = {}): JSX.Element => <NotFoundCoverView {...params} />
  },
  {
    path: '/error',
    renderer: (params = {}): JSX.Element => <ErrorCoverView {...params} />
  }
];

export default routes;
