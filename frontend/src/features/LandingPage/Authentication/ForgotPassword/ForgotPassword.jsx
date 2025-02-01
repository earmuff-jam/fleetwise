import { useState } from 'react';

import { produce } from 'immer';

import { useDispatch } from 'react-redux';

import { enqueueSnackbar } from 'notistack';
import { Button, Stack } from '@mui/material';
import { authActions } from '@features/LandingPage/authSlice';
import { RESET_PASSWORD_FIELDS } from '@features/LandingPage/constants';
import ForgotPasswordFormFields from '@features/LandingPage/Authentication/ForgotPassword/ForgotPasswordFormFields';

export default function ForgotPassword({ handleClose, handlePageTransition }) {
  const dispatch = useDispatch();
  const [formFields, setFormFields] = useState(RESET_PASSWORD_FIELDS);

  const handleInput = (event) => {
    const { name, value } = event.target;
    setFormFields(
      produce(formFields, (draft) => {
        draft[name].value = value;
        draft[name].errorMsg = '';

        for (const validator of draft[name].validators) {
          if (validator.validate(value)) {
            draft[name].errorMsg = validator.message;
            break;
          }
        }
      })
    );
  };

  const validate = (formFields) => {
    const containsErr = Object.values(formFields).reduce((acc, el) => {
      if (el.errorMsg) {
        return true;
      }
      return acc;
    }, false);

    const requiredFormFields = Object.values(formFields).filter((v) => v.required);
    return containsErr || requiredFormFields.some((el) => el.value.trim() === '');
  };

  const submit = (e) => {
    e.preventDefault();

    if (validate(formFields)) {
      return;
    } else {
      const formattedData = Object.values(formFields).reduce((acc, el) => {
        if (el.value) {
          acc[el.name] = el.value;
        }
        return acc;
      }, {});

      dispatch(authActions.resetPassword(formattedData));
      enqueueSnackbar('Sent email notification to reset password.', {
        variant: 'success',
      });
    }
    handleClose();
    handlePageTransition();
  };

  return (
    <Stack spacing={1}>
      <ForgotPasswordFormFields formFields={formFields} handleInput={handleInput} />
      <Button
        variant="text"
        disabled={validate(formFields)}
        disableRipple={true}
        disableFocusRipple={true}
        onClick={submit}
      >
        Reset password
      </Button>
    </Stack>
  );
}
