import { useContext, useState } from 'react';

import { Container, Stack } from '@mui/material';

import SimpleModal from '@common/SimpleModal';
import FloatingBarChart from '@common/Chart/FloatingBarChart';
import Login from '@features/LandingPage/Authentication/Login/Login';
import HeroContent from '@features/LandingPage/HeroContent/HeroContent';

import Signup from '@features/LandingPage/Authentication/Signup/Signup';
import StyledAppBar from '@features/LandingPage/StyledAppBar/StyledAppBar';
import { COLORS, MODAL_STATE, SAMPLE_DATA } from '@features/LandingPage/constants';
import ForgotPassword from '@features/LandingPage/Authentication/ForgotPassword/ForgotPassword';
import { PageContext } from '@src/ApplicationValidator';

export default function LandingPage() {
  const { setPage } = useContext(PageContext);
  const [modalState, setModalState] = useState(MODAL_STATE.NONE);

  const handleCloseModal = () => setModalState(MODAL_STATE.NONE);
  const openSignupModal = () => setModalState(MODAL_STATE.SIGN_UP);
  const openLoginModal = () => setModalState(MODAL_STATE.SIGN_IN);

  const handleForgotPasswordModal = () => setModalState(MODAL_STATE.FORGOT_PASSWORD);

  const handlePageTransition = () => setPage('reset');

  const formattedData = SAMPLE_DATA.map((v, index) => ({
    label: v.name,
    start: index === 0 ? 0 : SAMPLE_DATA[index - 1].price,
    end: v.price,
    color: COLORS[index % COLORS.length],
  }));

  return (
    <>
      <StyledAppBar />
      <Container maxWidth="md">
        <Stack spacing={2}>
          <HeroContent openSignupModal={openSignupModal} openLoginModal={openLoginModal} />
          <FloatingBarChart
            data={formattedData}
            backgroundColor={formattedData.map((d) => d.color)}
            borderColor={formattedData.map((d) => d.color)}
          />
        </Stack>

        {modalState === MODAL_STATE.SIGN_UP && (
          <SimpleModal
            title="Sign up"
            subtitle="Create an account to keep track of all your inventories."
            handleClose={handleCloseModal}
            maxSize="sm"
          >
            <Signup handleClose={handleCloseModal} />
          </SimpleModal>
        )}

        {modalState === MODAL_STATE.SIGN_IN && (
          <SimpleModal
            title="Sign in"
            subtitle="Login and manage your account."
            handleClose={handleCloseModal}
            maxSize="sm"
          >
            <Login handleClose={handleCloseModal} handleForgotPassword={handleForgotPasswordModal} />
          </SimpleModal>
        )}

        {modalState === MODAL_STATE.FORGOT_PASSWORD && (
          <SimpleModal
            title="Forgot Password"
            subtitle="Email sent to reset password"
            handleClose={handleCloseModal}
            maxSize="sm"
          >
            <ForgotPassword handleClose={handleCloseModal} handlePageTransition={handlePageTransition} />
          </SimpleModal>
        )}
      </Container>
    </>
  );
}
