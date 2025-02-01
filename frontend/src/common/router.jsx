import { lazy } from 'react';

import { createBrowserRouter } from 'react-router-dom';

import Layout from '@features/Layout/Layout'; // can't lazy load this
import ResetPassword from '@features/LandingPage/Authentication/ResetPassword/ResetPassword';

const Overview = lazy(() => import('@features/Home/Overview'));
const Reports = lazy(() => import('@features/Reports/Reports'));
const NotesList = lazy(() => import('@features/Notes/NotesList'));

const ProfilePage = lazy(() => import('@features/Profile/ProfilePage'));
const CategoryList = lazy(() => import('@features/Categories/CategoryList'));
const AssetList = lazy(() => import('@features/Assets/AssetList'));
const RecentActivityList = lazy(() => import('@features/RecentActivities/RecentActivityList'));

const SelectedAsset = lazy(() => import('@features/SelectedAsset/SelectedAsset'));
const MaintenancePlanList = lazy(() => import('@features/MaintenancePlan/MaintenancePlanList'));
const CategoryItemDetails = lazy(() => import('@features/CategoryItemDetails/CategoryItemDetails'));

const MaintenancePlanItemDetails = lazy(() =>
  import('@features/MaintenancePlanItemDetails/MaintenancePlanItemDetails')
);

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        path: '/',
        element: <Overview />,
      },
      {
        path: '/inventories/list',
        element: <AssetList />,
      },
      {
        path: '/inventories/:id/update',
        element: <SelectedAsset />,
      },
      {
        path: '/categories/list',
        element: <CategoryList />,
      },
      {
        path: '/category/:id',
        element: <CategoryItemDetails />,
      },
      {
        path: '/plans/list',
        element: <MaintenancePlanList />,
      },
      {
        path: '/plan/:id',
        element: <MaintenancePlanItemDetails />,
      },
      {
        path: '/reports',
        element: <Reports />,
      },
      {
        path: '/profile',
        element: <ProfilePage />,
      },
      {
        path: '/profile/notes',
        element: <NotesList />,
      },
      {
        path: 'recent/activities',
        element: <RecentActivityList />,
      },
    ],
  },
  {
    path: '/reset',
    element: <ResetPassword />,
  },
]);
