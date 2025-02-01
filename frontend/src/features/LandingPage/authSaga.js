import { authActions } from './authSlice';

import { takeLatest, put, call, takeEvery } from 'redux-saga/effects';

import instance from '@utils/Instances';
import { REACT_APP_LOCALHOST_URL } from '@utils/Common';

const BASEURL = `${REACT_APP_LOCALHOST_URL}/api/v1`;

export function* fetchUserLogin(action) {
  try {
    const response = yield call(instance.post, `${BASEURL}/signin`, { ...action.payload });

    localStorage.setItem('userID', response.data.id);
    localStorage.setItem('isVerified', response.data.is_verified);
    yield put(authActions.getUserIDSuccess(response.data));
  } catch (e) {
    yield put(authActions.getUserIDFailure(e));
  }
}

export function* fetchUserSignup(action) {
  try {
    const response = yield call(instance.post, `${BASEURL}/signup`, { ...action.payload });

    yield put(authActions.getSignupSuccess(response.data));
  } catch (e) {
    yield put(authActions.getSignupFailure(e));
  }
}

export function* fetchIsValidUserEmail(action) {
  try {
    const response = yield call(instance.post, `${BASEURL}/isValidEmail`, { ...action.payload });

    yield put(authActions.isValidUserEmailSuccess(response.data));
  } catch (e) {
    yield put(authActions.isValidUserEmailFailure(e));
  }
}

export function* revalidateEmail(action) {
  try {
    const userID = localStorage.getItem('userID');
    const isVerified = localStorage.getItem('isVerified');
    const response = yield call(instance.post, `${BASEURL}/reset`, {
      ...action.payload,
      is_verified: isVerified === true,
      id: userID,
    });
    yield put(authActions.revalidateEmailSuccess(response.data));
  } catch (e) {
    yield put(authActions.revalidateEmailFailure(e));
  }
}

export function* fetchResetPassword(action) {
  try {
    yield call(instance.post, `${BASEURL}/resetPassword`, {
      ...action.payload,
    });
    yield put(authActions.resetPasswordSuccess());
  } catch (e) {
    yield put(authActions.resetPasswordFailure(e));
  }
}

export function* fetchUserLogout() {
  try {
    const response = yield call(instance.get, `${BASEURL}/logout`);
    yield put(authActions.getLogoutSuccess(response.data));
  } catch (e) {
    yield put(authActions.getLogoutFailure(e));
  }
}

export function* watchFetchUserLogin() {
  yield takeLatest('auth/getUserID', fetchUserLogin);
}

export function* watchFetchUserSignup() {
  yield takeLatest('auth/getSignup', fetchUserSignup);
}

export function* watchFetchResetPassword() {
  yield takeLatest('auth/resetPassword', fetchResetPassword);
}

export function* watchFetchIsValidUserEmail() {
  yield takeEvery(`auth/isValidUserEmail`, fetchIsValidUserEmail);
}

export function* watchRvalidateEmail() {
  yield takeEvery(`auth/revalidateEmail`, revalidateEmail);
}

export function* watchFetchUserLogout() {
  yield takeLatest(`auth/getLogout`, fetchUserLogout);
}

export default [
  watchFetchUserLogin,
  watchFetchUserSignup,
  watchFetchResetPassword,
  watchFetchIsValidUserEmail,
  watchRvalidateEmail,
  watchFetchUserLogout,
];
